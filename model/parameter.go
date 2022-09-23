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
	"fmt"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/observiq/bindplane-op/model/validation"
	"github.com/observiq/stanza/errors"
	"gopkg.in/yaml.v3"
)

const (
	stringType   = "string"
	boolType     = "bool"
	intType      = "int"
	stringsType  = "strings"
	enumType     = "enum"
	enumsType    = "enums"
	yamlType     = "yaml"
	mapType      = "map"
	timezoneType = "timezone"
	metricsType  = "metrics"
)

// ParameterDefinition is a basic description of a definition's parameter. This implementation comes directly from
// stanza plugin parameters with slight modifications for mapstructure.
type ParameterDefinition struct {
	Name        string `json:"name" yaml:"name"`
	Label       string `json:"label" yaml:"label"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`

	// "string", "int", "bool", "strings", or "enum"
	Type string `json:"type" yaml:"type"`

	// only useable if Type == "enum"
	ValidValues []string `json:"validValues,omitempty" yaml:"validValues,omitempty" mapstructure:"validValues"`

	// Must be valid according to Type & ValidValues
	Default        interface{}           `json:"default,omitempty" yaml:"default,omitempty"`
	RelevantIf     []RelevantIfCondition `json:"relevantIf,omitempty" yaml:"relevantIf,omitempty" mapstructure:"relevantIf"`
	Hidden         bool                  `json:"hidden" yaml:"hidden"`
	AdvancedConfig bool                  `json:"advancedConfig" yaml:"advancedConfig" mapstructure:"advancedConfig"`

	Options ParameterOptions `json:"options" yaml:"options"`

	Documentation []DocumentationLink `json:"documentation,omitempty" yaml:"documentation,omitempty"`
}

// DocumentationLink contains the text and url for documentation of a ParameterDefinition
type DocumentationLink struct {
	Text string `json:"text" yaml:"text"`
	URL  string `json:"url" yaml:"url"`
}

// ParameterOptions specify further customization for input parameters
type ParameterOptions struct {
	// Creatable will modify the "enum" parameter from a select to
	// a creatable select where a user can specify a custom value
	Creatable bool `json:"creatable" yaml:"creatable"`

	// TrackUnchecked will modify the "enums" parameter to store the
	// unchecked values as the value.
	TrackUnchecked bool `json:"trackUnchecked" yaml:"trackUnchecked"`

	// GridColumns will specify the number of flex-grid columns the
	// control will take up, must be an integer between 1 and 12 or
	// unspecified.
	GridColumns *int `json:"gridColumns,omitempty" yaml:"gridColumns,omitempty"`

	// SectionHeader is used to indicate that the bool parameter input is
	// a switch for further configuration for UI styling.
	SectionHeader *bool `json:"sectionHeader,omitempty" yaml:"sectionHeader,omitempty"`

	MetricCategories []MetricCategory `json:"metricCategories" yaml:"metricCategories"`

	// Multiline indicates that a multiline textarea should be used for editing a "string" parameter.
	Multiline bool `json:"multiline,omitempty" yaml:"multiline,omitempty"`
}

// MetricCategory consists of the label, optional column, and metrics for a metricsType Parameter
type MetricCategory struct {
	Label string `json:"label" yaml:"label"`
	// 0 or 1
	Column int `json:"column" yaml:"column"`

	Metrics []MetricOption `json:"metrics" yaml:"metrics"`
}

// MetricOption is an individual metric that can be specified for a metricsType Parameter
type MetricOption struct {
	Name        string  `json:"name" yaml:"name"`
	Description *string `json:"description" yaml:"description"`
	KPI         *bool   `json:"kpi" yaml:"kpi"`
}

// RelevantIfCondition specifies a condition under which a parameter is deemed relevant.
type RelevantIfCondition struct {
	Name     string `json:"name" yaml:"name" mapstructure:"name"`
	Operator string `json:"operator" yaml:"operator" mapstructure:"operator"`
	Value    any    `json:"value" yaml:"value" mapstructure:"value"`
}

