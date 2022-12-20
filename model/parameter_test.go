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
	"encoding/json"
	"strconv"
	"testing"

	"github.com/observiq/bindplane-op/model/validation"
	"github.com/stretchr/testify/require"
)

func TestValidateDefault(t *testing.T) {
	testCases := []struct {
		name      string
		expectErr bool
		param     ParameterDefinition
	}{
		{
			"ValidStringDefault",
			false,
			ParameterDefinition{
				Type:    "string",
				Default: "test",
			},
		},
		{
			"InvalidStringDefault",
			true,
			ParameterDefinition{
				Type:    "string",
				Default: 5,
			},
		},
		{
			"ValidIntDefault",
			false,
			ParameterDefinition{
				Type:    "int",
				Default: 5,
			},
		},
		{
			"InvalidStringDefault",
			true,
			ParameterDefinition{
				Type:    "int",
				Default: "test",
			},
		},
		{
			"ValidBoolDefault",
			false,
			ParameterDefinition{
				Type:    "bool",
				Default: true,
			},
		},
		{
			"InvalidBoolDefault",
			true,
			ParameterDefinition{
				Type:    "bool",
				Default: "test",
			},
		},
		{
			"ValidStringsDefault",
			false,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{"test"},
			},
		},
		{
			"InvalidStringsDefault",
			true,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{5},
			},
		},
		{
			"ValidEnumDefault",
			false,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     "test",
			},
		},
		{
			"InvalidEnumDefault",
			true,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     "invalid",
			},
		},
		{
			"ValidEnumsDefaultEmpty",
			false,
			ParameterDefinition{
				Type:        "enums",
				ValidValues: []string{"foo", "bar", "baz", "blah"},
				Default:     []any{},
			},
		},
		{
			"ValidEnumsDefaultAll",
			false,
			ParameterDefinition{
				Type:        "enums",
				ValidValues: []string{"foo", "bar", "baz", "blah"},
				Default:     []any{"foo", "bar", "baz", "blah"},
			},
		},
		{
			"NonStringEnumDefault",
			true,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     5,
			},
		},
		{
			"InvalidTypeDefault",
			true,
			ParameterDefinition{
				Type:    "float",
				Default: 5,
			},
		},
		{
			"MetricsTypeNoDefault",
			true,
			ParameterDefinition{
				Type: "metrics",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.param.validateDefault()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateOptions(t *testing.T) {
	testCases := []struct {
		name      string
		expectErr bool
		errorMsg  string
		param     ParameterDefinition
	}{
		{
			"Enum Creatable OK",
			false,
			"",
			ParameterDefinition{
				Type: "enum",
				Options: ParameterOptions{
					Creatable: true,
				},
			},
		},
		{
			"Enum Not Creatable OK",
			false,
			"",
			ParameterDefinition{
				Type: "enum",
				Options: ParameterOptions{
					Creatable: false,
				},
			},
		},
		{
			"Non-Enum Creatable Error",
			true,
			"1 error occurred:\n\t* creatable is true for parameter of type 'string'\n\n",
			ParameterDefinition{
				Type: "string",
				Options: ParameterOptions{
					Creatable: true,
				},
			},
		},
		{
			"Non-Enum Not Creatable OK",
			false,
			"",
			ParameterDefinition{
				Type: "enum",
				Options: ParameterOptions{
					Creatable: true,
				},
			},
		},
		{
			"Enums TrackUnchecked OK",
			false,
			"",
			ParameterDefinition{
				Type: "enums",
				Options: ParameterOptions{
					TrackUnchecked: true,
				},
			},
		},
		{
			"Enums !TrackUnchecked OK",
			false,
			"",
			ParameterDefinition{
				Type: "enums",
			},
		},
		{
			"Non Enums TrackUnchecked Error",
			true,
			"1 error occurred:\n\t* trackUnchecked is true for parameter of type `map`\n\n",
			ParameterDefinition{
				Type: "map",
				Options: ParameterOptions{
					TrackUnchecked: true,
				},
			},
		},
		{
			"Non MapToSubForm TrackUnchecked Error",
			true,
			"1 error occurred:\n\t* trackUnchecked is true for parameter of type `map`\n\n",
			ParameterDefinition{
				Type: "map",
				Options: ParameterOptions{
					TrackUnchecked: true,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			errs := validation.NewErrors()
			test.param.validateOptions(errs)
			if test.expectErr {
				require.Error(t, errs.Result())
				require.Equal(t, test.errorMsg, errs.Result().Error())
			} else {
				require.NoError(t, errs.Result())
			}
		})
	}
}

func TestValidateValue(t *testing.T) {
	testCases := []struct {
		name      string
		expectErr bool
		param     ParameterDefinition
		value     interface{}
	}{
		{
			"ValidString",
			false,
			ParameterDefinition{
				Type:    "string",
				Default: "test",
			},
			"string",
		},
		{
			"InvalidString",
			true,
			ParameterDefinition{
				Type:    "string",
				Default: "test",
			},
			5,
		},
		{
			"ValidInt",
			false,
			ParameterDefinition{
				Type:    "int",
				Default: 5,
			},
			5,
		},
		{
			"InvalidInt",
			true,
			ParameterDefinition{
				Type:    "int",
				Default: 5,
			},
			"test",
		},
		{
			"Valid json.Number",
			false,
			ParameterDefinition{
				Type:    "int",
				Default: 5,
			},
			json.Number(strconv.FormatInt(60, 10)),
		},
		{
			"Invalid json.Number",
			true,
			ParameterDefinition{
				Type:    "int",
				Default: 5,
			},
			json.Number(strconv.FormatFloat(60.104824, 'e', -1, 64)),
		},
		{
			"ValidBool",
			false,
			ParameterDefinition{
				Type:    "bool",
				Default: true,
			},
			false,
		},
		{
			"InvalidBool",
			true,
			ParameterDefinition{
				Type:    "bool",
				Default: true,
			},
			"test",
		},
		{
			"ValidStringsAsInterface",
			false,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{"test"},
			},
			[]interface{}{"test"},
		},
		{
			"ValidStrings",
			false,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{"test"},
			},
			[]string{"test"},
		},
		{
			"InvalidStringsAsInterface",
			true,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{"test"},
			},
			[]interface{}{5},
		},
		{
			"InvalidStrings",
			true,
			ParameterDefinition{
				Type:    "strings",
				Default: []interface{}{"test"},
			},
			[]int{5},
		},
		{
			"ValidEnum",
			false,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     "test",
			},
			"test",
		},
		{
			"ValidEnum - Creatable",
			false,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Options: ParameterOptions{
					Creatable: true,
				},
			},
			"not-test",
		},
		{
			"InvalidEnumValue",
			true,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     "test",
			},
			"missing",
		},
		{
			"InvalidEnumtype",
			true,
			ParameterDefinition{
				Type:        "enum",
				ValidValues: []string{"test"},
				Default:     "test",
			},
			5,
		},
		{
			"InvalidType",
			true,
			ParameterDefinition{
				Type:    "float",
				Default: 5,
			},
			5,
		},
		{
			"ValidMap",
			false,
			ParameterDefinition{
				Type: "map",
			},
			map[string]string{
				"foo":  "bar",
				"blah": "baz",
			},
		},
		{
			"InvalidMap",
			true,
			ParameterDefinition{
				Type: "map",
			},
			5,
		},
		{
			"InvalidMapType",
			true,
			ParameterDefinition{
				Type: "map",
			},
			map[string]interface{}{
				"blah": 1,
				"foo":  "blah",
			},
		},
		{
			"ValidYaml",
			false,
			ParameterDefinition{
				Type: "yaml",
			},
			`blah: foo
bar: baz
baz:
- one
- two
`,
		}, {
			"ValidYaml",
			false,
			ParameterDefinition{
				Type: "yaml",
			},
			`- one
- two
`,
		},
		{
			"InvalidYaml",
			true,
			ParameterDefinition{
				Type: "yaml",
			},
			`one: two
three: four
- five: 6
seven:
	- eight
	- nine
	- 10
eleven:
	- twelve: thirteen
	fourteen: 15
`,
		},
		{
			"InvalidYaml",
			true,
			ParameterDefinition{
				Type: "yaml",
			},
			`{{{}}}`,
		},
		{
			"ValidEnums",
			false,
			ParameterDefinition{
				Type:        "enums",
				ValidValues: []string{"one", "two", "three", "four"},
			},
			[]any{
				"two", "four",
			},
		},
		{
			"InvalidEnums",
			true,
			ParameterDefinition{
				Type:        "enums",
				ValidValues: []string{"one", "two", "three", "four"},
			},
			[]any{
				"one", "seven",
			},
		},
		{
			"ValidTimezone",
			false,
			ParameterDefinition{
				Type: "timezone",
			},
			"UTC",
		},
		{
			"InvalidTimezone",
			true,
			ParameterDefinition{
				Type: "timezone",
			},
			10,
		},
		{
			"InvalidTimezoneValue",
			true,
			ParameterDefinition{
				Type: "timezone",
			},
			"America/NewJersey",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set the name for better error messages
			tc.param.Name = tc.name
			err := tc.param.validateValue(tc.value)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

var kpi = true
var boolPtr = &kpi
var desc = "description"
var descPtr = &desc

func TestValidateMetricCategory(t *testing.T) {
	testCases := []struct {
		description    string
		metricCategory MetricCategory
		expectErr      bool
		errorMsg       string
	}{
		{
			"No label, bad",
			MetricCategory{
				Metrics: []MetricOption{
					{Name: "one.two.three"},
				},
			},
			true,
			"1 error occurred:\n\t* missing required field Label in metric category\n\n",
		},
		{
			"No metrics, bad",
			MetricCategory{
				Label: "Blah",
			},
			true,
			"1 error occurred:\n\t* missing required field metrics on metricCategory\n\n",
		},
		{
			"Column not zero or one",
			MetricCategory{
				Label:   "Blah",
				Metrics: []MetricOption{{Name: "Something"}},
				Column:  3,
			},
			true,
			"1 error occurred:\n\t* metric category value is neither 0 nor 1\n\n",
		},
		{
			"OK",
			MetricCategory{
				Label:   "Blah",
				Column:  1,
				Metrics: []MetricOption{{Name: "First"}, {Name: "Seconds"}},
			},
			false,
			"",
		},
		{
			"Metrics missing required name field",
			MetricCategory{
				Label:   "Foo",
				Column:  0,
				Metrics: []MetricOption{{}, {KPI: boolPtr}, {Description: descPtr}},
			},
			true,
			"3 errors occurred:\n\t* missing required name field for metric option\n\t* missing required name field for metric option\n\t* missing required name field for metric option\n\n",
		},
	}

	for _, test := range testCases {
		errs := validation.NewErrors()
		test.metricCategory.validateMetricCategory(errs)
		if test.expectErr {
			require.Error(t, errs.Result())
			require.Equal(t, test.errorMsg, errs.Result().Error())
		} else {
			require.NoError(t, errs.Result())
		}
	}
}

func TestValidateAwsCloudwatchNamedFieldType(t *testing.T) {
	tests := []struct {
		description   string
		parameterType ParameterDefinition
		value         any
		wantErr       string
	}{
		{
			"bad value type int",
			ParameterDefinition{
				Type: awsCloudwatchNamedFieldType,
			},
			4,
			"malformed value for parameter of type awsCloudwatchNamedField",
		},
		{
			"malformed slice item 'prefixes'",
			ParameterDefinition{
				Type: awsCloudwatchNamedFieldType,
			},
			[]struct {
				ID       string `mapstructure:"id"`
				Names    []interface{}
				Prefixes string
			}{
				{
					"id",
					[]interface{}{"one", "two", "three"},
					"foo",
				},
			},
			"incorrect type included in Prefixes field",
		},
		{
			"malformed slice item 'id'",
			ParameterDefinition{
				Type: awsCloudwatchNamedFieldType,
			},
			[]struct {
				ID       int
				Names    []interface{}
				Prefixes []interface{}
			}{
				{
					4,
					[]interface{}{"one", "two", "three"},
					[]interface{}{"one", "two", "three"},
				},
			},
			"incorrect type included in 'id' field",
		},
		{
			"malformed slice item 'names'",
			ParameterDefinition{
				Type: awsCloudwatchNamedFieldType,
			},
			[]struct {
				ID       string `mapstructure:"id"`
				Names    string
				Prefixes []interface{}
			}{
				{
					"id",
					"blah",
					[]interface{}{"one", "two", "three"},
				},
			},
			"incorrect type included in Names field",
		},
		{
			"ok value",
			ParameterDefinition{
				Type: awsCloudwatchNamedFieldType,
			},
			[]struct {
				ID       string
				Names    []interface{}
				Prefixes []interface{}
			}{
				{
					"id",
					make([]interface{}, 0),
					make([]interface{}, 0),
				},
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			gotErr := test.parameterType.validateAwsCloudwatchNamedFieldType(awsCloudwatchNamedFieldType, test.value)
			if test.wantErr != "" {
				require.Equal(t, test.wantErr, gotErr.Error())
			} else {
				require.NoError(t, gotErr)
			}

		})
	}
}
