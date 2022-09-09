// Copyright  observIQ, Inc
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

// Package client provides a go client for interacting with the BindPlane OP server. Most of the functions depend on the
// BindPlane REST API internally.
package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/observiq/bindplane-op/common"
	"github.com/observiq/bindplane-op/internal/rest"
	"github.com/observiq/bindplane-op/internal/version"
	"github.com/observiq/bindplane-op/model"
)

// AgentInstallOptions contains configuration options used for installing an agent.
type AgentInstallOptions struct {
	// Platform is the platform the agent will run on, e.g. "linux"
	Platform string

	// Version is the agent release version to install. Available release versions of the observiq-otel-collector are
	// available at [observiq-otel-collector Releases]
	//
	// [observiq-otel-collector Releases]: https://github.com/observIQ/observiq-otel-collector/releases
	Version string

	// Labels is a string representation of the agents labels, e.g. "platform=dev,os=windows,app=nginx"
	Labels string

	// SecretKey is the secret key used to authenticate agents with BindPlane OP
	SecretKey string

	// RemoteURL is the URL that the agent will use to connect to BindPlane OP
	RemoteURL string
}

// ----------------------------------------------------------------------

// queryOptions represents the set of options available for a store query
type queryOptions struct {
	selector string
	query    string
	offset   int
	limit    int
	sort     string
}

func makeQueryOptions(options []QueryOption) queryOptions {
	opts := queryOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}

// QueryOption is an option used in many Store queries
type QueryOption func(*queryOptions)

// WithSelector adds a selector to the query options
func WithSelector(selector string) QueryOption {
	return func(opts *queryOptions) {
		opts.selector = selector
	}
}

// WithQuery adds a search query string to the query options
func WithQuery(query string) QueryOption {
	return func(opts *queryOptions) {
		opts.query = query
	}
}

// WithOffset sets the offset for the results to return. For paging, if the pages have 10 items per page and this is the
// 3rd page, set the offset to 20.
func WithOffset(offset int) QueryOption {
	return func(opts *queryOptions) {
		opts.offset = offset
	}
}

// WithLimit sets the maximum number of results to return. For paging, if the pages have 10 items per page, set the
// limit to 10.
func WithLimit(limit int) QueryOption {
	return func(opts *queryOptions) {
		opts.limit = limit
	}
}

// WithSort sets the sort order for the request. The sort value is the name of the field, sorted ascending. To sort
// descending, prefix the field with a minus sign (-). Some Stores only allow sorting by certain fields. Sort values not
// supported will be ignored.
func WithSort(field string) QueryOption {
	return func(opts *queryOptions) {
		opts.sort = field
	}
}

