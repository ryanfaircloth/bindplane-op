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

package record

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// Metric is a metric record sent to bindplane
type Metric struct {
	Name           string                 `json:"name"`
	Timestamp      time.Time              `json:"timestamp"`
	StartTimestamp time.Time              `json:"start_timestamp,omitempty"`
	Value          interface{}            `json:"value"`
	Unit           string                 `json:"unit"`
	Type           string                 `json:"type"`
	Attributes     map[string]interface{} `json:"attributes"`
	Resource       map[string]interface{} `json:"resource"`
}

// Log is a log record sent to bindplane
type Log struct {
	Timestamp  time.Time              `json:"timestamp"`
	Body       interface{}            `json:"body"`
	Severity   string                 `json:"severity"`
	Attributes map[string]interface{} `json:"attributes"`
	Resource   map[string]interface{} `json:"resource"`
}

// Trace is a trace record sent to bindplane
type Trace struct {
	Name         string                 `json:"name"`
	TraceID      string                 `json:"trace_id"`
	SpanID       string                 `json:"span_id"`
	ParentSpanID string                 `json:"parent_span_id"`
	Start        time.Time              `json:"start"`
	End          time.Time              `json:"end"`
	Attributes   map[string]interface{} `json:"attributes"`
	Resource     map[string]interface{} `json:"resource"`
}

// ConvertMetrics gets metric records from pmetrics
func ConvertMetrics(metrics pmetric.Metrics) []*Metric {
	records := []*Metric{}
	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		resourceMetrics := metrics.ResourceMetrics().At(i)
		resourceAttributes := resourceMetrics.Resource().Attributes().AsRaw()
		for k := 0; k < resourceMetrics.ScopeMetrics().Len(); k++ {
			scopeMetrics := resourceMetrics.ScopeMetrics().At(k)
			for j := 0; j < scopeMetrics.Metrics().Len(); j++ {
				metric := scopeMetrics.Metrics().At(j)
				recordSlice := getRecordsFromMetric(metric, resourceAttributes)
				records = append(records, recordSlice...)
			}
		}
	}
	return records
}

// getRecords gets metric records from a pmetric
func getRecordsFromMetric(metric pmetric.Metric, resourceAttributes map[string]interface{}) []*Metric {
	switch metric.DataType() {
	case pmetric.MetricDataTypeSum:
		return getMetricRecordsFromSum(metric, resourceAttributes)
	case pmetric.MetricDataTypeGauge:
		return getMetricRecordsFromGauge(metric, resourceAttributes)
	case pmetric.MetricDataTypeSummary:
		return getMetricRecordsFromSummary(metric, resourceAttributes)
	}

	return nil
}

// getMetricRecordsFromSum converts a sum into metric records
func getMetricRecordsFromSum(metric pmetric.Metric, resourceAttributes map[string]interface{}) []*Metric {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []*Metric{}
	sum := metric.Sum()
	points := sum.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := Metric{
			Name:           metricName,
			Timestamp:      point.Timestamp().AsTime(),
			StartTimestamp: point.StartTimestamp().AsTime(),
			Value:          getDataPointValue(point),
			Unit:           metricUnit,
			Type:           metricType,
			Attributes:     point.Attributes().AsRaw(),
			Resource:       resourceAttributes,
		}
		metricRecords = append(metricRecords, &record)
	}
	return metricRecords
}

// getMetricRecordsFromGauge converts a gauge into metric records
func getMetricRecordsFromGauge(metric pmetric.Metric, resourceAttributes map[string]interface{}) []*Metric {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []*Metric{}
	gauge := metric.Gauge()
	points := gauge.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := Metric{
			Name:       metricName,
			Timestamp:  point.Timestamp().AsTime(),
			Value:      getDataPointValue(point),
			Unit:       metricUnit,
			Type:       metricType,
			Attributes: point.Attributes().AsRaw(),
			Resource:   resourceAttributes,
		}
		metricRecords = append(metricRecords, &record)
	}
	return metricRecords
}

