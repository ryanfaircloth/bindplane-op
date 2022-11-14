// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	eventbus "github.com/observiq/bindplane-op/internal/eventbus"
	mock "github.com/stretchr/testify/mock"

	model "github.com/observiq/bindplane-op/model"

	search "github.com/observiq/bindplane-op/internal/store/search"

	sessions "github.com/gorilla/sessions"

	stats "github.com/observiq/bindplane-op/internal/store/stats"

	store "github.com/observiq/bindplane-op/internal/store"

	time "time"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Agent provides a mock function with given fields: ctx, name
func (_m *Store) Agent(ctx context.Context, name string) (*model.Agent, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Agent
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Agent); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Agent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AgentConfiguration provides a mock function with given fields: ctx, agentID
func (_m *Store) AgentConfiguration(ctx context.Context, agentID string) (*model.Configuration, error) {
	ret := _m.Called(ctx, agentID)

	var r0 *model.Configuration
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Configuration); ok {
		r0 = rf(ctx, agentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Configuration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, agentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AgentIndex provides a mock function with given fields: ctx
func (_m *Store) AgentIndex(ctx context.Context) search.Index {
	ret := _m.Called(ctx)

	var r0 search.Index
	if rf, ok := ret.Get(0).(func(context.Context) search.Index); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(search.Index)
		}
	}

	return r0
}

// AgentVersion provides a mock function with given fields: ctx, name
func (_m *Store) AgentVersion(ctx context.Context, name string) (*model.AgentVersion, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.AgentVersion
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.AgentVersion); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AgentVersion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AgentVersions provides a mock function with given fields: ctx
func (_m *Store) AgentVersions(ctx context.Context) ([]*model.AgentVersion, error) {
	ret := _m.Called(ctx)

	var r0 []*model.AgentVersion
	if rf, ok := ret.Get(0).(func(context.Context) []*model.AgentVersion); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AgentVersion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Agents provides a mock function with given fields: ctx, options
func (_m *Store) Agents(ctx context.Context, options ...store.QueryOption) ([]*model.Agent, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*model.Agent
	if rf, ok := ret.Get(0).(func(context.Context, ...store.QueryOption) []*model.Agent); ok {
		r0 = rf(ctx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Agent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...store.QueryOption) error); ok {
		r1 = rf(ctx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AgentsCount provides a mock function with given fields: _a0, _a1
func (_m *Store) AgentsCount(_a0 context.Context, _a1 ...store.QueryOption) (int, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, ...store.QueryOption) int); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...store.QueryOption) error); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AgentsIDsMatchingConfiguration provides a mock function with given fields: ctx, conf
func (_m *Store) AgentsIDsMatchingConfiguration(ctx context.Context, conf *model.Configuration) ([]string, error) {
	ret := _m.Called(ctx, conf)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, *model.Configuration) []string); ok {
		r0 = rf(ctx, conf)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Configuration) error); ok {
		r1 = rf(ctx, conf)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ApplyResources provides a mock function with given fields: ctx, resources
func (_m *Store) ApplyResources(ctx context.Context, resources []model.Resource) ([]model.ResourceStatus, error) {
	ret := _m.Called(ctx, resources)

	var r0 []model.ResourceStatus
	if rf, ok := ret.Get(0).(func(context.Context, []model.Resource) []model.ResourceStatus); ok {
		r0 = rf(ctx, resources)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ResourceStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []model.Resource) error); ok {
		r1 = rf(ctx, resources)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CleanupDisconnectedAgents provides a mock function with given fields: ctx, since
func (_m *Store) CleanupDisconnectedAgents(ctx context.Context, since time.Time) error {
	ret := _m.Called(ctx, since)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) error); ok {
		r0 = rf(ctx, since)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Clear provides a mock function with given fields:
func (_m *Store) Clear() {
	_m.Called()
}

// Configuration provides a mock function with given fields: ctx, name
func (_m *Store) Configuration(ctx context.Context, name string) (*model.Configuration, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Configuration
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Configuration); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Configuration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConfigurationIndex provides a mock function with given fields: ctx
func (_m *Store) ConfigurationIndex(ctx context.Context) search.Index {
	ret := _m.Called(ctx)

	var r0 search.Index
	if rf, ok := ret.Get(0).(func(context.Context) search.Index); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(search.Index)
		}
	}

	return r0
}

