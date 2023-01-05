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
	"strings"
	"testing"

	"github.com/observiq/bindplane-op/model/otel"
	"github.com/observiq/bindplane-op/model/validation"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEvalCabinDestination(t *testing.T) {
	dt := fileResource[*DestinationType](t, "testfiles/destinationtype-cabin.yaml")
	d := fileResource[*Destination](t, "testfiles/destination-cabin.yaml")
	values := dt.evalOutput(&dt.Spec.Logs, d, func(e error) {
		require.NoError(t, e)
	})
	require.Len(t, values.Receivers, 0)
	require.Len(t, values.Processors, 1)
	require.Len(t, values.Exporters, 1)
	require.Len(t, values.Extensions, 0)

	// we expect observiq-cloud to be rendered
	_, ok := values.Exporters[0][otel.NewComponentID("observiq", "cabin-production-logs")]
	require.True(t, ok)

	exportersYaml, err := yaml.Marshal(values.Exporters)
	require.NoError(t, err)

	expectYaml := strings.TrimLeft(`
- observiq/cabin-production-logs:
    endpoint: https://nozzle.app.observiq.com
    secret_key: 2c088c5e-2afc-483b-be52-e2b657fcff08
    timeout: 10s
`, "\n")

	require.Equal(t, expectYaml, string(exportersYaml))
}

func TestEvalGoogleCloud(t *testing.T) {
	dt := fileResource[*DestinationType](t, "testfiles/destinationtype-googlecloud.yaml")
	d := fileResource[*Destination](t, "testfiles/destination-googlecloud.yaml")
	values := dt.eval(d, func(e error) {
		require.NoError(t, e)
	})
	require.Len(t, values[otel.Logs].Receivers, 0)
	require.Len(t, values[otel.Logs].Processors, 1)
	require.Len(t, values[otel.Logs].Exporters, 1)
	require.Len(t, values[otel.Logs].Extensions, 0)
	require.Len(t, values[otel.Metrics].Receivers, 0)
	require.Len(t, values[otel.Metrics].Processors, 1)
	require.Len(t, values[otel.Metrics].Exporters, 1)
	require.Len(t, values[otel.Metrics].Extensions, 0)
	require.Len(t, values[otel.Traces].Receivers, 0)
	require.Len(t, values[otel.Traces].Processors, 1)
	require.Len(t, values[otel.Traces].Exporters, 1)
	require.Len(t, values[otel.Traces].Extensions, 0)
}

func TestValidateNoDuplicateParamterNames(t *testing.T) {
	testCases := []struct {
		description string
		spec        ResourceTypeSpec
		expectErr   bool
		errorMsg    string
	}{
		{
			"duplicate names, error",
			ResourceTypeSpec{
				Parameters: []ParameterDefinition{
					{
						Name: "param-1",
					},
					{
						Name: "param-1",
					},
					{
						Name: "param-2",
					},
				},
			},
			true,
			"1 error occurred:\n\t* found multiple parameters with name param-1\n\n",
		},
		{
			"no duplicates, ok",
			ResourceTypeSpec{
				Parameters: []ParameterDefinition{

					{
						Name: "param-1",
					},
					{
						Name: "param-2",
					},
				},
			},
			false,
			"1 error occurred:\n\t* found multiple parameters with name param-1\n\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			errs := validation.NewErrors()
			test.spec.validateNoDuplicateParameterNames(errs)
			if test.expectErr {
				require.Error(t, errs.Result())
				require.Equal(t, test.errorMsg, errs.Result().Error())
			} else {
				require.NoError(t, errs.Result())
			}

		})
	}
}

func TestTelemetryTypes(t *testing.T) {
	macosSourceType := fileResource[*SourceType](t, "testfiles/sourcetype-macos.yaml")
	otlpSourceType := fileResource[*SourceType](t, "testfiles/sourcetype-otlp.yaml")

	tests := []struct {
		description string
		sourceType  *SourceType
		expect      []otel.PipelineType
	}{
		{
			description: "macos supports logs and metrics",
			sourceType:  macosSourceType,
			expect:      []otel.PipelineType{otel.Logs, otel.Metrics},
		},
		{
			description: "otlp supports logs, metrics, and traces",
			sourceType:  otlpSourceType,
			expect:      []otel.PipelineType{otel.Logs, otel.Metrics, otel.Traces},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.sourceType.Spec.TelemetryTypes()
			require.ElementsMatch(t, test.expect, got)
		},
		)
	}
}

