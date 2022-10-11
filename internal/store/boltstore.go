// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-multierror"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/observiq/bindplane-op/internal/eventbus"
	"github.com/observiq/bindplane-op/internal/otlp/record"
	"github.com/observiq/bindplane-op/internal/store/search"
	"github.com/observiq/bindplane-op/internal/store/stats"
	"github.com/observiq/bindplane-op/model"
)

// bucket names
const (
	bucketResources    = "Resources"
	bucketTasks        = "Tasks"
	bucketAgents       = "Agents"
	bucketMeasurements = "Measurements"
)

type boltstore struct {
	db                 *bbolt.DB
	updates            *storeUpdates
	agentIndex         search.Index
	configurationIndex search.Index
	logger             *zap.Logger
	sync.RWMutex
	sessionStorage sessions.Store
}

var _ Store = (*boltstore)(nil)

// NewBoltStore returns a new store boltstore struct that implements the store.Store interface.
func NewBoltStore(ctx context.Context, db *bbolt.DB, options Options, logger *zap.Logger) Store {
	store := &boltstore{
		db:                 db,
		updates:            newStoreUpdates(ctx, options.MaxEventsToMerge),
		agentIndex:         search.NewInMemoryIndex("agent"),
		configurationIndex: search.NewInMemoryIndex("configuration"),
		logger:             logger,

		sessionStorage: newBPCookieStore(options.SessionsSecret),
	}

	// boltstore is not used for clusters, disconnect all agents
	store.disconnectAllAgents(context.Background())

	// start the timer that runs cleanup on measurements
	store.startMeasurements(ctx)

	// run the migration that fixes keys that were forcefully lowercased
	err := store.migrateBackToCaseSensitive(context.Background())
	if err != nil {
		logger.Error("failed to migrateBackToCaseSensitive", zap.Error(err))
	}

	return store
}

// InitDB takes in the full path to a storage file and returns an opened bbolt database.
// It will return an error if the file cannot be opened.
func InitDB(storageFilePath string) (*bbolt.DB, error) {
	var db, err = bbolt.Open(storageFilePath, 0640, nil)
	if err != nil {
		return nil, fmt.Errorf("error while opening bbolt storage file: %s, %w", storageFilePath, err)
	}

	buckets := []string{
		bucketResources,
		bucketTasks,
		bucketAgents,
		bucketMeasurements,
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			_, _ = tx.CreateBucketIfNotExists([]byte(bucket))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create bbolt storage bucket: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketMeasurements))
		for _, metric := range supportedMetricNames {
			_, _ = b.CreateBucketIfNotExists([]byte(metric))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create bbolt metrics bucket: %w", err)
	}

	return db, nil
}