func (p ParameterDefinition) validateValue(value interface{}) error {
	return p.validateValueType(parameterFieldValue, value)
}

func (p ParameterDefinition) validateDefinition(kind Kind, errs validation.Errors) {
	if err := p.validateName(); err != nil {
		errs.Add(err)
	}

	if err := p.validateType(); err != nil {
		errs.Add(err)
	}

	if err := p.validateValidValues(); err != nil {
		errs.Add(err)
	}

	if err := p.validateDefault(); err != nil {
		errs.Add(err)
	}

	p.validateOptions(errs)

	p.validateSpecialParameters(kind, errs)
}

// validateSpecialParameters ensures that for consistency, common parameters like start_at appear the same in all sources
func (p ParameterDefinition) validateSpecialParameters(kind Kind, errs validation.Errors) {
	if kind == KindSourceType {
		switch p.Name {
		case "start_at":
			p.validateSpecialParameter(errs, ParameterDefinition{
				Name:           "start_at",
				Label:          "Start At",
				Description:    "Start reading logs from 'beginning' or 'end'.",
				Type:           "enum",
				ValidValues:    []string{"beginning", "end"},
				Default:        "end",
				AdvancedConfig: true,
			})
		case "collection_interval":
			// special case for vmware_vcenter which needs a longer collection interval
			if p.Description != "How often (minutes) to scrape for metrics." {
				p.validateSpecialParameter(errs, ParameterDefinition{
					Name:           "collection_interval",
					Label:          "Collection Interval",
					Description:    "How often (seconds) to scrape for metrics.",
					Type:           "int",
					Default:        60,
					AdvancedConfig: true,
				})
			}
		case "jar_path":
			p.validateSpecialParameter(errs, ParameterDefinition{
				Label:          "JMX Metrics Collection Jar Path",
				Description:    "Full path to the JMX metrics jar.",
				Type:           "string",
				Default:        "/opt/opentelemetry-java-contrib-jmx-metrics.jar",
				AdvancedConfig: true,
				RelevantIf: []RelevantIfCondition{
					{
						Name:     "enable_metrics",
						Operator: "equals",
						Value:    true,
					},
				},
			})
		}

		// use full width for paths
		if p.Name != "jar_path" && (strings.HasSuffix(p.Name, "_path") || strings.HasSuffix(p.Name, "_paths")) {
			if p.Options.GridColumns == nil || *p.Options.GridColumns != 12 {
				errs.Warn(errors.NewError(fmt.Sprintf("%s parameter appears to be a path and should use the full width. ", p.Name), "specify gridColumns: 12 in options"))
			}
		}
	}
}

func (p ParameterDefinition) validateSpecialParameter(errs validation.Errors, expect ParameterDefinition) {
	// for consistency, %s should be the same anywhere it appears in sources
	if p.Label != expect.Label {
		errs.Warn(fmt.Errorf("%s parameter with label: %s should use label: %s", p.Name, p.Label, expect.Label))
	}
	if p.Description != expect.Description {
		errs.Warn(fmt.Errorf("%s parameter with description: %s should use description: %s", p.Name, p.Description, expect.Description))
	}
	if p.Type != expect.Type {
		errs.Warn(fmt.Errorf("%s parameter with type: %s should have type: %s", p.Name, p.Type, expect.Type))
	}
	pValidValues := strings.Join(p.ValidValues, ",")
	eValidValues := strings.Join(expect.ValidValues, ",")
	if pValidValues != eValidValues {
		errs.Warn(fmt.Errorf("%s parameter with validValues: [%s] should have validValues: [%s]", p.Name, pValidValues, eValidValues))
	}
	pDefault := fmt.Sprintf("%v", p.Default)
	eDefault := fmt.Sprintf("%v", expect.Default)
	if expect.Default != nil && pDefault != eDefault {
		errs.Warn(fmt.Errorf("%s parameter with default: %s should have default: %s", p.Name, pDefault, eDefault))
	}
	if p.AdvancedConfig != expect.AdvancedConfig {
		errs.Warn(fmt.Errorf("%s parameter with advancedConfig: %t should have advancedConfig: %t", p.Name, p.AdvancedConfig, expect.AdvancedConfig))
	}
}