func TestResourceType_templateFuncHasMetricsEnabled(t *testing.T) {
	type fields struct {
		ResourceMeta ResourceMeta
		Spec         ResourceTypeSpec
	}
	type args struct {
		parameterValue []any
		parameterName  string
		category       string
	}
	tests := []struct {
		name    string
		metrics ParameterDefinition
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "missing definition",
			metrics: ParameterDefinition{
				Name: "missing",
				Type: metricsType,
			},
			args: args{
				parameterName: "something-else",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "bad type",
			metrics: ParameterDefinition{
				Name: "wrongType",
				Type: stringType,
			},
			args: args{
				parameterName: "wrongType",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "no metrics",
			metrics: ParameterDefinition{
				Name: "metrics",
				Type: metricsType,
				Options: ParameterOptions{
					MetricCategories: []MetricCategory{
						{
							Metrics: nil,
						},
					},
				},
			},
			args: args{
				parameterName: "metrics",
			},
			want: false,
		},
		{
			name: "no disabled metrics",
			metrics: ParameterDefinition{
				Name: "metrics",
				Type: metricsType,
				Options: ParameterOptions{
					MetricCategories: []MetricCategory{
						{
							Label: "Network",
							Metrics: []MetricOption{
								{
									Name: "system.network.io",
								},
								{
									Name: "system.network.errors",
								},
							},
						},
					},
				},
			},
			args: args{
				parameterValue: nil,
				parameterName:  "metrics",
				category:       "Network",
			},
			want: true,
		},
		{
			name: "all disabled metrics",
			metrics: ParameterDefinition{
				Name: "metrics",
				Type: metricsType,
				Options: ParameterOptions{
					MetricCategories: []MetricCategory{
						{
							Label: "Network",
							Metrics: []MetricOption{
								{
									Name: "system.network.io",
								},
								{
									Name: "system.network.errors",
								},
							},
						},
					},
				},
			},
			args: args{
				parameterValue: []any{"system.network.io", "system.network.errors"},
				parameterName:  "metrics",
				category:       "Network",
			},
			want: false,
		},
		{
			name: "some disabled metrics",
			metrics: ParameterDefinition{
				Name: "metrics",
				Type: metricsType,
				Options: ParameterOptions{
					MetricCategories: []MetricCategory{
						{
							Label: "Network",
							Metrics: []MetricOption{
								{
									Name: "system.network.io",
								},
								{
									Name: "system.network.errors",
								},
							},
						},
					},
				},
			},
			args: args{
				parameterValue: []any{"system.network.io"},
				parameterName:  "metrics",
				category:       "Network",
			},
			want: true,
		},
		{
			name: "all disabled metrics",
			metrics: ParameterDefinition{
				Name: "metrics",
				Type: metricsType,
				Options: ParameterOptions{
					MetricCategories: []MetricCategory{
						{
							Label: "Network",
							Metrics: []MetricOption{
								{
									Name: "system.network.io",
								},
								{
									Name: "system.network.errors",
								},
							},
						},
					},
				},
			},
			args: args{
				parameterValue: []any{"system.network.io", "system.network.errors"},
				parameterName:  "metrics",
				category:       "Not a real category",
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rt := &ResourceType{
				ResourceMeta: ResourceMeta{
					Metadata: Metadata{
						Name: "test",
					},
				},
				Spec: ResourceTypeSpec{
					Parameters: []ParameterDefinition{
						test.metrics,
					},
				},
			}
			got, err := rt.templateFuncHasCategoryMetricsEnabled(test.args.parameterValue, test.args.parameterName, test.args.category)
			if (err != nil) != test.wantErr {
				t.Errorf("ResourceType.templateFuncHasMetricsEnabled() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("ResourceType.templateFuncHasMetricsEnabled() = %v, want %v", got, test.want)
			}
		})
	}
}