// BindPlane is a REST client for BindPlane OP.
type BindPlane interface {
	// Agents returns a list of Agents.
	Agents(ctx context.Context, options ...QueryOption) ([]*model.Agent, error)
	// Agent returns a single Agent.
	Agent(ctx context.Context, id string) (*model.Agent, error)
	// DeleteAgents deletes multiple agents by ID.
	DeleteAgents(ctx context.Context, agentIDs []string) ([]*model.Agent, error)

	// AgentVersions returns a list of AgentVersion resources.
	AgentVersions(ctx context.Context) ([]*model.AgentVersion, error)
	// AgentVersion returns a single AgentVersion resources by name.
	AgentVersion(ctx context.Context, name string) (*model.AgentVersion, error)
	// DeleteAgentVersion deletes an AgentVersion resource by name.
	DeleteAgentVersion(ctx context.Context, name string) error

	// SyncAgentVersions builds agent-version from the release data in GitHub.
	// If version is empty, it syncs the last 10 releases.
	SyncAgentVersions(ctx context.Context, version string) ([]*model.AnyResourceStatus, error)

	// Configurations returns a list of Configuration resources.
	Configurations(ctx context.Context) ([]*model.Configuration, error)
	// Configuration returns a single Configuration resource from GET /v1/configurations/:name
	Configuration(ctx context.Context, name string) (*model.Configuration, error)
	// Delete configuration deletes a single configuration reseource.
	DeleteConfiguration(ctx context.Context, name string) error
	// RawConfiguration returns the raw OpenTelemetry configuration for the configuration with
	// the specified name. This can either be the raw value of a raw configuration or the
	// rendered value of a configuration with sources and destinations.
	RawConfiguration(ctx context.Context, name string) (string, error)
	// CopyConfig creates a deep copy of an existing resource under a new name.
	CopyConfig(ctx context.Context, name, copyName string) error

	// Sources returns a list of all Source resources.
	Sources(ctx context.Context) ([]*model.Source, error)
	// Source returns a single Source resource by name.
	Source(ctx context.Context, name string) (*model.Source, error)
	// DeleteSource deletes a single Source resource by name.
	DeleteSource(ctx context.Context, name string) error

	// SourceTypes returns a list of all SourceType resources.
	SourceTypes(ctx context.Context) ([]*model.SourceType, error)
	// SourceType returns a single SourceType resource by name.
	SourceType(ctx context.Context, name string) (*model.SourceType, error)
	// DeleteSourceType deletes a single SourceType resource by name.
	DeleteSourceType(ctx context.Context, name string) error

	// Processors returns a list of all Processor resources.
	Processors(ctx context.Context) ([]*model.Processor, error)
	// Processor returns a single Processor resource by name.
	Processor(ctx context.Context, name string) (*model.Processor, error)
	// DeleteProcessor deletes a single Processor resource by name.
	DeleteProcessor(ctx context.Context, name string) error

	// ProcessorTypes returns a list of all ProcessorType resources.
	ProcessorTypes(ctx context.Context) ([]*model.ProcessorType, error)
	// ProcessorType returns a single ProcessorType resource by name.
	ProcessorType(ctx context.Context, name string) (*model.ProcessorType, error)
	// DeleteProcessorType deletes a single ProcessorType resource by name.
	DeleteProcessorType(ctx context.Context, name string) error

	// Destinations returns a list of all Destination resources.
	Destinations(ctx context.Context) ([]*model.Destination, error)
	// Destination returns a single Destination resource by name.
	Destination(ctx context.Context, name string) (*model.Destination, error)
	// DeleteDestination deletes a single Destination resource by name.
	DeleteDestination(ctx context.Context, name string) error

	// DestinationTypes returns a list of all DestinationType resources.
	DestinationTypes(ctx context.Context) ([]*model.DestinationType, error)
	// DestinationType returns a single DestinationType by name.
	DestinationType(ctx context.Context, name string) (*model.DestinationType, error)
	// DeleteDestinationType deletes a single Destination resource by name.
	DeleteDestinationType(ctx context.Context, name string) error

	// Apply upserts multiple resources of any kind.
	Apply(ctx context.Context, r []*model.AnyResource) ([]*model.AnyResourceStatus, error)
	// Delete deletes multiple resources, minimum required fields to delete are Kind and Metadata.Name.
	Delete(ctx context.Context, r []*model.AnyResource) ([]*model.AnyResourceStatus, error)

	// Version returns the version of the BindPlane-OP server.
	Version(ctx context.Context) (version.Version, error)

	// AgentInstallCommand returns the installation command for the given AgentInstallationOptions.
	AgentInstallCommand(ctx context.Context, options AgentInstallOptions) (string, error)
	// AgentUpgrade upgrades the agent with given ID to the specified version.
	AgentUpgrade(ctx context.Context, id string, version string) error

	// AgentLabels gets the labels for an agent
	AgentLabels(ctx context.Context, id string) (*model.Labels, error)
	// ApplyAgentLabels applies the specified labels to an agent, merging the specified labels with the existing labels
	// and returning the labels of the agent
	ApplyAgentLabels(ctx context.Context, id string, labels *model.Labels, override bool) (*model.Labels, error)
}

type bindplaneClient struct {
	client *resty.Client
	config *common.Client
	*zap.Logger
}

var _ BindPlane = (*bindplaneClient)(nil)

