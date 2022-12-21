// Copyright observIQ, Inc.
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

package ui

import (
	"fmt"
	"testing"

	"github.com/observiq/bindplane-op/common"
	"github.com/observiq/bindplane-op/internal/version"
	"github.com/stretchr/testify/require"
)

func TestGenerateGlobalJS(t *testing.T) {
	type args struct {
		config *common.Server
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			"renders template",
			args{
				config: &common.Server{},
			},
			false,
			fmt.Sprintf("var __BINDPLANE_VERSION__ = \"%s\";\n", version.NewVersion().String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateGlobalJS(tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