func (p ParameterDefinition) validateName() error {
	if p.Name == "" {
		return errors.NewError(
			"missing name for parameter",
			"ensure that the name is a valid go identifier that can be used in go templates",
		)
	}
	if !token.IsIdentifier(p.Name) {
		return errors.NewError(
			fmt.Sprintf("invalid name '%s' for parameter", p.Name),
			"ensure that the name is a valid go identifier that can be used in go templates",
		)
	}
	return nil
}

func (p ParameterDefinition) validateType() error {
	if p.Type == "" {
		return errors.NewError(
			fmt.Sprintf("missing type for '%s'", p.Name),
			"ensure that the type is one of 'string', 'int', 'bool', 'strings', or 'enum'",
		)
	}
	switch p.Type {
	case stringType, intType, boolType, stringsType, enumType, enumsType, mapType, yamlType, timezoneType, metricsType: // ok
	default:
		return errors.NewError(
			fmt.Sprintf("invalid type '%s' for '%s'", p.Type, p.Name),
			"ensure that the type is one of 'string', 'int', 'bool', 'strings', or 'enum'",
		)
	}
	return nil
}

func (p ParameterDefinition) validateOptions(errs validation.Errors) {
	if p.Options.Creatable && p.Type != "enum" {
		errs.Add(
			errors.NewError(
				fmt.Sprintf("creatable is true for parameter of type '%s'", p.Type),
				"remove 'creatable' field or change type to 'enum'",
			))
	}

	if p.Options.TrackUnchecked && p.Type != "enums" {
		errs.Add(
			errors.NewError(
				fmt.Sprintf("trackUnchecked is true for parameter of type `%s`", p.Type),
				"remove 'trackUnchecked' field or change type to 'enums`",
			),
		)
	}

	if p.Options.Multiline && p.Type != "string" {
		errs.Add(
			errors.NewError(
				fmt.Sprintf("multiline is true for parameter of type `%s`", p.Type),
				"remove 'multiline' field or change type to 'string`",
			),
		)
	}

	p.validateMetricCategories(errs)
}

func (p ParameterDefinition) validateMetricCategories(errs validation.Errors) {
	switch p.Type {
	case metricsType:
		if p.Options.MetricCategories == nil {
			errs.Add(
				errors.NewError("options.metricCategories is required for type metrics",
					"include a metricCategories field under options or change the type from 'metrics'",
				),
			)
		}

		for _, category := range p.Options.MetricCategories {
			category.validateMetricCategory(errs)
		}

	default:
		if p.Options.MetricCategories != nil {
			errs.Add(
				errors.NewError(fmt.Sprintf("options.metricCategories is not a valid option for type '%s'", p.Type),
					"remove metricCategories field under options or change the type to 'metrics'",
				),
			)
		}

	}

}

func (m *MetricCategory) validateMetricCategory(errs validation.Errors) {
	if m.Label == "" {
		errs.Add(
			errors.NewError(
				"missing required field Label in metric category",
				"make sure all metric categories contain a label field",
			))
	}

	if m.Column != 0 && m.Column != 1 {
		errs.Add(
			errors.NewError(
				"metric category value is neither 0 nor 1",
				"make sure metric category column field is either 0 or 1",
			))
	}

	if m.Metrics == nil || len(m.Metrics) == 0 {
		errs.Add(
			errors.NewError(
				"missing required field metrics on metricCategory",
				"add a an array of MetricOptions to the metricCategory",
			))
	}

	for _, metricOption := range m.Metrics {
		if metricOption.Name == "" {
			errs.Add(
				errors.NewError(
					"missing required name field for metric option",
					"add a name field for each metric option in a metric category",
				))
		}
	}
}

