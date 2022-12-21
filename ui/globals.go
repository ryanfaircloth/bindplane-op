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
	"bytes"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/observiq/bindplane-op/common"
	"github.com/observiq/bindplane-op/internal/server"
	"github.com/observiq/bindplane-op/internal/version"
)

const templateStr = `var __BINDPLANE_VERSION__ = "{{.Version}}";
`

const fallbackJs = `var globals = {
	version: "unknown"
}
`

type configOptions struct {
	Version string
}

func newConfigOptions(config *common.Server) *configOptions {
	return &configOptions{
		Version: version.NewVersion().String(),
	}
}

// generateGlobalJS generates the static javascript file for the UI.
func generateGlobalJS(config *common.Server) (string, error) {
	tmp, err := template.New("globals").Parse(templateStr)
	if err != nil {
		return fallbackJs, err
	}

	opts := newConfigOptions(config)

	w := bytes.NewBufferString("")
	err = tmp.Execute(w, opts)
	if err != nil {
		return fallbackJs, err
	}

	return w.String(), nil
}

func globalJS(ctx *gin.Context, bindplane server.BindPlane) {
	js, err := generateGlobalJS(bindplane.Config())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", "application/javascript")
	ctx.String(http.StatusOK, js)

}
