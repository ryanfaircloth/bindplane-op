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

package otlp

import (
	"testing"

	"github.com/observiq/bindplane-op/internal/otlp/record"
	"github.com/observiq/bindplane-op/internal/store/stats"
	"github.com/stretchr/testify/require"
)

func TestIncludeAgentMetrics(t *testing.T) {
	makeMetric := func(processor string) *record.Metric {
		return &record.Metric{
			Name: "name",
			Attributes: map[string]any{
				stats.ProcessorAttributeName: processor,
			},
		}
	}

	const (
		s  = "throughputmeasurement/_s_metrics_source0"
		s0 = "throughputmeasurement/_s0_metrics_source0"
		s1 = "throughputmeasurement/_s1_metrics_source0"
		d  = "throughputmeasurement/_d_metrics_source0"
		d0 = "throughputmeasurement/_d0_metrics_source0"
		d1 = "throughputmeasurement/_d1_metrics_source0"
	)

	t.Run("no replacement", func(t *testing.T) {
		metric := makeMetric(s0)
		metrics := []*record.Metric{}
		metrics = includeAgentMetrics(metrics, metric)
		require.Len(t, metrics, 1)
		require.Equal(t, s0, metrics[0].Attributes[stats.ProcessorAttributeName])
	})

	t.Run("s0/s1 replacement", func(t *testing.T) {
		metric := makeMetric(s)
		metrics := []*record.Metric{}
		metrics = includeAgentMetrics(metrics, metric)
		require.Len(t, metrics, 2)
		require.Equal(t, s0, metrics[0].Attributes[stats.ProcessorAttributeName])
		require.Equal(t, s1, metrics[1].Attributes[stats.ProcessorAttributeName])
	})

	t.Run("d0/d1 replacement", func(t *testing.T) {
		metric := makeMetric(d)
		metrics := []*record.Metric{}
		metrics = includeAgentMetrics(metrics, metric)
		require.Len(t, metrics, 2)
		require.Equal(t, d0, metrics[0].Attributes[stats.ProcessorAttributeName])
		require.Equal(t, d1, metrics[1].Attributes[stats.ProcessorAttributeName])
	})
}
