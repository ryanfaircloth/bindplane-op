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

package model

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func validateResource[T Resource](t *testing.T, name string) T {
	return fileResource[T](t, filepath.Join("testfiles", "validate", name))
}
func testResource[T Resource](t *testing.T, name string) T {
	return fileResource[T](t, filepath.Join("testfiles", name))
}
func fileResource[T Resource](t *testing.T, path string) T {
	resources, err := ResourcesFromFile(path)
	require.NoError(t, err)

	parsed, err := ParseResources(resources)
	require.NoError(t, err)
	require.Len(t, parsed, 1)

	resource, ok := parsed[0].(T)
	require.True(t, ok)
	return resource
}

type testConfiguration struct {
	bindplaneURL string
}

func newTestConfiguration() BindPlaneConfiguration {
	return &testConfiguration{}
}

func (c *testConfiguration) BindPlaneURL() string {
	return c.bindplaneURL
}

var _ BindPlaneConfiguration = (*testConfiguration)(nil)

type testResourceStore struct {
	sources          map[string]*Source
	sourceTypes      map[string]*SourceType
	processors       map[string]*Processor
	processorTypes   map[string]*ProcessorType
	destinations     map[string]*Destination
	destinationTypes map[string]*DestinationType
}

func newTestResourceStore() *testResourceStore {
	return &testResourceStore{
		sources:          map[string]*Source{},
		sourceTypes:      map[string]*SourceType{},
		processors:       map[string]*Processor{},
		processorTypes:   map[string]*ProcessorType{},
		destinations:     map[string]*Destination{},
		destinationTypes: map[string]*DestinationType{},
	}
}

var _ ResourceStore = (*testResourceStore)(nil)

func (s *testResourceStore) Source(ctx context.Context, name string) (*Source, error) {
	return s.sources[name], nil
}
func (s *testResourceStore) SourceType(ctx context.Context, name string) (*SourceType, error) {
	return s.sourceTypes[name], nil
}
func (s *testResourceStore) Processor(ctx context.Context, name string) (*Processor, error) {
	return s.processors[name], nil
}
func (s *testResourceStore) ProcessorType(ctx context.Context, name string) (*ProcessorType, error) {
	return s.processorTypes[name], nil
}
func (s *testResourceStore) Destination(ctx context.Context, name string) (*Destination, error) {
	return s.destinations[name], nil
}
func (s *testResourceStore) DestinationType(ctx context.Context, name string) (*DestinationType, error) {
	return s.destinationTypes[name], nil
}

func TestParseConfiguration(t *testing.T) {
	path := filepath.Join("testfiles", "configuration-raw.yaml")
	bytes, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read the testfile")
	var configuration Configuration
	err = yaml.Unmarshal(bytes, &configuration)
	require.NoError(t, err)
	require.Equal(t, "cabin-production-configuration", configuration.Metadata.Name)
	require.Equal(t, "receivers:", strings.Split(configuration.Spec.Raw, "\n")[0])
}