// NewBindPlane takes a client configuration, logger and returns a new BindPlane.
func NewBindPlane(config *common.Client, logger *zap.Logger) (BindPlane, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 20)
	client.SetBasicAuth(config.Username, config.Password)
	client.SetBaseURL(fmt.Sprintf("%s/v1", config.BindPlaneURL()))

	tlsConfig, err := tlsClient(config.Certificate, config.PrivateKey, config.CertificateAuthority, config.InsecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to configure TLS client: %w", err)
	}
	client.SetTLSClientConfig(tlsConfig)

	return &bindplaneClient{
		client: client,
		config: config,
		Logger: logger.Named("bindplane-client"),
	}, nil
}

func (c *bindplaneClient) Agents(ctx context.Context, options ...QueryOption) ([]*model.Agent, error) {
	c.Debug("Agents called")

	opts := makeQueryOptions(options)
	ar := &model.AgentsResponse{}
	resp, err := c.client.R().
		SetResult(ar).
		SetQueryParam("selector", opts.selector).
		SetQueryParam("query", opts.query).
		SetQueryParam("offset", fmt.Sprintf("%d", opts.offset)).
		SetQueryParam("limit", fmt.Sprintf("%d", opts.limit)).
		SetQueryParam("sort", opts.sort).
		Get("/agents")
	if err != nil {
		logRequestError(c.Logger, err, "/agents")
		return nil, err
	}

	return ar.Agents, c.statusError(resp, err, "unable to get agents")
}

func (c *bindplaneClient) Agent(ctx context.Context, id string) (*model.Agent, error) {
	c.Debug("Agent called")

	ar := &model.AgentResponse{}
	agentsEndpoint := fmt.Sprintf("/agents/%s", id)
	resp, err := c.client.R().SetResult(ar).Get(agentsEndpoint)
	if err != nil {
		logRequestError(c.Logger, err, agentsEndpoint)
		return nil, err
	}

	return ar.Agent, c.statusError(resp, err, "unable to get agents")
}

func (c *bindplaneClient) DeleteAgents(ctx context.Context, ids []string) ([]*model.Agent, error) {
	c.Debug("DeleteAgents called")

	body := &model.DeleteAgentsPayload{
		IDs: ids,
	}
	result := &model.DeleteAgentsResponse{}
	resp, err := c.client.R().SetBody(body).SetResult(result).Delete("/agents")
	return result.Agents, c.statusError(resp, err, "unable to delete agents")
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) AgentVersions(ctx context.Context) ([]*model.AgentVersion, error) {
	result := model.AgentVersionsResponse{}
	err := c.resources(ctx, "/agent-versions", &result)
	return result.AgentVersions, err
}

func (c *bindplaneClient) AgentVersion(ctx context.Context, name string) (*model.AgentVersion, error) {
	result := model.AgentVersionResponse{}
	err := c.resource(ctx, "/agent-versions", name, &result)
	return result.AgentVersion, err
}

func (c *bindplaneClient) DeleteAgentVersion(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/agent-versions", name)
}

func (c *bindplaneClient) SyncAgentVersions(ctx context.Context, version string) ([]*model.AnyResourceStatus, error) {
	ar := &model.ApplyResponseClientSide{}
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(ar).
		Post(fmt.Sprintf("/agent-versions/%s/sync", version))
	if err != nil {
		logRequestError(c.Logger, err, "/agent-versions/:name/sync")
		return nil, err
	}
	return ar.Updates, c.statusError(resp, err, "unable to sync agent-versions")
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) Configurations(ctx context.Context) ([]*model.Configuration, error) {
	c.Debug("Configurations called")

	pr := &model.ConfigurationsResponse{}
	resp, err := c.client.R().SetResult(pr).Get("/configurations")
	return pr.Configurations, c.statusError(resp, err, "unable to get configurations")
}

func (c *bindplaneClient) Configuration(ctx context.Context, name string) (*model.Configuration, error) {
	result := model.ConfigurationResponse{}
	err := c.resource(ctx, "/configurations", name, &result)
	return result.Configuration, err
}

func (c *bindplaneClient) DeleteConfiguration(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/configurations", name)
}