// Configurations provides a mock function with given fields: ctx, options
func (_m *Store) Configurations(ctx context.Context, options ...store.QueryOption) ([]*model.Configuration, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*model.Configuration
	if rf, ok := ret.Get(0).(func(context.Context, ...store.QueryOption) []*model.Configuration); ok {
		r0 = rf(ctx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Configuration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...store.QueryOption) error); ok {
		r1 = rf(ctx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAgentVersion provides a mock function with given fields: ctx, name
func (_m *Store) DeleteAgentVersion(ctx context.Context, name string) (*model.AgentVersion, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.AgentVersion
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.AgentVersion); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AgentVersion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAgents provides a mock function with given fields: ctx, agentIDs
func (_m *Store) DeleteAgents(ctx context.Context, agentIDs []string) ([]*model.Agent, error) {
	ret := _m.Called(ctx, agentIDs)

	var r0 []*model.Agent
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*model.Agent); ok {
		r0 = rf(ctx, agentIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Agent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, agentIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteConfiguration provides a mock function with given fields: ctx, name
func (_m *Store) DeleteConfiguration(ctx context.Context, name string) (*model.Configuration, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Configuration
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Configuration); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Configuration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDestination provides a mock function with given fields: ctx, name
func (_m *Store) DeleteDestination(ctx context.Context, name string) (*model.Destination, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Destination
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Destination); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Destination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDestinationType provides a mock function with given fields: ctx, name
func (_m *Store) DeleteDestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.DestinationType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.DestinationType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DestinationType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProcessor provides a mock function with given fields: ctx, name
func (_m *Store) DeleteProcessor(ctx context.Context, name string) (*model.Processor, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Processor
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Processor); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Processor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProcessorType provides a mock function with given fields: ctx, name
func (_m *Store) DeleteProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.ProcessorType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.ProcessorType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProcessorType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteResources provides a mock function with given fields: ctx, resources
func (_m *Store) DeleteResources(ctx context.Context, resources []model.Resource) ([]model.ResourceStatus, error) {
	ret := _m.Called(ctx, resources)

	var r0 []model.ResourceStatus
	if rf, ok := ret.Get(0).(func(context.Context, []model.Resource) []model.ResourceStatus); ok {
		r0 = rf(ctx, resources)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ResourceStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []model.Resource) error); ok {
		r1 = rf(ctx, resources)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSource provides a mock function with given fields: ctx, name
func (_m *Store) DeleteSource(ctx context.Context, name string) (*model.Source, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Source
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Source); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Source)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSourceType provides a mock function with given fields: ctx, name
func (_m *Store) DeleteSourceType(ctx context.Context, name string) (*model.SourceType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.SourceType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SourceType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SourceType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destination provides a mock function with given fields: ctx, name
func (_m *Store) Destination(ctx context.Context, name string) (*model.Destination, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Destination
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Destination); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Destination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DestinationType provides a mock function with given fields: ctx, name
func (_m *Store) DestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.DestinationType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.DestinationType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DestinationType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DestinationTypes provides a mock function with given fields: ctx
func (_m *Store) DestinationTypes(ctx context.Context) ([]*model.DestinationType, error) {
	ret := _m.Called(ctx)

	var r0 []*model.DestinationType
	if rf, ok := ret.Get(0).(func(context.Context) []*model.DestinationType); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DestinationType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destinations provides a mock function with given fields: ctx
func (_m *Store) Destinations(ctx context.Context) ([]*model.Destination, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Destination
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Destination); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Destination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Measurements provides a mock function with given fields:
func (_m *Store) Measurements() stats.Measurements {
	ret := _m.Called()

	var r0 stats.Measurements
	if rf, ok := ret.Get(0).(func() stats.Measurements); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stats.Measurements)
		}
	}

	return r0
}

// Processor provides a mock function with given fields: ctx, name
func (_m *Store) Processor(ctx context.Context, name string) (*model.Processor, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Processor
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Processor); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Processor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessorType provides a mock function with given fields: ctx, name
func (_m *Store) ProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.ProcessorType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.ProcessorType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProcessorType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessorTypes provides a mock function with given fields: ctx
func (_m *Store) ProcessorTypes(ctx context.Context) ([]*model.ProcessorType, error) {
	ret := _m.Called(ctx)

	var r0 []*model.ProcessorType
	if rf, ok := ret.Get(0).(func(context.Context) []*model.ProcessorType); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ProcessorType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Processors provides a mock function with given fields: ctx
func (_m *Store) Processors(ctx context.Context) ([]*model.Processor, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Processor
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Processor); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Processor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Source provides a mock function with given fields: ctx, name
func (_m *Store) Source(ctx context.Context, name string) (*model.Source, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Source
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Source); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Source)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SourceType provides a mock function with given fields: ctx, name
func (_m *Store) SourceType(ctx context.Context, name string) (*model.SourceType, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.SourceType
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SourceType); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SourceType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SourceTypes provides a mock function with given fields: ctx
func (_m *Store) SourceTypes(ctx context.Context) ([]*model.SourceType, error) {
	ret := _m.Called(ctx)

	var r0 []*model.SourceType
	if rf, ok := ret.Get(0).(func(context.Context) []*model.SourceType); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.SourceType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Sources provides a mock function with given fields: ctx
func (_m *Store) Sources(ctx context.Context) ([]*model.Source, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Source
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Source); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Source)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Updates provides a mock function with given fields:
func (_m *Store) Updates() eventbus.Source[*store.Updates] {
	ret := _m.Called()

	var r0 eventbus.Source[*store.Updates]
	if rf, ok := ret.Get(0).(func() eventbus.Source[*store.Updates]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eventbus.Source[*store.Updates])
		}
	}

	return r0
}

// UpsertAgent provides a mock function with given fields: ctx, agentID, updater
func (_m *Store) UpsertAgent(ctx context.Context, agentID string, updater store.AgentUpdater) (*model.Agent, error) {
	ret := _m.Called(ctx, agentID, updater)

	var r0 *model.Agent
	if rf, ok := ret.Get(0).(func(context.Context, string, store.AgentUpdater) *model.Agent); ok {
		r0 = rf(ctx, agentID, updater)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Agent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, store.AgentUpdater) error); ok {
		r1 = rf(ctx, agentID, updater)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAgents provides a mock function with given fields: ctx, agentIDs, updater
func (_m *Store) UpsertAgents(ctx context.Context, agentIDs []string, updater store.AgentUpdater) ([]*model.Agent, error) {
	ret := _m.Called(ctx, agentIDs, updater)

	var r0 []*model.Agent
	if rf, ok := ret.Get(0).(func(context.Context, []string, store.AgentUpdater) []*model.Agent); ok {
		r0 = rf(ctx, agentIDs, updater)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Agent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, store.AgentUpdater) error); ok {
		r1 = rf(ctx, agentIDs, updater)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserSessions provides a mock function with given fields:
func (_m *Store) UserSessions() sessions.Store {
	ret := _m.Called()

	var r0 sessions.Store
	if rf, ok := ret.Get(0).(func() sessions.Store); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sessions.Store)
		}
	}

	return r0
}

type mockConstructorTestingTNewStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStore(t mockConstructorTestingTNewStore) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}