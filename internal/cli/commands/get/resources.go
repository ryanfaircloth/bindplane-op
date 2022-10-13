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

package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/observiq/bindplane-op/client"
	"github.com/observiq/bindplane-op/internal/cli"
	"github.com/observiq/bindplane-op/internal/cli/printer"
	"github.com/observiq/bindplane-op/model"
	"github.com/spf13/cobra"
)

// ResourcesCommand returns the BindPlane get source-types cobra command
func ResourcesCommand(bindplane *cli.BindPlane) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resources",
		Short: "Displays all resources",
		Long:  `Use -o yaml to export all resources to yaml.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := bindplane.Client()
			if err != nil {
				return fmt.Errorf("error creating client: %w", err)
			}

			ctx := cmd.Context()
			p := bindplane.Printer()

			resources := []resource[model.Printable]{
				{name: "configurations", get: func() ([]model.Printable, error) { return cvt(c.Configurations(ctx)) }},
				{name: "sources", get: func() ([]model.Printable, error) { return cvt(c.Sources(ctx)) }},
				{name: "processors", get: func() ([]model.Printable, error) { return cvt(c.Processors(ctx)) }},
				{name: "destinations", get: func() ([]model.Printable, error) { return cvt(c.Destinations(ctx)) }},
				{name: "source-types", get: func() ([]model.Printable, error) { return cvt(c.SourceTypes(ctx)) }},
				{name: "processor-types", get: func() ([]model.Printable, error) { return cvt(c.ProcessorTypes(ctx)) }},
				{name: "destination-types", get: func() ([]model.Printable, error) { return cvt(c.DestinationTypes(ctx)) }},
				{name: "agent-versions", get: func() ([]model.Printable, error) { return cvt(c.AgentVersions(ctx)) }},
			}

			switch p.(type) {
			case *printer.YamlPrinter:
				err = printAll(ctx, c, resources, p)
			case *printer.JSONPrinter:
				err = printAll(ctx, c, resources, p)
			default:
				err = printEachTable(ctx, resources, p)
			}

			return nil
		},
	}
	return cmd
}

func printAll(ctx context.Context, c client.BindPlane, resources []resource[model.Printable], p printer.Printer) error {
	var list []model.Printable

	for _, r := range resources {
		items, err := r.get()
		if err != nil {
			return err
		}
		// note that we prepend to reverse the order of resources
		list = append(items, list...)
	}

	p.PrintResources(list)
	return nil
}

func printEachTable(ctx context.Context, resources []resource[model.Printable], p printer.Printer) error {
	for _, r := range resources {
		items, err := r.get()
		if err != nil {
			return err
		}
		fmt.Printf("\n%s\n%s\n", r.name, strings.Repeat("-", len(r.name)))
		printer.PrintResources(p, items)
	}

	return nil
}

type getResources[T any] func() ([]T, error)

type resource[T any] struct {
	name string
	get  getResources[T]
}

func cvt[T model.Printable](resources []T, err error) ([]model.Printable, error) {
	if err != nil {
		return nil, err
	}
	result := []model.Printable{}
	for _, item := range resources {
		result = append(result, item)
	}
	return result, nil
}