func (c *bindplaneClient) RawConfiguration(ctx context.Context, name string) (string, error) {
	result := model.ConfigurationResponse{}
	err := c.resource(ctx, "/configurations", name, &result)
	return result.Raw, err
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) Sources(ctx context.Context) ([]*model.Source, error) {
	result := model.SourcesResponse{}
	err := c.resources(ctx, "/sources", &result)
	return result.Sources, err
}

func (c *bindplaneClient) Source(ctx context.Context, name string) (*model.Source, error) {
	result := model.SourceResponse{}
	err := c.resource(ctx, "/sources", name, &result)
	return result.Source, err
}

func (c *bindplaneClient) DeleteSource(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/sources", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) SourceTypes(ctx context.Context) ([]*model.SourceType, error) {
	result := model.SourceTypesResponse{}
	err := c.resources(ctx, "/source-types", &result)
	return result.SourceTypes, err
}

func (c *bindplaneClient) SourceType(ctx context.Context, name string) (*model.SourceType, error) {
	result := model.SourceTypeResponse{}
	err := c.resource(ctx, "/source-types", name, &result)
	return result.SourceType, err
}

func (c *bindplaneClient) DeleteSourceType(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/source-types", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) Processors(ctx context.Context) ([]*model.Processor, error) {
	result := model.ProcessorsResponse{}
	err := c.resources(ctx, "/processors", &result)
	return result.Processors, err
}

func (c *bindplaneClient) Processor(ctx context.Context, name string) (*model.Processor, error) {
	result := model.ProcessorResponse{}
	err := c.resource(ctx, "/processors", name, &result)
	return result.Processor, err
}

func (c *bindplaneClient) DeleteProcessor(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/processors", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) ProcessorTypes(ctx context.Context) ([]*model.ProcessorType, error) {
	result := model.ProcessorTypesResponse{}
	err := c.resources(ctx, "/processor-types", &result)
	return result.ProcessorTypes, err
}

func (c *bindplaneClient) ProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	result := model.ProcessorTypeResponse{}
	err := c.resource(ctx, "/processor-types", name, &result)
	return result.ProcessorType, err
}

func (c *bindplaneClient) DeleteProcessorType(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/processor-types", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) Destinations(ctx context.Context) ([]*model.Destination, error) {
	result := model.DestinationsResponse{}
	err := c.resources(ctx, "/destinations", &result)
	return result.Destinations, err
}

func (c *bindplaneClient) Destination(ctx context.Context, name string) (*model.Destination, error) {
	result := model.DestinationResponse{}
	err := c.resource(ctx, "/destinations", name, &result)
	return result.Destination, err
}

func (c *bindplaneClient) DeleteDestination(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/destinations", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) DestinationTypes(ctx context.Context) ([]*model.DestinationType, error) {
	result := model.DestinationTypesResponse{}
	err := c.resources(ctx, "/destination-types", &result)
	return result.DestinationTypes, err
}

func (c *bindplaneClient) DestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	result := model.DestinationTypeResponse{}
	err := c.resource(ctx, "/destination-types", name, &result)
	return result.DestinationType, err
}

func (c *bindplaneClient) DeleteDestinationType(ctx context.Context, name string) error {
	return c.deleteResource(ctx, "/destination-types", name)
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) Apply(ctx context.Context, resources []*model.AnyResource) ([]*model.AnyResourceStatus, error) {
	c.Debug("Apply called")

	payload := model.ApplyPayload{
		Resources: resources,
	}

	data, err := jsoniter.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("client apply: %w", err)
	}

	ar := &model.ApplyResponseClientSide{}
	resp, err := c.client.R().SetHeader("Content-Type", "application/json").
		SetBody(data).SetResult(ar).Post("/apply")
	return ar.Updates, c.statusError(resp, err, "unable to apply resources")
}

