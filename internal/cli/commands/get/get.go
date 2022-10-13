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

	"github.com/spf13/cobra"

	"github.com/observiq/bindplane-op/client"
	"github.com/observiq/bindplane-op/internal/cli"
	"github.com/observiq/bindplane-op/internal/cli/printer"
	"github.com/observiq/bindplane-op/model"
)

// Command returns the BindPlane get cobra command.
func Command(bindplane *cli.BindPlane) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more resources",
	}

	cmd.AddCommand(
		ResourcesCommand(bindplane),
		AgentsCommand(bindplane),
		AgentVersionsCommand(bindplane),
		ConfigurationsCommand(bindplane),
		DestinationsCommand(bindplane),
		DestinationTypesCommand(bindplane),
		ProcessorsCommand(bindplane),
		ProcessorTypesCommand(bindplane),
		SourcesCommand(bindplane),
		SourceTypesCommand(bindplane),
	)

	return cmd
}

// ----------------------------------------------------------------------
// generic implementations for get

type getter[T model.Printable] struct {
	one func(ctx context.Context, client client.BindPlane, name string) (T, bool, error)
	all func(ctx context.Context, client client.BindPlane) ([]T, error)
}

func getImpl[T model.Printable](bindplane *cli.BindPlane, resourceName string, g getter[T]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := bindplane.Client()
		if err != nil {
			return fmt.Errorf("error creating client: %w", err)
		}

		if len(args) > 0 {
			name := args[0]
			item, exists, err := g.one(cmd.Context(), c, name)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("no %s found with name %s", resourceName, name)
			}

			printer.PrintResource(bindplane.Printer(), item)
			return nil
		}

		items, err := g.all(cmd.Context(), c)
		if err != nil {
			return err
		}

		printer.PrintResources(bindplane.Printer(), items)
		return nil
	}
}
