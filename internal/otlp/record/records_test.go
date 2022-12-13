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
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestConvertLogs(t *testing.T) {
	testCases := []struct {
		name     string
		input    func() plog.Logs
		expected []*Log
	}{
		{
			name: "no logs",
			input: func() plog.Logs {
				return plog.NewLogs()
			},
			expected: []*Log{},
		},
		{
			name: "single log with string body",
			input: func() plog.Logs {
				l := plog.NewLogs()

				resource := l.ResourceLogs().AppendEmpty()
				resource.Resource().Attributes().PutStr("resource_id", "unique")

				logRecord := resource.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
				logRecord.Attributes().PutStr("custom_field", "custom_value")
				logRecord.Attributes().PutInt("db_id", 22)

				logRecord.SetSeverityText("ERROR")
				logRecord.SetSeverityNumber(plog.SeverityNumberError)
				logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Date(2022, time.September, 15, 1, 1, 1, 1, time.UTC)))
				pcommon.NewValueStr("log message").CopyTo(logRecord.Body())
				return l
			},
			expected: []*Log{
				{
					Attributes: map[string]interface{}{
						"custom_field": "custom_value",
						"db_id":        int64(22),
					},
					Body: "log message",
					Resource: map[string]interface{}{
						"resource_id": "unique",
					},
					Timestamp: time.Date(2022, time.September, 15, 1, 1, 1, 1, time.UTC),
					Severity:  "error",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ConvertLogs(tc.input()))
		})
	}
}

func TestGetLogMessage(t *testing.T) {

	testCases := []struct {
		name     string
		input    func() pcommon.Value
		expected string
	}{
		{
			name: "string value",
			input: func() pcommon.Value {
				return pcommon.NewValueStr("plain string")
			},
			expected: "plain string",
		},
		{
			name: "double value",
			input: func() pcommon.Value {
				return pcommon.NewValueDouble(1248.16)
			},
			expected: "1248.16",
		},
		{
			name: "int value",
			input: func() pcommon.Value {
				return pcommon.NewValueInt(4096)
			},
			expected: "4096",
		},
		{
			name: "bool value",
			input: func() pcommon.Value {
				return pcommon.NewValueBool(true)
			},
			expected: "true",
		},
		{
			name: "bytes value",
			input: func() pcommon.Value {
				// slice := pcommon.NewByteSlice()
				// slice.Append([]byte("log message in bytes")...)
				value := pcommon.NewValueBytes()
				value.SetEmptyBytes().Append([]byte("log message in bytes")...)
				return value
			},
			expected: base64.StdEncoding.EncodeToString([]byte("log message in bytes")),
		},
		{
			name: "empty value",
			input: func() pcommon.Value {
				return pcommon.NewValueEmpty()
			},
			expected: "",
		},
		{
			name: "slice value",
			input: func() pcommon.Value {
				body := pcommon.NewValueSlice()
				body.Slice().AppendEmpty().SetInt(30)
				body.Slice().AppendEmpty().SetInt(60)
				body.Slice().AppendEmpty().SetBool(false)
				body.Slice().AppendEmpty().SetStr("single string")
				return body
			},
			expected: `[30,60,false,"single string"]`,
		},
		{
			name: "map value",
			input: func() pcommon.Value {
				body := pcommon.NewValueMap()
				body.Map().PutStr("key1", "value1")
				body.Map().PutStr("message", "log message")
				body.Map().PutInt("pid", 333)
				return body
			},
			expected: `{
	-> key1: STRING(value1)
	-> message: STRING(log message)
	-> pid: INT(333)
}`,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, getLogMessage(tc.input()))
		})
	}
}