func (c *bindplaneClient) Delete(ctx context.Context, resources []*model.AnyResource) ([]*model.AnyResourceStatus, error) {
	c.Debug("Batch Delete called")

	payload := model.DeletePayload{
		Resources: resources,
	}

	data, err := jsoniter.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data to json: %w", err)
	}

	resp, err := c.client.R().SetHeader("Content-Type", "application/json").
		SetBody(data).Post("/delete")
	if err != nil {
		logRequestError(c.Logger, err, "/delete")
		return nil, err
	}

	dr := &model.DeleteResponseClientSide{}

	switch resp.StatusCode() {
	case http.StatusAccepted:
		return dr.Updates, nil
	case http.StatusUnauthorized:
		return nil, c.unauthorizedError(resp)
	case http.StatusBadRequest:
		if dr.Errors != nil {
			return nil, errors.New(dr.Errors[0])
		}
		return nil, errors.New("bad request")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("%s", dr.Errors[0])
	}

	err = json.Unmarshal(resp.Body(), dr)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("unknown response from bindplane server")
}

func (c *bindplaneClient) Version(ctx context.Context) (version.Version, error) {
	c.Debug("Version called")

	v := version.Version{}
	resp, err := c.client.R().SetResult(&v).Get("/version")
	return v, c.statusError(resp, err, "unable to get version")
}

func (c *bindplaneClient) AgentInstallCommand(ctx context.Context, options AgentInstallOptions) (string, error) {
	c.Debug("AgentInstallCommand called")

	var command model.InstallCommandResponse
	endpoint := fmt.Sprintf("/agent-versions/%s/install-command", options.Version)

	resp, err := c.client.R().
		SetQueryParam("platform", options.Platform).
		SetQueryParam("version", options.Version).
		SetQueryParam("labels", options.Labels).
		SetQueryParam("remote-url", options.RemoteURL).
		SetQueryParam("secret-key", options.SecretKey).
		SetResult(&command).
		Get(endpoint)

	return command.Command, c.statusError(resp, err, "unable to get install command")
}

func (c *bindplaneClient) AgentUpgrade(ctx context.Context, id string, version string) error {
	endpoint := fmt.Sprintf("/agents/%s/version", id)
	resp, err := c.client.R().
		SetBody(model.PostAgentVersionRequest{
			Version: version,
		}).
		Post(endpoint)

	if err != nil {
		return err
	}

	// look for errors
	if resp.StatusCode() != http.StatusNoContent {
		errResponse := &model.ErrorResponse{}
		err = json.Unmarshal(resp.Body(), errResponse)
		if err != nil {
			return fmt.Errorf("unable to parse api response: %w", err)
		}

		if len(errResponse.Errors) > 0 {
			var errs error
			for _, e := range errResponse.Errors {
				errs = multierror.Append(errs, errors.New(e))
			}
			return errs
		}
	}

	return nil
}

func logRequestError(logger *zap.Logger, err error, endpoint string) {
	logger.Error("Error making request", zap.Error(err), zap.String("endpoint", endpoint))
}

func (c *bindplaneClient) AgentLabels(ctx context.Context, id string) (*model.Labels, error) {
	var response model.AgentLabelsResponse
	endpoint := fmt.Sprintf("/agents/%s/labels", id)

	resp, err := c.client.R().
		SetResult(&response).
		Get(endpoint)

	return response.Labels, c.statusError(resp, err, "unable to get agent labels")
}

func (c *bindplaneClient) ApplyAgentLabels(ctx context.Context, id string, labels *model.Labels, overwrite bool) (*model.Labels, error) {
	payload := model.AgentLabelsPayload{
		Labels: labels.AsMap(),
	}

	data, err := jsoniter.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal labels to apply: %w", err)
	}

	endpoint := fmt.Sprintf("/agents/%s/labels", id)
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("overwrite", strconv.FormatBool(overwrite)).
		SetBody(data).
		Patch(endpoint)

	if resp.StatusCode() != http.StatusConflict {
		err = c.statusError(resp, err, "unable to apply labels")
		if err != nil {
			return nil, err
		}
	}

	var response model.AgentLabelsResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse api response: %w", err)
	}

	if response.Errors != nil {
		err = fmt.Errorf(strings.Join(response.Errors, "\n"))
	}

	return response.Labels, err
}