func (p ParameterDefinition) validateValidValues() error {
	switch p.Type {
	case stringType, intType, boolType, stringsType, yamlType, mapType:
		if len(p.ValidValues) > 0 {
			return errors.NewError(
				fmt.Sprintf("validValues is undefined for parameter of type '%s'", p.Type),
				"remove 'validValues' field or change type to 'enum' or 'enums",
			)
		}
	case enumType, enumsType:
		if len(p.ValidValues) == 0 {
			return errors.NewError(
				"parameter of type 'enum' or 'enums' must have 'validValues' specified",
				"specify an array that includes one or more valid values",
			)
		}
	}
	return nil
}

func (p ParameterDefinition) validateDefault() error {
	switch {
	case p.Type == metricsType && p.Default == nil:
		return errors.NewError(
			"default is required for parameter type 'metrics'",
			"set the default value to an empty array",
		)
	case p.Default == nil:
		return nil
	default:
		return p.validateValueType(parameterFieldDefault, p.Default)
	}
}

type parameterFieldType string

const (
	parameterFieldValue      parameterFieldType = "parameter"
	parameterFieldDefault                       = "default"
	parameterFieldRelevantIf                    = "relevantIf"
)

// validateValueType determines if the specified value is of the right type.
func (p ParameterDefinition) validateValueType(fieldType parameterFieldType, value any) error {
	switch p.Type {
	case stringType:
		return p.validateStringValue(fieldType, value)
	case intType:
		return p.validateIntValue(fieldType, value)
	case boolType:
		return p.validateBoolValue(fieldType, value)
	case stringsType:
		return p.validateStringArrayValue(fieldType, value)
	case enumType:
		return p.validateEnumValue(fieldType, value)
	case enumsType:
		return p.validateEnumsValue(fieldType, value)
	case mapType:
		return p.validateMapValue(fieldType, value)
	case yamlType:
		return p.validateYamlValue(fieldType, value)
	case timezoneType:
		return p.validateTimezoneType(fieldType, value)
	case metricsType:
		return p.validateMetricsType(fieldType, value)
	default:
		return errors.NewError(
			"invalid type for parameter",
			"ensure that the type is one of 'string', 'int', 'bool', 'strings', or 'enum'",
		)
	}
}