// AgentConfiguration returns the configuration that should be applied to an agent.
func (s *boltstore) AgentConfiguration(ctx context.Context, agentID string) (*model.Configuration, error) {
	if agentID == "" {
		return nil, fmt.Errorf("cannot return configuration for empty agentID")
	}

	agent, err := s.Agent(ctx, agentID)
	if agent == nil {
		return nil, nil
	}

	// check for configuration= label and use that
	if configurationName, ok := agent.Labels.Set["configuration"]; ok {
		// if there is a configuration label, this takes precedence and we don't need to look any further
		return s.Configuration(ctx, configurationName)
	}

	var result *model.Configuration

	err = s.db.View(func(tx *bbolt.Tx) error {
		// iterate over the configurations looking for one that applies
		prefix := []byte(model.KindConfiguration)
		cursor := resourcesBucket(tx).Cursor()

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			configuration := &model.Configuration{}
			if err := json.Unmarshal(v, configuration); err != nil {
				s.logger.Error("unable to unmarshal configuration, ignoring", zap.Error(err))
				continue
			}
			if configuration.IsForAgent(agent) {
				result = configuration
				break
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve agent configuration: %w", err)
	}

	return result, nil
}

// AgentsIDsMatchingConfiguration returns the list of agent IDs that are using the specified configuration
func (s *boltstore) AgentsIDsMatchingConfiguration(ctx context.Context, configuration *model.Configuration) ([]string, error) {
	ids := s.AgentIndex(ctx).Select(configuration.Spec.Selector.MatchLabels)
	return ids, nil
}

func (s *boltstore) Updates() eventbus.Source[*Updates] {
	return s.updates.Updates()
}

// DeleteResources iterates threw a slice of resources, and removes them from storage by name.
// Sends any successful pipeline deletes to the pipelineDeletes channel, to be handled by the manager.
// Exporter and receiver deletes are sent to the manager via notifyUpdates.
func (s *boltstore) DeleteResources(ctx context.Context, resources []model.Resource) ([]model.ResourceStatus, error) {
	updates := NewUpdates()

	// track deleteStatuses to return
	deleteStatuses := make([]model.ResourceStatus, 0)

	for _, r := range resources {
		empty, err := model.NewEmptyResource(r.GetKind())
		if err != nil {
			deleteStatuses = append(deleteStatuses, *model.NewResourceStatusWithReason(r, model.StatusError, err.Error()))
			continue
		}

		deleted, exists, err := deleteResource(ctx, s, r.GetKind(), r.Name(), empty)

		switch err.(type) {
		case *DependencyError:
			deleteStatuses = append(
				deleteStatuses,
				*model.NewResourceStatusWithReason(r, model.StatusInUse, err.Error()))
			continue

		case nil:
			break

		default:
			deleteStatuses = append(deleteStatuses, *model.NewResourceStatusWithReason(r, model.StatusError, err.Error()))
			continue
		}

		if !exists {
			continue
		}

		deleteStatuses = append(deleteStatuses, *model.NewResourceStatus(r, model.StatusDeleted))
		updates.IncludeResource(deleted, EventTypeRemove)
	}

	s.notify(ctx, updates)
	return deleteStatuses, nil
}

// Apply resources iterates through a slice of resources, then adds them to storage,
// and calls notify updates on the updated resources.
func (s *boltstore) ApplyResources(ctx context.Context, resources []model.Resource) ([]model.ResourceStatus, error) {
	updates := NewUpdates()

	// resourceStatuses to return for the applied resources
	resourceStatuses := make([]model.ResourceStatus, 0)

	var errs error
	for _, resource := range resources {
		// Set the resource's initial ID, which wil be overwritten if
		// the resource already exists (using the existing resource ID)
		resource.EnsureID()

		warn, err := resource.ValidateWithStore(ctx, s)
		if err != nil {
			resourceStatuses = append(resourceStatuses, *model.NewResourceStatusWithReason(resource, model.StatusInvalid, err.Error()))
			continue
		}

		err = s.db.Update(func(tx *bbolt.Tx) error {
			// update the resource in the database
			status, err := upsertResource(tx, resource, resource.GetKind())
			if err != nil {
				resourceStatuses = append(resourceStatuses, *model.NewResourceStatusWithReason(resource, model.StatusError, err.Error()))
				return err
			}
			resourceStatuses = append(resourceStatuses, *model.NewResourceStatusWithReason(resource, status, warn))

			switch status {
			case model.StatusCreated:
				updates.IncludeResource(resource, EventTypeInsert)
			case model.StatusConfigured:
				updates.IncludeResource(resource, EventTypeUpdate)
			}

			// some resources need special treatment
			switch r := resource.(type) {
			case *model.Configuration:
				// update the index
				err = s.configurationIndex.Upsert(r)
				if err != nil {
					s.logger.Error("failed to update the search index", zap.String("configuration", r.Name()))
				}
			}
			return nil
		})
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	s.notify(ctx, updates)

	return resourceStatuses, errs
}

// ----------------------------------------------------------------------

func (s *boltstore) notify(ctx context.Context, updates *Updates) {
	err := updates.addTransitiveUpdates(ctx, s)
	if err != nil {
		// TODO: if we can't notify about all updates, what do we do?
		s.logger.Error("unable to add transitive updates", zap.Any("updates", updates), zap.Error(err))
	}
	if !updates.Empty() {
		s.updates.Send(updates)
	}
}

// ----------------------------------------------------------------------

// Clear clears the db store of resources, agents, and tasks.  Mostly used for testing.
func (s *boltstore) Clear() {
	// Disregarding error from update because these actions errors are known and prevented
	_ = s.db.Update(func(tx *bbolt.Tx) error {
		// Delete all the buckets.
		// Disregarding errors because it will only error if the bucket doesn't exist
		// or isn't a bucket key - which we're confident its not.
		_ = tx.DeleteBucket([]byte(bucketResources))
		_ = tx.DeleteBucket([]byte(bucketTasks))
		_ = tx.DeleteBucket([]byte(bucketAgents))
		_ = tx.DeleteBucket([]byte(bucketMeasurements))

		// create them again
		// Disregarding errors because bucket names are valid.
		_, _ = tx.CreateBucketIfNotExists([]byte(bucketResources))
		_, _ = tx.CreateBucketIfNotExists([]byte(bucketTasks))
		_, _ = tx.CreateBucketIfNotExists([]byte(bucketAgents))
		b, _ := tx.CreateBucketIfNotExists([]byte(bucketMeasurements))

		for _, metric := range supportedMetricNames {
			_, _ = b.CreateBucketIfNotExists([]byte(metric))
		}
		return nil
	})
}

func (s *boltstore) UpsertAgents(ctx context.Context, agentIDs []string, updater AgentUpdater) ([]*model.Agent, error) {
	ctx, span := tracer.Start(ctx, "store/UpsertAgents")
	defer span.End()

	agents := make([]*model.Agent, 0, len(agentIDs))
	updates := NewUpdates()

	err := s.db.Update(func(tx *bbolt.Tx) error {
		for _, agentID := range agentIDs {
			agent, err := upsertAgentTx(tx, agentID, updater, updates)
			if err != nil {
				return err
			}

			agents = append(agents, agent)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// update the search index with changes
	for _, a := range agents {
		if err := s.agentIndex.Upsert(a); err != nil {
			s.logger.Error("failed to update the search index", zap.String("agentID", a.ID))
		}
	}

	// notify results
	s.notify(ctx, updates)
	return agents, nil
}

// UpsertAgent creates or updates the given agent and calls the updater method on it.
func (s *boltstore) UpsertAgent(ctx context.Context, id string, updater AgentUpdater) (*model.Agent, error) {
	ctx, span := tracer.Start(ctx, "store/UpsertAgent")
	defer span.End()

	var updatedAgent *model.Agent
	updates := NewUpdates()

	err := s.db.Update(func(tx *bbolt.Tx) error {
		agent, err := upsertAgentTx(tx, id, updater, updates)
		updatedAgent = agent
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// update the index
	err = s.agentIndex.Upsert(updatedAgent)
	if err != nil {
		s.logger.Error("failed to update the search index", zap.String("agentID", updatedAgent.ID))
	}

	s.notify(ctx, updates)

	return updatedAgent, nil
}

func (s *boltstore) Agents(ctx context.Context, options ...QueryOption) ([]*model.Agent, error) {
	opts := makeQueryOptions(options)

	// search is implemented using the search index
	if opts.query != nil {
		ids, err := s.agentIndex.Search(ctx, opts.query)
		if err != nil {
			return nil, err
		}
		return s.agentsByID(ids, opts)
	}

	agents := []*model.Agent{}

	err := s.db.View(func(tx *bbolt.Tx) error {
		cursor := agentBucket(tx).Cursor()
		prefix := []byte("Agent")

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			agent := &model.Agent{}
			if err := json.Unmarshal(v, agent); err != nil {
				return fmt.Errorf("agents: %w", err)
			}

			if opts.selector.Matches(agent.Labels) {
				agents = append(agents, agent)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if opts.sort == "" {
		opts.sort = "name"
	}
	return applySortOffsetAndLimit(agents, opts, func(field string, item *model.Agent) string {
		switch field {
		case "id":
			return item.ID
		default:
			return item.Name
		}
	}), nil
}

func (s *boltstore) DeleteAgents(ctx context.Context, agentIDs []string) ([]*model.Agent, error) {
	updates := NewUpdates()
	deleted := make([]*model.Agent, 0, len(agentIDs))

	err := s.db.Update(func(tx *bbolt.Tx) error {
		c := agentBucket(tx).Cursor()

		for _, id := range agentIDs {
			agentKey := agentKey(id)
			k, v := c.Seek(agentKey)

			if k != nil && bytes.Equal(k, agentKey) {

				// Save the agent to return and set its status to deleted.
				agent := &model.Agent{}
				err := json.Unmarshal(v, agent)
				if err != nil {
					return err
				}

				agent.Status = model.Deleted
				deleted = append(deleted, agent)

				// delete it
				err = c.Delete()
				if err != nil {
					return err
				}

				// include it in updates
				updates.IncludeAgent(agent, EventTypeRemove)
			}
		}

		return nil
	})

	if err != nil {
		return deleted, err
	}

	// remove deleted agents from the index
	for _, agent := range deleted {
		if err := s.agentIndex.Remove(agent); err != nil {
			s.logger.Error("failed to remove from the search index", zap.String("agentID", agent.ID))
		}
	}

	// notify updates
	s.notify(ctx, updates)

	return deleted, nil
}

func (s *boltstore) agentsByID(ids []string, opts queryOptions) ([]*model.Agent, error) {
	var agents []*model.Agent

	err := s.db.View(func(tx *bbolt.Tx) error {
		for _, id := range ids {
			data := agentBucket(tx).Get(agentKey(id))
			if data == nil {
				return nil
			}
			agent := &model.Agent{}
			if err := json.Unmarshal(data, agent); err != nil {
				return fmt.Errorf("agents: %w", err)
			}

			if opts.selector.Matches(agent.Labels) {
				agents = append(agents, agent)
			}
		}
		return nil
	})

	return agents, err
}

func (s *boltstore) AgentsCount(ctx context.Context, options ...QueryOption) (int, error) {
	agents, err := s.Agents(ctx, options...)
	if err != nil {
		return -1, err
	}
	return len(agents), nil
}

func (s *boltstore) Agent(ctx context.Context, id string) (*model.Agent, error) {
	var agent *model.Agent

	err := s.db.View(func(tx *bbolt.Tx) error {
		data := agentBucket(tx).Get(agentKey(id))
		if data == nil {
			return nil
		}
		agent = &model.Agent{}
		return json.Unmarshal(data, agent)
	})

	return agent, err
}

func (s *boltstore) AgentVersion(ctx context.Context, name string) (*model.AgentVersion, error) {
	item, exists, err := resource[*model.AgentVersion](s, model.KindAgentVersion, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) AgentVersions(ctx context.Context) ([]*model.AgentVersion, error) {
	result, err := resources[*model.AgentVersion](s, model.KindAgentVersion)
	if err == nil {
		model.SortAgentVersionsLatestFirst(result)
	}
	return result, err
}
func (s *boltstore) DeleteAgentVersion(ctx context.Context, name string) (*model.AgentVersion, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindAgentVersion, name, &model.AgentVersion{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) Configurations(ctx context.Context, options ...QueryOption) ([]*model.Configuration, error) {
	opts := makeQueryOptions(options)
	// search is implemented using the search index
	if opts.query != nil {
		names, err := s.configurationIndex.Search(ctx, opts.query)
		if err != nil {
			return nil, err
		}
		return resourcesByName[*model.Configuration](s, model.KindConfiguration, names, opts)
	}

	return resourcesWithFilter(s, model.KindConfiguration, func(c *model.Configuration) bool {
		return opts.selector.Matches(c.GetLabels())
	})
}
func (s *boltstore) Configuration(ctx context.Context, name string) (*model.Configuration, error) {
	item, exists, err := resource[*model.Configuration](s, model.KindConfiguration, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) DeleteConfiguration(ctx context.Context, name string) (*model.Configuration, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindConfiguration, name, &model.Configuration{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) Source(ctx context.Context, name string) (*model.Source, error) {
	item, exists, err := resource[*model.Source](s, model.KindSource, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) Sources(ctx context.Context) ([]*model.Source, error) {
	return resources[*model.Source](s, model.KindSource)
}
func (s *boltstore) DeleteSource(ctx context.Context, name string) (*model.Source, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindSource, name, &model.Source{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) SourceType(ctx context.Context, name string) (*model.SourceType, error) {
	item, exists, err := resource[*model.SourceType](s, model.KindSourceType, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) SourceTypes(ctx context.Context) ([]*model.SourceType, error) {
	return resources[*model.SourceType](s, model.KindSourceType)
}
func (s *boltstore) DeleteSourceType(ctx context.Context, name string) (*model.SourceType, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindSourceType, name, &model.SourceType{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) Processor(ctx context.Context, name string) (*model.Processor, error) {
	item, exists, err := resource[*model.Processor](s, model.KindProcessor, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) Processors(ctx context.Context) ([]*model.Processor, error) {
	return resources[*model.Processor](s, model.KindProcessor)
}
func (s *boltstore) DeleteProcessor(ctx context.Context, name string) (*model.Processor, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindProcessor, name, &model.Processor{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) ProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	item, exists, err := resource[*model.ProcessorType](s, model.KindProcessorType, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) ProcessorTypes(ctx context.Context) ([]*model.ProcessorType, error) {
	return resources[*model.ProcessorType](s, model.KindProcessorType)
}
func (s *boltstore) DeleteProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindProcessorType, name, &model.ProcessorType{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) Destination(ctx context.Context, name string) (*model.Destination, error) {
	item, exists, err := resource[*model.Destination](s, model.KindDestination, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) Destinations(ctx context.Context) ([]*model.Destination, error) {
	return resources[*model.Destination](s, model.KindDestination)
}
func (s *boltstore) DeleteDestination(ctx context.Context, name string) (*model.Destination, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindDestination, name, &model.Destination{})
	if !exists {
		return nil, err
	}
	return item, err
}

func (s *boltstore) DestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	item, exists, err := resource[*model.DestinationType](s, model.KindDestinationType, name)
	if !exists {
		item = nil
	}
	return item, err
}
func (s *boltstore) DestinationTypes(ctx context.Context) ([]*model.DestinationType, error) {
	return resources[*model.DestinationType](s, model.KindDestinationType)
}
func (s *boltstore) DeleteDestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	item, exists, err := deleteResourceAndNotify(ctx, s, model.KindDestinationType, name, &model.DestinationType{})
	if !exists {
		return nil, err
	}
	return item, err
}

// CleanupDisconnectedAgents removes agents that have disconnected before the specified time
func (s *boltstore) CleanupDisconnectedAgents(ctx context.Context, since time.Time) error {
	agents, err := s.Agents(ctx)
	if err != nil {
		return err
	}
	changes := NewUpdates()

	for _, agent := range agents {
		if agent.DisconnectedSince(since) {
			err := s.db.Update(func(tx *bbolt.Tx) error {
				return agentBucket(tx).Delete(agentKey(agent.ID))
			})
			if err != nil {
				return err
			}
			changes.IncludeAgent(agent, EventTypeRemove)

			// update the index
			if err := s.agentIndex.Remove(agent); err != nil {
				s.logger.Error("failed to remove from the search index", zap.String("agentID", agent.ID))
			}
		}
	}

	s.notify(ctx, changes)
	return nil
}

// Index provides access to the search Index implementation managed by the Store
func (s *boltstore) AgentIndex(ctx context.Context) search.Index {
	return s.agentIndex
}

// ConfigurationIndex provides access to the search Index for Configurations
func (s *boltstore) ConfigurationIndex(ctx context.Context) search.Index {
	return s.configurationIndex
}

func (s *boltstore) UserSessions() sessions.Store {
	return s.sessionStorage
}

// Measurements stores stats for agents and configurations
func (s *boltstore) Measurements() stats.Measurements {
	return s
}

// ----------------------------------------------------------------------

func (s *boltstore) disconnectAllAgents(ctx context.Context) {
	if agents, err := s.Agents(ctx); err != nil {
		s.logger.Error("error while disconnecting all agents on startup", zap.Error(err))
	} else {
		s.logger.Info("disconnecting all agents on startup", zap.Int("count", len(agents)))
		for _, agent := range agents {
			_, err := s.UpsertAgent(ctx, agent.ID, func(a *model.Agent) {
				a.Disconnect()
			})
			if err != nil {
				s.logger.Error("error while disconnecting agent on startup", zap.Error(err))
			}
		}
	}
}

/* ---------------------------- helper functions ---------------------------- */
func resourcesPrefix(kind model.Kind) []byte {
	return []byte(fmt.Sprintf("%s|", kind))
}

func resourceKey(kind model.Kind, name string) []byte {
	return []byte(fmt.Sprintf("%s|%s", kind, name))
}

func agentKey(id string) []byte {
	return []byte(fmt.Sprintf("%s|%s", "Agent", id))
}

func keyFromResource(r model.Resource) []byte {
	if r == nil || r.GetKind() == model.KindUnknown {
		return make([]byte, 0)
	}
	return resourceKey(r.GetKind(), r.Name())
}

/* --------------------------- transaction helpers -------------------------- */
/* ------- These helper functions happen inside of a bbolt transaction ------ */
func agentBucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(bucketAgents))
}

func resourcesBucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(bucketResources))
}

func measurementsBucket(tx *bbolt.Tx, metric string) *bbolt.Bucket {
	b := tx.Bucket([]byte(bucketMeasurements))
	if b != nil {
		return b.Bucket([]byte(metric))
	}
	return nil
}

func upsertResource(tx *bbolt.Tx, r model.Resource, kind model.Kind) (model.UpdateStatus, error) {
	key := resourceKey(kind, r.Name())
	bucket := resourcesBucket(tx)
	existing := bucket.Get(key)

	// preserve the id (if possible)
	if len(existing) > 0 {
		var cur model.AnyResource
		if err := json.Unmarshal(existing, &cur); err == nil {
			r.SetID(cur.ID())
		}
	}

	data, err := json.Marshal(r)
	if err != nil {
		// error, status unchanged
		return model.StatusUnchanged, fmt.Errorf("upsert resource: %w", err)
	}
	if bytes.Equal(existing, data) {
		return model.StatusUnchanged, nil
	}

	if err = bucket.Put(key, data); err != nil {
		// error, status unchanged
		return model.StatusUnchanged, fmt.Errorf("upsert resource: %w", err)
	}

	if len(existing) == 0 {
		return model.StatusCreated, nil
	}
	return model.StatusConfigured, nil
}

// upsertAgentTx is a transaction helper that updates the given agent,
// puts it into the agent bucket  and includes it in the passed updates.
// it does *not* update the search index or notify any subscribers of updates.
func upsertAgentTx(tx *bbolt.Tx, agentID string, updater AgentUpdater, updates *Updates) (*model.Agent, error) {
	bucket := agentBucket(tx)
	key := agentKey(agentID)

	agentEventType := EventTypeInsert
	agent := &model.Agent{ID: agentID}

	// load the existing agent or create it
	if data := bucket.Get(key); data != nil {
		// existing agent, unmarshal
		if err := json.Unmarshal(data, agent); err != nil {
			return agent, err
		}
		agentEventType = EventTypeUpdate
	}

	// compare labels before/after and notify if they change
	labelsBefore := agent.Labels.String()

	// update the agent
	updater(agent)

	labelsAfter := agent.Labels.String()

	// if the labels changes is this is just an update, use EventTypeLabel
	if labelsAfter != labelsBefore && agentEventType == EventTypeUpdate {
		agentEventType = EventTypeLabel
	}

	// marshal it back to to json
	data, err := json.Marshal(agent)
	if err != nil {
		return agent, err
	}

	err = bucket.Put(key, data)
	if err != nil {
		return agent, err
	}

	updates.IncludeAgent(agent, agentEventType)
	return agent, nil
}

// ----------------------------------------------------------------------
// generic resource accessors

func resource[R model.Resource](s *boltstore, kind model.Kind, name string) (resource R, exists bool, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		key := resourceKey(kind, name)
		data := resourcesBucket(tx).Get(key)
		if data == nil {
			return nil
		}
		exists = true
		return json.Unmarshal(data, &resource)
	})
	return resource, exists, err
}

func resources[R model.Resource](s *boltstore, kind model.Kind) ([]R, error) {
	return resourcesWithFilter[R](s, kind, nil)
}
func resourcesWithFilter[R model.Resource](s *boltstore, kind model.Kind, include func(R) bool) ([]R, error) {
	var resources []R

	err := s.db.View(func(tx *bbolt.Tx) error {
		prefix := resourcesPrefix(kind)
		cursor := resourcesBucket(tx).Cursor()

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			var resource R
			if err := json.Unmarshal(v, &resource); err != nil {
				// TODO(andy): if it can't be unmarshaled, it should probably be removed from the store. ignore it for now.
				s.logger.Error("failed to unmarshal resource", zap.String("key", string(k)), zap.String("kind", string(kind)), zap.Error(err))
				continue
			}
			if include == nil || include(resource) {
				resources = append(resources, resource)
			}
		}

		return nil
	})

	return resources, err
}

// resourcesByName returns the resources of the specified name with the specified names. If requesting some resources
// results in an error, the errors will be accumulated and return with the list of resources successfully retrieved.
func resourcesByName[R model.Resource](s *boltstore, kind model.Kind, names []string, opts queryOptions) ([]R, error) {
	var errs error
	var results []R

	for _, name := range names {
		if result, exists, err := resource[R](s, kind, name); err != nil {
			errs = multierror.Append(errs, err)
		} else {
			if exists && opts.selector.Matches(result.GetLabels()) {
				results = append(results, result)
			}
		}
	}
	return results, errs
}

func deleteResourceAndNotify[R model.Resource](ctx context.Context, s *boltstore, kind model.Kind, name string, emptyResource R) (resource R, exists bool, err error) {
	deleted, exists, err := deleteResource(ctx, s, kind, name, emptyResource)

	if err == nil && exists {
		updates := NewUpdates()
		updates.IncludeResource(deleted, EventTypeRemove)
		s.notify(ctx, updates)
	}

	return deleted, exists, err
}

// deleteResource removes the resource with the given kind and name. Returns ResourceMissingError if the resource wasn't
// found. Returns DependencyError if the resource is referenced by another.
// emptyResource will be populated with the deleted resource. For convenience, if the delete is successful, the
// populated resource will also be returned. If there was an error, nil will be returned for the resource.
func deleteResource[R model.Resource](ctx context.Context, s *boltstore, kind model.Kind, name string, emptyResource R) (resource R, exists bool, err error) {
	var dependencies DependentResources

	err = s.db.Update(func(tx *bbolt.Tx) error {
		key := resourceKey(kind, name)

		c := resourcesBucket(tx).Cursor()
		k, v := c.Seek(key)

		if bytes.Equal(k, key) {
			// populate the emptyResource with the data before deleting
			err := json.Unmarshal(v, emptyResource)
			if err != nil {
				return err
			}

			exists = true

			// Check if the resources is referenced by another
			dependencies, err = FindDependentResources(ctx, s, emptyResource)
			if !dependencies.empty() {
				return ErrResourceInUse
			}

			// Delete the key from the store
			return c.Delete()
		}

		return ErrResourceMissing
	})

	switch {
	case errors.Is(err, ErrResourceMissing):
		return resource, exists, nil
	case errors.Is(err, ErrResourceInUse):
		return emptyResource, exists, newDependencyError(dependencies)
	case err != nil:
		return resource, exists, err
	}

	if emptyResource.GetKind() == model.KindConfiguration {
		if err := s.configurationIndex.Remove(emptyResource); err != nil {
			s.logger.Error("failed to remove configuration from the search index", zap.String("name", emptyResource.Name()))
		}
	}

	return emptyResource, exists, nil
}

// getObjectIds will retrieve identifiers for all objects in a bucket where the keys are formatted KIND|IDENTIFIER
func (s *boltstore) getObjectIds(bucketFunc func(*bbolt.Tx) *bbolt.Bucket, kind model.Kind) ([]string, error) {
	ids := []string{}
	prefix := []byte(fmt.Sprintf("%s|", kind))
	err := s.db.View(func(tx *bbolt.Tx) error {
		cursor := bucketFunc(tx).Cursor()
		for k, _ := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			if _, id, found := strings.Cut(string(k), "|"); found {
				ids = append(ids, id)
			}
		}
		return nil
	})
	return ids, err
}

// ----------------------------------------------------------------------
// Measurements implementation

const logDataSizeName = "otelcol_processor_throughputmeasurement_log_data_size"
const metricDataSizeName = "otelcol_processor_throughputmeasurement_metric_data_size"
const traceDataSizeName = "otelcol_processor_throughputmeasurement_trace_data_size"
const measurementsDateFormat = "2006-01-02T15:04:05"

var supportedMetricNames = []string{
	logDataSizeName,
	metricDataSizeName,
	traceDataSizeName,
}

// AgentMetrics provides metrics for an individual agents. They are essentially configuration metrics filtered to a
// list of agents.
//
// Note: While the same record.Metric struct is used to return the metrics, these are not the same metrics provided to
// Store. They will be aggregated and counter metrics will be converted into rates.
func (s *boltstore) AgentMetrics(ctx context.Context, ids []string, options ...stats.QueryOption) (stats.MetricData, error) {
	// Empty string single key or empty array of ids is a request for all Agents
	if len(ids) == 0 || (len(ids) == 1 && ids[0] == "") {
		var err error
		ids, err = s.getObjectIds(agentBucket, model.KindAgent)
		if err != nil {
			return nil, err
		}
	}
	return s.retrieveMetrics(ctx, supportedMetricNames, string(model.KindAgent), ids, options...)
}

// ConfigurationMetrics provides all metrics associated with a configuration aggregated from all agents using the
// configuration.
//
// Note: While the same record.Metric struct is used to return the metrics, these are not the same metrics provided to
// Store. They will be aggregated and counter metrics will be converted into rates.
func (s *boltstore) ConfigurationMetrics(ctx context.Context, name string, options ...stats.QueryOption) (stats.MetricData, error) {
	names := []string{name}
	var err error
	// Empty name is a request for all configurations
	if name == "" {
		names, err = s.getObjectIds(resourcesBucket, model.KindConfiguration)
		if err != nil {
			return nil, err
		}
	}

	baseMetrics, err := s.retrieveMetrics(ctx, supportedMetricNames, string(model.KindConfiguration), names, options...)
	if err != nil {
		return nil, err
	}

	groupedMetrics := map[string]stats.MetricData{}
	for _, m := range baseMetrics {
		// since multiple configurations may be returned for the overview page, group by configuration and processor
		key := fmt.Sprintf("%s|%s", stats.Configuration(m), stats.Processor(m))
		groupedMetrics[key] = append(groupedMetrics[key], m)
	}

	finalMetrics := stats.MetricData{}
	for _, metrics := range groupedMetrics {
		sum := 0.0
		for _, m := range metrics {
			val, _ := stats.Value(m)
			sum += val
		}

		attributes := map[string]interface{}{
			stats.ConfigurationAttributeName: metrics[0].Attributes[stats.ConfigurationAttributeName],
			stats.ProcessorAttributeName:     metrics[0].Attributes[stats.ProcessorAttributeName],
		}

		m := generateRecord(metrics[0], sum, attributes)

		finalMetrics = append(finalMetrics, m)
	}

	return finalMetrics, nil
}

// OverviewMetrics provides all metrics needed for the overview page. This page shows configs and destinations.
func (s *boltstore) OverviewMetrics(ctx context.Context, options ...stats.QueryOption) (stats.MetricData, error) {
	return s.ConfigurationMetrics(ctx, "", options...)
}

func (s *boltstore) retrieveMetrics(ctx context.Context, metricNames []string, objectType string, ids []string, options ...stats.QueryOption) (stats.MetricData, error) {
	result := stats.MetricData{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		var errs error
		opts := stats.MakeQueryOptions(options)
		for _, metricName := range metricNames {
			cursor := measurementsBucket(tx, metricName).Cursor()

			// Always start looking at a recent non-current chunk data point, as current
			// bucket is not going to have all data for all resources if the frame is Now()
			frame := opts.EndTime.Add(-10 * time.Second).UTC()
			startDate := frame.Add(-1 * opts.Period).Truncate(getRollupFromPeriod(opts)).Format(measurementsDateFormat)
			endDate := frame.Truncate(getRollupFromPeriod(opts))
			endDateString := endDate.Format(measurementsDateFormat)

			for _, id := range ids {
				endMetrics, err := findEndMetrics(cursor, endDateString, objectType, id)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}

				if len(endMetrics) == 0 {
					continue
				}

				// List the keys for which we want to find matching start points.
				desiredKeys := map[string]interface{}{}
				for k := range endMetrics {
					desiredKeys[k] = true
				}

				startMetrics, err := findStartMetrics(cursor, startDate, endDate, objectType, id, desiredKeys)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}

				// Only calculate rates if we found a start point to match the end points desired
				for key, first := range startMetrics {
					if last, ok := endMetrics[key]; ok {
						if metric := calculateRateMetric(first, last); metric != nil {
							result = append(result, metric)
						}
					}
				}
			}
		}

		return errs
	})
	return result, err
}

func findEndMetrics(c *bbolt.Cursor, time, objectType, id string) (map[string]*record.Metric, error) {
	identifier := fmt.Sprintf("%s|%s", objectType, sanitizeKey(id))
	prefix := []byte(fmt.Sprintf("%s|%s|", identifier, time))
	metrics := map[string]*record.Metric{}

	var k, v []byte
	for k, v = c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var m *record.Metric
		m = &record.Metric{}
		if err := json.Unmarshal(v, m); err != nil {
			return nil, err
		}
		metrics[removeTimestampFromKey(k)] = m
	}

	return metrics, nil
}

func findStartMetrics(c *bbolt.Cursor, time string, endTime time.Time, objectType, id string, desiredKeys map[string]interface{}) (map[string]*record.Metric, error) {
	identifier := fmt.Sprintf("%s|%s", objectType, sanitizeKey(id))
	startingPrefix := []byte(fmt.Sprintf("%s|%s|", identifier, time))
	metrics := map[string]*record.Metric{}

	continueToSeek := func(k []byte) bool {
		// Stop searching if we've satisfied all of the keys we were looking for
		if len(desiredKeys) == 0 {
			return false
		}
		// Stop searching if we've hit the end of the data
		if k == nil {
			return false
		}

		// Stop searching if we're no longer looking at data for the right object
		if !bytes.HasPrefix(k, []byte(identifier)) {
			return false
		}

		return true
	}

	var k, v []byte
	for k, v = c.Seek(startingPrefix); continueToSeek(k); k, v = c.Next() {
		// If we've been passed a desiredKeys list, only unmarshal & process
		// data points on the requested list
		keyWithoutTimestamp := removeTimestampFromKey(k)
		if _, ok := desiredKeys[keyWithoutTimestamp]; !ok {
			continue
		}
		var m *record.Metric
		m = &record.Metric{}
		if err := json.Unmarshal(v, m); err != nil {
			return nil, err
		}

		// If we've reached or passed the endTime for this query, stop searching
		if m.Timestamp.Sub(endTime) >= 0 {
			break
		}

		metrics[keyWithoutTimestamp] = m
		delete(desiredKeys, keyWithoutTimestamp)
	}

	return metrics, nil
}

func removeTimestampFromKey(k []byte) string {
	keyParts := strings.Split(string(k), "|")
	keyParts = append(keyParts[:2], keyParts[3:]...)
	return strings.Join(keyParts, "|")
}

func getRollupFromPeriod(opts stats.QueryOptions) time.Duration {
	switch opts.Period {
	case 1 * time.Minute:
		return 10 * time.Second
	case 5 * time.Minute:
		return 1 * time.Minute
	case 1 * time.Hour:
		return 5 * time.Minute
	case 24 * time.Hour:
		return 1 * time.Hour
	default:
		return 1 * time.Hour
	}
}

// SaveAgentMetrics saves new metrics. These metrics will be aggregated to determine metrics associated with agents and configurations.
func (s *boltstore) SaveAgentMetrics(ctx context.Context, metrics []*record.Metric) error {
	groupedMetrics := map[string][]*record.Metric{}
	for _, m := range metrics {
		groupedMetrics[m.Name] = append(groupedMetrics[m.Name], m)
	}

	var metricNames []string
	for _, metric := range metrics {
		metricNames = append(metricNames, metric.Name)
	}

	var errs error
	for _, metricName := range supportedMetricNames {
		if group, ok := groupedMetrics[metricName]; ok {
			err := s.storeMeasurements(ctx, metricName, group)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

// ProcessMetrics is called in the background at regular intervals and performs metric roll-up and removes old data
func (s *boltstore) ProcessMetrics(ctx context.Context) error {
	var errs error
	for _, m := range supportedMetricNames {
		if err := s.cleanupMeasurements(ctx, m); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (s *boltstore) startMeasurements(ctx context.Context) {
	// start the measurements timer for writing & rolling up
	go func() {
		measurementsTicker := time.NewTicker(time.Minute)
		defer measurementsTicker.Stop()

		for {
			select {
			case <-measurementsTicker.C:
				// periodically clean up old measurements
				if err := s.ProcessMetrics(ctx); err != nil {
					s.logger.Error("error cleaning up measurements", zap.Error(err))
				}
			case <-ctx.Done():
				// send anything left in the buffer before stopping
				return
			}
		}
	}()
}

func (s *boltstore) storeMeasurements(ctx context.Context, metricName string, metrics stats.MetricData) error {
	if len(metrics) == 0 {
		return nil
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		var errs error
		b := measurementsBucket(tx, metricName)

		for _, m := range metrics {
			data, err := json.Marshal(m)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			for _, key := range metricsKeys(m) {
				if err = b.Put(key, data); err != nil {
					errs = multierror.Append(errs, err)
				}
			}
		}
		return errs
	})
}

func (s *boltstore) cleanupMeasurements(ctx context.Context, metricName string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		var errs error
		c := measurementsBucket(tx, metricName).Cursor()

		// Capture now to re-use for all the date math
		now := time.Now()

		// Only look at data points older than 100 seconds, and we keep all the latest 10s intervals
		endCleanupDate := now.Add(-100 * time.Second)

		// Iterate through all points in the bucket, due to ordering we can't just scan by date
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if ts, err := time.Parse(measurementsDateFormat, keyTimestamp(k)); err != nil {
				errs = multierror.Append(errs, err)
			} else {
				// Assume anything older than the cutoff is being deleted, unless the
				// normalized timestamp matches one of the rollup times
				keep := ts.Sub(endCleanupDate) > 0
				// For last 10 minutes keep the minute data
				keep = keep || ts.Truncate(time.Minute*1).Equal(ts) && now.Sub(ts) <= 10*time.Minute
				// for the last 6 hours keep the 5 minute data
				keep = keep || ts.Truncate(time.Minute*5).Equal(ts) && now.Sub(ts) <= 6*time.Hour
				// for the last 24 hours keep the hourly data
				keep = keep || ts.Truncate(time.Hour*1).Equal(ts) && now.Sub(ts) <= 24*time.Hour
				// for the last 31 days keep the daily data
				keep = keep || ts.Truncate(time.Hour*24).Equal(ts) && now.Sub(ts) <= 24*31*time.Hour

				if !keep {
					if err := c.Delete(); err != nil {
						errs = multierror.Append(errs, err)
					}
				}
			}
		}

		return errs
	})
}

func calculateRateMetric(first, last *record.Metric) *record.Metric {
	// If we ended up with start & end being the same moment, or somehow
	// start is later than end, we don't want to provide a 0 or negative rate
	if first.Timestamp.Sub(last.Timestamp) >= 0 {
		return nil
	}

	var lastValue, firstValue float64
	var pass bool
	if lastValue, pass = stats.Value(last); !pass {
		return nil
	}
	if firstValue, pass = stats.Value(first); !pass {
		return nil
	}

	delta := lastValue - firstValue
	if delta < 0 {
		// TODO - handle rollover better than this
		delta = lastValue
	}

	duration := last.Timestamp.Sub(first.Timestamp)
	if duration <= 0*time.Second {
		return nil
	}

	rate := delta / duration.Seconds()
	// Reduce to 2 decimal places
	rate = math.Round(rate*100) / 100

	return generateRecord(last, rate, last.Attributes)
}

func generateRecord(source *record.Metric, value interface{}, attributes map[string]interface{}) *record.Metric {
	return &record.Metric{
		Name:       source.Name,
		Timestamp:  source.Timestamp.UTC(),
		Value:      value,
		Unit:       "B/s",
		Type:       "Rate",
		Attributes: attributes,
		Resource:   source.Resource,
	}
}

func keyTimestamp(k []byte) string {
	return strings.Split(string(k), "|")[2]
}

// metricsKey provides a key to store *record.Metric based on the agentID, processor and configuration of a single metric. It is
// assumed that all metrics come from a single agent using a single configuration.
func metricsKeys(m *record.Metric) [][]byte {
	normalizedTime := m.Timestamp.UTC().Truncate(10 * time.Second).Format(measurementsDateFormat)
	return [][]byte{
		[]byte(fmt.Sprintf("%s|%s|%s|%s|%s", model.KindAgent, stats.Agent(m), normalizedTime, stats.Configuration(m), stats.Processor(m))),
		[]byte(fmt.Sprintf("%s|%s|%s|%s|%s", model.KindConfiguration, stats.Configuration(m), normalizedTime, stats.Agent(m), stats.Processor(m))),
	}
}

func sanitizeKey(key string) string {
	return strings.ReplaceAll(key, "|", "")
}

// ----------------------------------------------------------------------

func getNameFromResource[R model.Resource](ctx context.Context, v []byte) (string, error) {
	var resource R
	if err := json.Unmarshal(v, &resource); err != nil {
		return "", err
	}
	return resource.Name(), nil
}

func (s *boltstore) migrateBackToCaseSensitive(ctx context.Context) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := resourcesBucket(tx)
		cursor := bucket.Cursor()

		for k, v := cursor.Seek(nil); k != nil; k, v = cursor.Next() {
			oldKey := string(k)
			kind, _, found := strings.Cut(oldKey, "|")
			if !found {
				continue
			}
			foundKind := model.ParseKind(string(kind))
			if foundKind == model.KindUnknown {
				continue
			}

			var foundName string
			var err error
			switch foundKind {
			case model.KindAgentVersion:
				foundName, err = getNameFromResource[*model.AgentVersion](ctx, v)
			case model.KindConfiguration:
				foundName, err = getNameFromResource[*model.Configuration](ctx, v)
			case model.KindSource:
				foundName, err = getNameFromResource[*model.Source](ctx, v)
			case model.KindSourceType:
				foundName, err = getNameFromResource[*model.SourceType](ctx, v)
			case model.KindProcessor:
				foundName, err = getNameFromResource[*model.Processor](ctx, v)
			case model.KindProcessorType:
				foundName, err = getNameFromResource[*model.ProcessorType](ctx, v)
			case model.KindDestination:
				foundName, err = getNameFromResource[*model.Destination](ctx, v)
			case model.KindDestinationType:
				foundName, err = getNameFromResource[*model.DestinationType](ctx, v)
			}

			if err != nil {
				continue
			}

			newKey := fmt.Sprintf("%s|%s", foundKind, foundName)
			if newKey != oldKey {
				err := cursor.Delete()
				if err != nil {
					return err
				}
				err = bucket.Put([]byte(newKey), v)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