// ----------------------------------------------------------------------

func (c *bindplaneClient) CopyConfig(ctx context.Context, name, copyName string) error {
	payload := model.PostCopyConfigRequest{
		Name: copyName,
	}

	endpoint := fmt.Sprintf("/configurations/%s/copy", name)

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(endpoint)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("a configuration with name '%s' already exists", copyName)
	default:
		err := &multierror.Error{}
		multierror.Append(err, fmt.Errorf("failed to copy configuration, got status %v", resp.StatusCode()))

		// check for errors field in response
		errResponse := &model.ErrorResponse{}
		if err := json.Unmarshal(resp.Body(), errResponse); err != nil {
			c.Logger.Error("failed to unmarshal error response when copying config", zap.Error(err))
		}

		if len(errResponse.Errors) > 0 {
			for _, e := range errResponse.Errors {
				multierror.Append(err, errors.New(e))
			}
		}

		return err.ErrorOrNil()
	}
}

// ----------------------------------------------------------------------

// resources gets the resources from the REST server and stores them in the provided result.
func (c *bindplaneClient) resources(ctx context.Context, resourcesURL string, result any) error {
	return c.get(ctx, resourcesURL, result)
}

// resource gets the resource with the specified name from the REST server and stores it in the provided result.
func (c *bindplaneClient) resource(ctx context.Context, resourcesURL string, name string, result any) error {
	resourceURL := fmt.Sprintf("%s/%s", resourcesURL, name)
	return c.get(ctx, resourceURL, result)
}

func (c *bindplaneClient) get(ctx context.Context, url string, result any) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(result).
		Get(url)

	if err != nil {
		logRequestError(c.Logger, err, url)
		return err
	}

	return c.statusError(resp, err, fmt.Sprintf("unable to get %s", url))
}

func (c *bindplaneClient) deleteResource(ctx context.Context, resourcesURL string, name string) error {
	deleteEndpoint := fmt.Sprintf("%s/%s", resourcesURL, name)
	resp, err := c.client.R().Delete(deleteEndpoint)
	if err != nil {
		logRequestError(c.Logger, err, deleteEndpoint)
		return fmt.Errorf("error making request to remote bindplane server, %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return c.unauthorizedError(resp)
	case http.StatusNotFound:
		return fmt.Errorf("%s not found", deleteEndpoint)
	case http.StatusBadRequest:
		dr := &model.DeleteResponse{}
		err = json.Unmarshal(resp.Body(), dr)
		if err != nil {
			return err
		}

		if dr.Errors != nil {
			return errors.New(dr.Errors[0])
		}

		return errors.New("bad request")
	case http.StatusConflict:
		errorResponse := &rest.ErrorResponse{}
		err = json.Unmarshal(resp.Body(), errorResponse)
		if err != nil {
			return errors.New("failed to parse response, status 409 Conflict")
		}

		if errorResponse.Errors != nil {
			return errors.New(errorResponse.Errors[0])
		}

		return errors.New("got status 409 Conflict")
	default:
		c.Logger.Error("unexpected status code received while trying to delete resource", zap.Int("statusCode", resp.StatusCode()), zap.String("endpoint", deleteEndpoint))
		return fmt.Errorf("unexpected status code received while trying to delete resource '%s': %s", name, resp.Status())
	}

}

func (c *bindplaneClient) unauthorizedError(resp *resty.Response) error {
	if resp.StatusCode() == http.StatusUnauthorized {
		err := fmt.Errorf(resp.Status())
		logRequestError(c.Logger, err, resp.Request.URL)
		return err
	}
	return nil
}

func (c *bindplaneClient) statusError(resp *resty.Response, err error, message string) error {
	if err != nil {
		logRequestError(c.Logger, err, resp.Request.URL)
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusCreated:
		return nil
	case http.StatusAccepted:
		return nil
	case http.StatusNoContent:
		return nil

	default:
		err := fmt.Errorf("%s, got %s", message, resp.Status())
		logRequestError(c.Logger, err, resp.Request.URL)
		return err
	}
}
