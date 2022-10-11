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

package stats

import "testing"

func TestParseProcessorName(t *testing.T) {
	tests := []struct {
		processorName    string
		wantPosition     string
		wantPipelineType string
		wantName         string
	}{
		{
			processorName:    "something_else/_s0_metrics_source0",
			wantPosition:     "",
			wantPipelineType: "",
			wantName:         "",
		},
		{
			processorName:    "throughputmeasurement/_s0_metrics_source0",
			wantPosition:     "s0",
			wantPipelineType: "metrics",
			wantName:         "source0",
		},
		{
			processorName:    "throughputmeasurement/_d1_metrics_error-logging",
			wantPosition:     "d1",
			wantPipelineType: "metrics",
			wantName:         "error-logging",
		},
		{
			processorName:    "throughputmeasurement/_d0_metrics_error_logging",
			wantPosition:     "d0",
			wantPipelineType: "metrics",
			wantName:         "error_logging",
		},
		{
			processorName:    "throughputmeasurement/_d0_metrics_error_logging_and_other_stuff",
			wantPosition:     "d0",
			wantPipelineType: "metrics",
			wantName:         "error_logging_and_other_stuff",
		},
	}
	for _, test := range tests {
		t.Run(test.processorName, func(t *testing.T) {
			gotPosition, gotPipelineType, gotName := ParseProcessorName(test.processorName)
			if gotPosition != test.wantPosition {
				t.Errorf("ParseProcessorName() gotPosition = %v, want %v", gotPosition, test.wantPosition)
			}
			if gotPipelineType != test.wantPipelineType {
				t.Errorf("ParseProcessorName() gotPipelineType = %v, want %v", gotPipelineType, test.wantPipelineType)
			}
			if gotName != test.wantName {
				t.Errorf("ParseProcessorName() gotName = %v, want %v", gotName, test.wantName)
			}
		})
	}
}