func TestEvalConfiguration(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	macos := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[macos.Name()] = macos

	cabin := testResource[*Destination](t, "destination-cabin.yaml")
	store.destinations[cabin.Name()] = cabin

	cabinType := testResource[*DestinationType](t, "destinationtype-cabin.yaml")
	store.destinationTypes[cabinType.Name()] = cabinType

	configuration := testResource[*Configuration](t, "configuration-macos-sources.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
    plugin/source1__journald:
        plugin:
            name: journald
    plugin/source1__macos:
        parameters:
            - name: enable_system_log
              value: true
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/cabin-production-logs: null
exporters:
    observiq/cabin-production-logs:
        endpoint: https://nozzle.app.observiq.com
        secret_key: 2c088c5e-2afc-483b-be52-e2b657fcff08
        timeout: 10s
service:
    pipelines:
        logs/source0__cabin-production-logs:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - batch/cabin-production-logs
            exporters:
                - observiq/cabin-production-logs
        logs/source1__cabin-production-logs:
            receivers:
                - plugin/source1__macos
                - plugin/source1__journald
            processors:
                - batch/cabin-production-logs
            exporters:
                - observiq/cabin-production-logs
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration2(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	macos := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[macos.Name()] = macos

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	configuration := testResource[*Configuration](t, "configuration-macos-googlecloud.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    hostmetrics/source1:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
    plugin/source1__journald:
        plugin:
            name: journald
    plugin/source1__macos:
        parameters:
            - name: enable_system_log
              value: true
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/destination0: null
    normalizesums/destination0: null
exporters:
    googlecloud/destination0: null
service:
    pipelines:
        logs/source0__destination0:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - batch/destination0
            exporters:
                - googlecloud/destination0
        logs/source1__destination0:
            receivers:
                - plugin/source1__macos
                - plugin/source1__journald
            processors:
                - batch/destination0
            exporters:
                - googlecloud/destination0
        metrics/source0__destination0:
            receivers:
                - hostmetrics/source0
            processors:
                - normalizesums/destination0
                - batch/destination0
            exporters:
                - googlecloud/destination0
        metrics/source1__destination0:
            receivers:
                - hostmetrics/source1
            processors:
                - normalizesums/destination0
                - batch/destination0
            exporters:
                - googlecloud/destination0
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration3(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	otlp := testResource[*SourceType](t, "sourcetype-otlp.yaml")
	store.sourceTypes[otlp.Name()] = otlp

	googleCloudType := testResource[*DestinationType](t, "destinationtype-otlp.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	configuration := testResource[*Configuration](t, "configuration-otlp.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    otlp/source0:
        protocols:
            grpc: null
            http: null
processors:
    batch/destination0: null
exporters:
    otlp/destination0:
        endpoint: otelcol:4317
service:
    pipelines:
        logs/source0__destination0:
            receivers:
                - otlp/source0
            processors:
                - batch/destination0
            exporters:
                - otlp/destination0
        metrics/source0__destination0:
            receivers:
                - otlp/source0
            processors:
                - batch/destination0
            exporters:
                - otlp/destination0
        traces/source0__destination0:
            receivers:
                - otlp/source0
            processors:
                - batch/destination0
            exporters:
                - otlp/destination0
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration4(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-postgresql.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	configuration := testResource[*Configuration](t, "configuration-postgresql-googlecloud.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0__postgresql:
        parameters:
            postgresql_log_path:
                - /var/log/postgresql/postgresql*.log
                - /var/lib/pgsql/data/log/postgresql*.log
                - /var/lib/pgsql/*/data/log/postgresql*.log
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/postgresql_logs.yaml
processors:
    batch/destination0: null
exporters:
    googlecloud/destination0: null
service:
    pipelines:
        logs/source0__destination0:
            receivers:
                - plugin/source0__postgresql
            processors:
                - batch/destination0
            exporters:
                - googlecloud/destination0
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration5(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType

	configuration := testResource[*Configuration](t, "configuration-macos-processors.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/googlecloud: null
    normalizesums/googlecloud: null
    resourceattributetransposer/source0__processor0:
        operations:
            - from: from.attribute
              to: to.attribute
    resourceattributetransposer/source0__processor1:
        operations:
            - from: from.attribute2
              to: to.attribute2
exporters:
    googlecloud/googlecloud: null
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
        metrics/source0__googlecloud:
            receivers:
                - hostmetrics/source0
            processors:
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - normalizesums/googlecloud
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfigurationDestinationProcessors(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType

	configuration := testResource[*Configuration](t, "configuration-macos-destination-processors.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/googlecloud: null
    normalizesums/googlecloud: null
    resourceattributetransposer/googlecloud__processor0:
        operations:
            - from: from.attribute3
              to: to.attribute3
    resourceattributetransposer/googlecloud__processor1:
        operations:
            - from: from.attribute4
              to: to.attribute4
    resourceattributetransposer/source0__processor0:
        operations:
            - from: from.attribute
              to: to.attribute
    resourceattributetransposer/source0__processor1:
        operations:
            - from: from.attribute2
              to: to.attribute2
exporters:
    googlecloud/googlecloud: null
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - resourceattributetransposer/googlecloud__processor0
                - resourceattributetransposer/googlecloud__processor1
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
        metrics/source0__googlecloud:
            receivers:
                - hostmetrics/source0
            processors:
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - resourceattributetransposer/googlecloud__processor0
                - resourceattributetransposer/googlecloud__processor1
                - normalizesums/googlecloud
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfigurationDestinationProcessorsWithMeasurements(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType

	agent := &Agent{
		ID:      "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Version: v1_9_2.String(),
	}

	configuration := testResource[*Configuration](t, "configuration-macos-destination-processors.yaml")
	result, err := configuration.Render(context.TODO(), agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: 01ARZ3NDEKTSV4RRFFQ69G5FAV
                        configuration: macos-xy
                      targets:
                        - 0.0.0.0:8888
processors:
    batch/_agent_metrics: null
    batch/googlecloud: null
    normalizesums/googlecloud: null
    resourceattributetransposer/googlecloud__processor0:
        operations:
            - from: from.attribute3
              to: to.attribute3
    resourceattributetransposer/googlecloud__processor1:
        operations:
            - from: from.attribute4
              to: to.attribute4
    resourceattributetransposer/source0__processor0:
        operations:
            - from: from.attribute
              to: to.attribute
    resourceattributetransposer/source0__processor1:
        operations:
            - from: from.attribute2
              to: to.attribute2
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d0_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_metrics_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_metrics_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - throughputmeasurement/_s0_logs_source0
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_googlecloud
                - resourceattributetransposer/googlecloud__processor0
                - resourceattributetransposer/googlecloud__processor1
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
        metrics/source0__googlecloud:
            receivers:
                - hostmetrics/source0
            processors:
                - throughputmeasurement/_s0_metrics_source0
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - throughputmeasurement/_s1_metrics_source0
                - throughputmeasurement/_d0_metrics_googlecloud
                - resourceattributetransposer/googlecloud__processor0
                - resourceattributetransposer/googlecloud__processor1
                - throughputmeasurement/_d1_metrics_googlecloud
                - normalizesums/googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfigurationMultiDestination(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	cabinType := testResource[*DestinationType](t, "destinationtype-cabin.yaml")
	store.destinationTypes[cabinType.Name()] = cabinType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	cabin := testResource[*Destination](t, "destination-cabin.yaml")
	store.destinations[cabin.Name()] = cabin

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType
	agent := &Agent{
		ID:      "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Version: v1_9_2.String(),
	}

	configuration := testResource[*Configuration](t, "configuration-macos-multi-destination.yaml")
	result, err := configuration.Render(context.TODO(), agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: 01ARZ3NDEKTSV4RRFFQ69G5FAV
                        configuration: macos-xy
                      targets:
                        - 0.0.0.0:8888
processors:
    batch/_agent_metrics: null
    batch/cabin-production-logs: null
    batch/googlecloud: null
    normalizesums/googlecloud: null
    resourceattributetransposer/source0__processor0:
        operations:
            - from: from.attribute
              to: to.attribute
    resourceattributetransposer/source0__processor1:
        operations:
            - from: from.attribute2
              to: to.attribute2
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_cabin-production-logs:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d0_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_cabin-production-logs:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_metrics_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_metrics_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    observiq/cabin-production-logs:
        endpoint: https://nozzle.app.observiq.com
        secret_key: 2c088c5e-2afc-483b-be52-e2b657fcff08
        timeout: 10s
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__cabin-production-logs:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - throughputmeasurement/_s0_logs_source0
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_cabin-production-logs
                - throughputmeasurement/_d1_logs_cabin-production-logs
                - batch/cabin-production-logs
                - snapshotprocessor
            exporters:
                - observiq/cabin-production-logs
        logs/source0__googlecloud:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - throughputmeasurement/_d0_logs_googlecloud
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
        metrics/source0__googlecloud:
            receivers:
                - hostmetrics/source0
            processors:
                - throughputmeasurement/_s0_metrics_source0
                - resourceattributetransposer/source0__processor0
                - resourceattributetransposer/source0__processor1
                - throughputmeasurement/_s1_metrics_source0
                - throughputmeasurement/_d0_metrics_googlecloud
                - throughputmeasurement/_d1_metrics_googlecloud
                - normalizesums/googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfigurationFailsMissingResource(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType

	configuration := testResource[*Configuration](t, "configuration-macos-processors.yaml")

	tests := []struct {
		name            string
		deleteResources func()
		expectError     string
		expect          string
	}{
		{
			name:            "deletes sourceType",
			deleteResources: func() { delete(store.sourceTypes, postgresql.Name()) },
			expectError:     "1 error occurred:\n\t* unknown SourceType: MacOS\n\n",
		},
		{
			name:            "deletes googleCloudType",
			deleteResources: func() { delete(store.destinationTypes, googleCloudType.Name()) },
			expectError:     "1 error occurred:\n\t* unknown DestinationType: googlecloud\n\n",
		},
		{
			name:            "deletes destination",
			deleteResources: func() { delete(store.destinations, googleCloud.Name()) },
			expectError:     "1 error occurred:\n\t* unknown Destination: googlecloud\n\n",
		},
		{
			name:            "deletes processorType",
			deleteResources: func() { delete(store.processorTypes, resourceAttributeTransposerType.Name()) },
			expectError:     "2 errors occurred:\n\t* unknown ProcessorType: resource-attribute-transposer\n\t* unknown ProcessorType: resource-attribute-transposer\n\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// before rendering, delete resources that we reference
			test.deleteResources()

			_, err := configuration.Render(context.TODO(), nil, config, store)
			require.Error(t, err)
			require.Equal(t, test.expectError, err.Error())

			// reset for next iteration
			store.sourceTypes[postgresql.Name()] = postgresql
			store.destinationTypes[googleCloudType.Name()] = googleCloudType
			store.destinations[googleCloud.Name()] = googleCloud
			store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType
		})
	}
}

func TestConfigurationRender_DisabledDestination(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	macos := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[macos.Name()] = macos

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	cabinType := testResource[*DestinationType](t, "destinationtype-cabin.yaml")
	store.destinationTypes[cabinType.Name()] = cabinType

	cabin := testResource[*Destination](t, "destination-cabin.yaml")
	store.destinations[cabin.Name()] = cabin

	configuration := testResource[*Configuration](t, "configuration-macos-googlecloud-disabled.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	// We expect the full pipeline, omitting the disabled googlecloud destination
	expect := strings.TrimLeft(`
receivers:
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
    plugin/source1__journald:
        plugin:
            name: journald
    plugin/source1__macos:
        parameters:
            - name: enable_system_log
              value: true
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/cabin-production-logs: null
exporters:
    observiq/cabin-production-logs:
        endpoint: https://nozzle.app.observiq.com
        secret_key: 2c088c5e-2afc-483b-be52-e2b657fcff08
        timeout: 10s
service:
    pipelines:
        logs/source0__cabin-production-logs:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - batch/cabin-production-logs
            exporters:
                - observiq/cabin-production-logs
        logs/source1__cabin-production-logs:
            receivers:
                - plugin/source1__macos
                - plugin/source1__journald
            processors:
                - batch/cabin-production-logs
            exporters:
                - observiq/cabin-production-logs
`, "\n")
	require.Equal(t, expect, result)
}
func TestConfigurationRender_DisabledSource(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	macos := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[macos.Name()] = macos

	fileLog := testResource[*SourceType](t, "sourcetype-filelog.yaml")
	store.sourceTypes[fileLog.Name()] = fileLog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	configuration := testResource[*Configuration](t, "configuration-macos-source-disabled.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	// We expect the full pipeline, omitting the disabled macOS source
	expect := strings.TrimLeft(`
receivers:
    plugin/source1:
        parameters:
            encoding: utf-8
            file_path:
                - /foo/bar/baz.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: none
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
processors:
    batch/googlecloud: null
    resourcedetection/source1:
        detectors:
            - system
        system:
            hostname_sources:
                - os
exporters:
    googlecloud/googlecloud: null
service:
    pipelines:
        logs/source1__googlecloud:
            receivers:
                - plugin/source1
            processors:
                - resourcedetection/source1
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
`, "\n")
	require.Equal(t, expect, result)
}

func TestEvalConfiguration_WithLogCount(t *testing.T) {
	t.Parallel()
	store := newTestResourceStore()
	config := newTestConfiguration()

	filelog := testResource[*SourceType](t, "sourcetype-filelog.yaml")
	store.sourceTypes[filelog.Name()] = filelog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	logCountProcessor := testResource[*ProcessorType](t, "processortype-countlogs.yaml")
	store.processorTypes[logCountProcessor.Name()] = logCountProcessor

	configuration := testResource[*Configuration](t, "configuration-file-count-logs.yaml")
	agent := Agent{
		Version: "v1.14.0",
	}
	result, err := configuration.Render(context.TODO(), &agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0:
        parameters:
            encoding: utf-8
            file_path:
                - /tmp/test.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: json
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: ""
                        configuration: file-countlogs
                      targets:
                        - 0.0.0.0:8888
    route/builtin: null
processors:
    batch/_agent_metrics: null
    batch/googlecloud: null
    logcount/source0__processor0:
        attributes:
            status_code: body.status
        interval: 60s
        match: "true"
        metric_name: custom.metric.count
        metric_unit: '{logs}'
        route: builtin
    normalizesums/googlecloud: null
    resourcedetection/source0:
        detectors:
            - system
        system:
            hostname_sources:
                - os
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d0_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0
            processors:
                - resourcedetection/source0
                - throughputmeasurement/_s0_logs_source0
                - logcount/source0__processor0
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_googlecloud
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
        metrics/route__googlecloud:
            receivers:
                - route/builtin
            processors:
                - throughputmeasurement/_d0_metrics_googlecloud
                - throughputmeasurement/_d1_metrics_googlecloud
                - normalizesums/googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestConfigurationRender_DisabledProcessor(t *testing.T) {
	store := newTestResourceStore()
	config := newTestConfiguration()

	postgresql := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[postgresql.Name()] = postgresql

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	resourceAttributeTransposerType := testResource[*ProcessorType](t, "processortype-resourceattributetransposer.yaml")
	store.processorTypes[resourceAttributeTransposerType.Name()] = resourceAttributeTransposerType

	configuration := testResource[*Configuration](t, "configuration-macos-processors-disabled.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    hostmetrics/source0:
        collection_interval: 1m
        scrapers:
            load: null
    plugin/source0__journald:
        plugin:
            name: journald
    plugin/source0__macos:
        parameters:
            - name: enable_system_log
              value: false
            - name: system_log_path
              value: /var/log/system.log
            - name: enable_install_log
              value: true
            - name: install_log_path
              value: /var/log/install.log
            - name: start_at
              value: end
        plugin:
            name: macos
processors:
    batch/googlecloud: null
    normalizesums/googlecloud: null
exporters:
    googlecloud/googlecloud: null
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0__macos
                - plugin/source0__journald
            processors:
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
        metrics/source0__googlecloud:
            receivers:
                - hostmetrics/source0
            processors:
                - normalizesums/googlecloud
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration_WithLogCountUnsupported(t *testing.T) {
	t.Parallel()
	store := newTestResourceStore()
	config := newTestConfiguration()

	filelog := testResource[*SourceType](t, "sourcetype-filelog.yaml")
	store.sourceTypes[filelog.Name()] = filelog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	logCountProcessor := testResource[*ProcessorType](t, "processortype-countlogs.yaml")
	store.processorTypes[logCountProcessor.Name()] = logCountProcessor

	configuration := testResource[*Configuration](t, "configuration-file-count-logs.yaml")
	agent := Agent{
		Version: "v1.13.0",
	}
	result, err := configuration.Render(context.TODO(), &agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0:
        parameters:
            encoding: utf-8
            file_path:
                - /tmp/test.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: json
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: ""
                        configuration: file-countlogs
                      targets:
                        - 0.0.0.0:8888
processors:
    batch/_agent_metrics: null
    batch/googlecloud: null
    logcount/source0__processor0:
        attributes:
            status_code: body.status
        interval: 60s
        match: "true"
        metric_name: custom.metric.count
        metric_unit: '{logs}'
        route: builtin
    resourcedetection/source0:
        detectors:
            - system
        system:
            hostname_sources:
                - os
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0
            processors:
                - resourcedetection/source0
                - throughputmeasurement/_s0_logs_source0
                - logcount/source0__processor0
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_googlecloud
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration_WithMetricExtract(t *testing.T) {
	t.Parallel()
	store := newTestResourceStore()
	config := newTestConfiguration()

	filelog := testResource[*SourceType](t, "sourcetype-filelog.yaml")
	store.sourceTypes[filelog.Name()] = filelog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	extractMetricProcessor := testResource[*ProcessorType](t, "processortype-extractmetric.yaml")
	store.processorTypes[extractMetricProcessor.Name()] = extractMetricProcessor

	configuration := testResource[*Configuration](t, "configuration-file-extract-metric.yaml")
	agent := Agent{
		Version: "v1.14.0",
	}
	result, err := configuration.Render(context.TODO(), &agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0:
        parameters:
            encoding: utf-8
            file_path:
                - /var/log/http.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: json
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: ""
                        configuration: file-extract-duration
                      targets:
                        - 0.0.0.0:8888
    route/builtin: null
processors:
    batch/_agent_metrics: null
    batch/googlecloud: null
    metricextract/source0__processor0:
        attributes:
            status_code: body.status
        extract: body.duration
        match: body.duration != nil
        metric_name: http.request.duration
        metric_unit: ms
        route: builtin
    normalizesums/googlecloud: null
    resourcedetection/source0:
        detectors:
            - system
        system:
            hostname_sources:
                - os
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d0_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_metrics_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0
            processors:
                - resourcedetection/source0
                - throughputmeasurement/_s0_logs_source0
                - metricextract/source0__processor0
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_googlecloud
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
        metrics/route__googlecloud:
            receivers:
                - route/builtin
            processors:
                - throughputmeasurement/_d0_metrics_googlecloud
                - throughputmeasurement/_d1_metrics_googlecloud
                - normalizesums/googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}

func TestEvalConfiguration_WithMetricExtractUnsupported(t *testing.T) {
	t.Parallel()
	store := newTestResourceStore()
	config := newTestConfiguration()

	filelog := testResource[*SourceType](t, "sourcetype-filelog.yaml")
	store.sourceTypes[filelog.Name()] = filelog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	extractMetricProcessor := testResource[*ProcessorType](t, "processortype-extractmetric.yaml")
	store.processorTypes[extractMetricProcessor.Name()] = extractMetricProcessor

	configuration := testResource[*Configuration](t, "configuration-file-extract-metric.yaml")
	agent := Agent{
		Version: "v1.13.22",
	}
	result, err := configuration.Render(context.TODO(), &agent, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0:
        parameters:
            encoding: utf-8
            file_path:
                - /var/log/http.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: json
            start_at: end
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
    prometheus/_agent_metrics:
        config:
            scrape_configs:
                - job_name: observiq-otel-collector
                  metric_relabel_configs:
                    - action: keep
                      regex: otelcol_processor_throughputmeasurement_.*
                      source_labels:
                        - __name__
                  scrape_interval: 10s
                  static_configs:
                    - labels:
                        agent: ""
                        configuration: file-extract-duration
                      targets:
                        - 0.0.0.0:8888
processors:
    batch/_agent_metrics: null
    batch/googlecloud: null
    metricextract/source0__processor0:
        attributes:
            status_code: body.status
        extract: body.duration
        match: body.duration != nil
        metric_name: http.request.duration
        metric_unit: ms
        route: builtin
    resourcedetection/source0:
        detectors:
            - system
        system:
            hostname_sources:
                - os
    snapshotprocessor: null
    throughputmeasurement/_d0_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_d1_logs_googlecloud:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s0_logs_source0:
        enabled: true
        sampling_ratio: 1
    throughputmeasurement/_s1_logs_source0:
        enabled: true
        sampling_ratio: 1
exporters:
    googlecloud/googlecloud: null
    otlphttp/_agent_metrics:
        endpoint: /v1/otlphttp
service:
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0
            processors:
                - resourcedetection/source0
                - throughputmeasurement/_s0_logs_source0
                - metricextract/source0__processor0
                - throughputmeasurement/_s1_logs_source0
                - throughputmeasurement/_d0_logs_googlecloud
                - throughputmeasurement/_d1_logs_googlecloud
                - batch/googlecloud
                - snapshotprocessor
            exporters:
                - googlecloud/googlecloud
        metrics/_agent_metrics:
            receivers:
                - prometheus/_agent_metrics
            processors:
                - batch/_agent_metrics
            exporters:
                - otlphttp/_agent_metrics
`, "\n")

	require.Equal(t, expect, result)
}

func TestDuplicate(t *testing.T) {
	duplicateName := "duplicate-config"

	configuration := testResource[*Configuration](t, "configuration-macos-googlecloud.yaml")
	require.NotNil(t, configuration)

	new := configuration.Duplicate(duplicateName)
	require.NotNil(t, new)

	t.Run("equal sources, destinations", func(t *testing.T) {
		require.Equal(t, configuration.Spec.Sources, new.Spec.Sources)
		require.Equal(t, configuration.Spec.Destinations, new.Spec.Destinations)
	})

	t.Run("replace name, id, and match labels", func(t *testing.T) {
		// Set the duplicate name
		require.Equal(t, new.Metadata.Name, duplicateName)

		// Set a new ID
		require.NotEqual(t, new.Metadata.ID, configuration.Metadata.ID)

		// Set the configuration matchLabel
		require.Contains(t, new.Spec.Selector.MatchLabels, "configuration")
		require.Equal(t, new.Spec.Selector.MatchLabels["configuration"], duplicateName)
	})
}

func TestConfigurationType(t *testing.T) {
	t.Run("raw configuration", func(t *testing.T) {
		c := NewConfigurationWithSpec("raw-config", ConfigurationSpec{
			Raw: "my raw configuration",
		})
		require.Equal(t, ConfigurationTypeRaw, c.Type())
	})
	t.Run("managed configuration", func(t *testing.T) {
		c := NewConfigurationWithSpec("managed-config", ConfigurationSpec{
			Sources: []ResourceConfiguration{
				{Name: "a source"},
			},
		})
		require.Equal(t, ConfigurationTypeModular, c.Type())
	})
}

func TestEvalConfiguration_FileLogStorage(t *testing.T) {
	t.Parallel()
	store := newTestResourceStore()
	config := newTestConfiguration()

	macos := testResource[*SourceType](t, "sourcetype-macos.yaml")
	store.sourceTypes[macos.Name()] = macos

	filelog := testResource[*SourceType](t, "sourcetype-filelog-storage.yaml")
	store.sourceTypes[filelog.Name()] = filelog

	googleCloudType := testResource[*DestinationType](t, "destinationtype-googlecloud.yaml")
	store.destinationTypes[googleCloudType.Name()] = googleCloudType

	googleCloud := testResource[*Destination](t, "destination-googlecloud.yaml")
	store.destinations[googleCloud.Name()] = googleCloud

	configuration := testResource[*Configuration](t, "configuration-filelog-storage.yaml")
	result, err := configuration.Render(context.TODO(), nil, config, store)
	require.NoError(t, err)

	expect := strings.TrimLeft(`
receivers:
    plugin/source0:
        parameters:
            encoding: utf-8
            file_path:
                - /foo/bar/baz.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: none
            start_at: end
            storage: file_storage/source0
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
    plugin/source1:
        parameters:
            encoding: utf-8
            file_path:
                - /foo/bar/baz2.log
            log_type: file
            multiline_line_start_pattern: ""
            parse_format: none
            start_at: end
            storage: file_storage/source1
        path: $OIQ_OTEL_COLLECTOR_HOME/plugins/file_logs.yaml
processors:
    batch/googlecloud: null
    resourcedetection/source0:
        detectors:
            - system
        system:
            hostname_sources:
                - os
    resourcedetection/source1:
        detectors:
            - system
        system:
            hostname_sources:
                - os
exporters:
    googlecloud/googlecloud: null
extensions:
    file_storage/source0:
        directory: /tmp/offset_storage_dir
    file_storage/source1:
        directory: /tmp/offset_storage_dir
service:
    extensions:
        - file_storage/source0
        - file_storage/source1
    pipelines:
        logs/source0__googlecloud:
            receivers:
                - plugin/source0
            processors:
                - resourcedetection/source0
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
        logs/source1__googlecloud:
            receivers:
                - plugin/source1
            processors:
                - resourcedetection/source1
                - batch/googlecloud
            exporters:
                - googlecloud/googlecloud
`, "\n")

	require.Equal(t, expect, result)
}