// getMetricRecordsFromSummary converts a summary into metric records
func getMetricRecordsFromSummary(metric pmetric.Metric, resourceAttributes map[string]interface{}) []*Metric {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []*Metric{}
	summary := metric.Summary()
	points := summary.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := Metric{
			Name:           metricName,
			Timestamp:      point.Timestamp().AsTime(),
			StartTimestamp: point.StartTimestamp().AsTime(),
			Value:          getSummaryPointValue(point),
			Unit:           metricUnit,
			Type:           metricType,
			Attributes:     point.Attributes().AsRaw(),
			Resource:       resourceAttributes,
		}
		metricRecords = append(metricRecords, &record)
	}
	return metricRecords
}

// getDataPointValue gets the value of a data point
func getDataPointValue(point pmetric.NumberDataPoint) interface{} {
	switch point.ValueType() {
	case pmetric.NumberDataPointValueTypeDouble:
		return point.DoubleVal()
	case pmetric.NumberDataPointValueTypeInt:
		return point.IntVal()
	default:
		return 0
	}
}

// getSummaryPointValue gets the value of a summary point
func getSummaryPointValue(point pmetric.SummaryDataPoint) map[float64]interface{} {
	value := make(map[float64]interface{})
	for i := 0; i < point.QuantileValues().Len(); i++ {
		q := point.QuantileValues().At(i)
		value[q.Quantile()] = q.Value()
	}
	return value
}

// ConvertLogs gets log records from plogs
func ConvertLogs(logs plog.Logs) []*Log {
	records := []*Log{}
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		resourceLogs := logs.ResourceLogs().At(i)
		resourceAttributes := resourceLogs.Resource().Attributes().AsRaw()
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				record := Log{
					Timestamp:  log.Timestamp().AsTime(),
					Body:       getLogMessage(log.Body()),
					Severity:   severityNumberToString(int32(log.SeverityNumber())),
					Attributes: log.Attributes().AsRaw(),
					Resource:   resourceAttributes,
				}
				records = append(records, &record)
			}
		}
	}
	return records
}

// getLogMessage retrieves the body of a plog.LogEntry to display in snapshots
//
// If the body is a Map, it's formatted the same way as in the logging exporter.
// Otherwise, AsString() is returned.
func getLogMessage(body pcommon.Value) string {
	if body.Type() != pcommon.ValueTypeMap {
		return body.AsString()
	}

	var b strings.Builder
	b.WriteString("{\n")
	// Sort to ensure logs with the same attributes display in the same order
	body.MapVal().Sort().Range(func(k string, v pcommon.Value) bool {
		fmt.Fprintf(&b, "\t-> %s: %s(%s)\n", k, v.Type(), v.AsString())
		return true
	})
	b.WriteByte('}')
	return b.String()
}

var severityMap = map[int32]string{
	1:  "trace",
	2:  "trace",
	3:  "trace",
	4:  "trace",
	5:  "debug",
	6:  "debug",
	7:  "debug",
	8:  "debug",
	9:  "info",
	10: "info",
	11: "info",
	12: "info",
	13: "warning",
	14: "warning",
	15: "warning",
	16: "warning",
	17: "error",
	18: "error",
	19: "error",
	20: "error",
	21: "fatal",
	22: "fatal",
	23: "fatal",
	24: "fatal",
}

func severityNumberToString(num int32) string {
	if sev, ok := severityMap[num]; ok {
		return sev
	}

	return "default"
}

// ConvertTraces gets trace records from ptraces
func ConvertTraces(traces ptrace.Traces) []*Trace {
	records := []*Trace{}
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		resourceAttributes := resourceSpans.Resource().Attributes().AsRaw()
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				span.TraceID().HexString()
				record := Trace{
					Name:         span.Name(),
					TraceID:      span.TraceID().HexString(),
					SpanID:       span.SpanID().HexString(),
					ParentSpanID: span.ParentSpanID().HexString(),
					Start:        span.StartTimestamp().AsTime(),
					End:          span.EndTimestamp().AsTime(),
					Attributes:   span.Attributes().AsRaw(),
					Resource:     resourceAttributes,
				}
				records = append(records, &record)
			}
		}
	}
	return records
}

// AttributeString returns an Attribute as a string or the defaultValue if either an attribute with that name does not
// exist or an attribute with that name does exist but is not a string.
func (m *Metric) AttributeString(name string, defaultValue string) string {
	if id, ok := m.Attributes[name]; ok {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return defaultValue
}

// Clone returns a deep copy of the object by using JSON Marshal/Unmarshal.
func Clone[T any](src *T) (*T, error) {
	s, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	var dest T
	if err = json.Unmarshal(s, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