func (p ParameterDefinition) validateStringValue(fieldType parameterFieldType, value any) error {
	if _, ok := value.(string); !ok {
		return errors.NewError(
			fmt.Sprintf("%s value for '%s' must be a string", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is a string", fieldType),
		)
	}
	return nil
}

func (p ParameterDefinition) validateIntValue(fieldType parameterFieldType, value any) error {
	isIntValue := false

	if _, ok := value.(int); ok {
		// obvious case of integer
		isIntValue = true
	} else if floatValue, ok := value.(float64); ok {
		// less obvious case of float64
		if floatValue == float64(int(floatValue)) {
			isIntValue = true
		}
	} else if stringValue, ok := value.(string); ok {
		_, err := strconv.Atoi(stringValue)
		isIntValue = err == nil
	}

	if !isIntValue {
		return errors.NewError(
			fmt.Sprintf("%s value for '%s' must be an integer", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is an integer", fieldType),
		)
	}
	return nil
}

func (p ParameterDefinition) validateBoolValue(fieldType parameterFieldType, value any) error {
	isBoolValue := false

	if _, ok := value.(bool); ok {
		isBoolValue = true
	} else if stringValue, ok := value.(string); ok {
		_, err := strconv.ParseBool(stringValue)
		isBoolValue = err == nil
	}

	if !isBoolValue {
		return errors.NewError(
			fmt.Sprintf("%s value for '%s' must be a bool", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is a bool", fieldType),
		)
	}
	return nil
}

func (p ParameterDefinition) validateStringArrayValue(fieldType parameterFieldType, value any) error {
	if _, ok := value.([]string); ok {
		return nil
	}
	valueList, ok := value.([]interface{})
	if !ok {
		return errors.NewError(
			fmt.Sprintf("%s value for '%s' must be an array of strings", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is an array of string", fieldType),
		)
	}
	for _, s := range valueList {
		if _, ok := s.(string); !ok {
			return errors.NewError(
				fmt.Sprintf("%s value for '%s' must be an array of strings", fieldType, p.Name),
				fmt.Sprintf("ensure that the %s value is an array of string", fieldType),
			)
		}
	}
	return nil
}

func (p ParameterDefinition) validateEnumValue(fieldType parameterFieldType, value any) error {
	def, ok := value.(string)
	if !ok {
		return errors.NewError(
			fmt.Sprintf("%s value for enumerated parameter '%s'.", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is a string", fieldType),
		)
	}

	// If the enum is creatable thats all we need to check - any string value is valid.
	if p.Options.Creatable {
		return nil
	}

	for _, val := range p.ValidValues {
		if val == def {
			return nil
		}
	}
	return errors.NewError(
		fmt.Sprintf("%s value for '%s' must be one of %v", fieldType, p.Name, p.ValidValues),
		fmt.Sprintf("ensure %s value is listed as a valid value", fieldType),
	)
}

func (p ParameterDefinition) validateEnumsValue(fieldType parameterFieldType, value any) error {
	def, ok := value.([]any)
	if !ok {
		return errors.NewError(
			fmt.Sprintf("%s value for enums parameter '%s'", fieldType, p.Name),
			fmt.Sprintf("ensure that the %s value is a string array", fieldType),
		)
	}

	// Make sure all strings in the value are a validValue
	err := &multierror.Error{}
	for _, str := range def {
		var containsValue bool
		for _, validValue := range p.ValidValues {
			if str == validValue {
				containsValue = true
				break
			}
		}

		if !containsValue {
			multierror.Append(err,
				errors.NewError(
					fmt.Sprintf("%s value for '%s' must be one of %v", fieldType, p.Name, p.ValidValues),
					fmt.Sprintf("ensure that all values for %s are in %v", p.Name, p.ValidValues),
				),
			)
		}
	}

	return err.ErrorOrNil()
}

func (p ParameterDefinition) validateTimezoneType(fieldType parameterFieldType, value any) error {
	tzErr := errors.NewError(fmt.Sprintf("invalid value for timezone for parameter %s", p.Name),
		"ensure that the value is one of the possible timezone values found here: https://github.com/observIQ/observiq-otel-collector/blob/main/receiver/pluginreceiver/timezone.go",
	)

	str, ok := value.(string)
	if !ok {
		return tzErr
	}

	if !validation.IsTimezone(str) {
		return tzErr
	}

	return nil
}

func (p ParameterDefinition) validateMetricsType(fieldType parameterFieldType, value any) error {
	return p.validateStringArrayValue(fieldType, value)
}

func (p ParameterDefinition) validateYamlValue(fieldType parameterFieldType, value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.NewError(
			fmt.Sprintf("expected a string for parameter %s", p.Name),
			fmt.Sprintf("ensure that the value is a string"),
		)
	}

	var into any
	return yaml.Unmarshal([]byte(str), &into)
}

func (p ParameterDefinition) validateMapValue(fieldType parameterFieldType, value any) error {
	reflectValue := reflect.ValueOf(value)
	kind := reflectValue.Kind()
	if kind != reflect.Map {
		return errors.NewError(
			fmt.Sprintf("expected type map for parameter %s but got %s", p.Name, kind),
			"ensure parameter is map[string]string",
		)
	}

	if m, ok := value.(map[string]any); ok {
		for _, v := range m {
			if k, ok := v.(string); !ok {
				return errors.NewError(
					fmt.Sprintf("expected type string for value for key %s in map", k),
					fmt.Sprintf("ensure all values in map are strings"),
				)
			}
		}
	}
	return nil
}

// Metrics returns the list of metrics associated with this parameter definition
func (p ParameterDefinition) metricNames(metricCategory string) []string {
	if p.Type != metricsType {
		return nil
	}
	results := []string{}

	for _, cat := range p.Options.MetricCategories {
		if cat.Label == metricCategory {
			for _, metric := range cat.Metrics {
				results = append(results, metric.Name)
			}
		}
	}

	return results
}
