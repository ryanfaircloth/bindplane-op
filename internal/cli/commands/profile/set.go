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

package profile

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/observiq/bindplane-op/internal/agent"
	"github.com/observiq/bindplane-op/internal/cli/flags"
	"github.com/observiq/bindplane-op/model"
)

// SetCommand returns the BindPlane profile set cobra command
func SetCommand(h Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <name>",
		Short: "set a parameter on a saved profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing required argument <name>")
			}

			name := args[0]
			f := h.Folder()

			profile, err := f.ReadProfile(name)
			if err != nil {
				profile = model.NewProfile(name, model.ProfileSpec{})
			}

			modifiers := registerModifiers()

			cmd.InheritedFlags().VisitAll(modifiers.handleFlag(profile))
			cmd.Flags().VisitAll(modifiers.handleFlag(profile))

			if modifiers.errs != nil {
				return modifiers.errs
			}

			err = f.WriteProfile(profile)
			if err != nil {
				return err
			}

			for _, m := range modifiers.messages() {
				fmt.Println(m)
			}
			return nil
		},
	}

	flags.Serve(cmd)

	return cmd
}

func registerModifiers() *profileSettingModifiers {
	p := newProfileSettingModifiers()

	p.register("port", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Port = f.Value.String()
		return nil
	})

	p.register("host", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Host = f.Value.String()
		return nil
	})

	p.register("server-url", func(name string, f *pflag.Flag, profile *model.Profile) error {
		serverAddress := f.Value.String()
		u, err := url.Parse(serverAddress)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		profile.Spec.Common.ServerURL = u.String()
		return nil
	})

	p.register("remote-url", func(name string, f *pflag.Flag, profile *model.Profile) error {
		remoteURL := f.Value.String()
		u, err := url.Parse(remoteURL)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "ws"
		}
		profile.Spec.Server.RemoteURL = u.String()
		return nil
	})

	p.register("secret-key", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.SecretKey = f.Value.String()
		return nil
	})

	p.register("username", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Username = f.Value.String()
		return nil
	})

	p.register("password", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Password = f.Value.String()
		return nil
	})

	p.register("storage-file-path", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.StorageFilePath = f.Value.String()
		return nil
	})

	p.register("tls-cert", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Common.Certificate = f.Value.String()
		return nil
	})

	p.register("tls-key", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Common.PrivateKey = f.Value.String()
		return nil
	})

	p.register("tls-ca", func(name string, f *pflag.Flag, profile *model.Profile) error {
		stringValue := f.Value.String()                                // In the p.register(of StringSlice this looks like `"[one,two]"
		value := strings.Split(stringValue[1:len(stringValue)-1], ",") // removes the brackets
		profile.Spec.Common.CertificateAuthority = value
		return nil
	})

	p.register("log-file-path", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Common.LogFilePath = f.Value.String()
		return nil
	})

	p.register("output", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Command.Output = f.Value.String()
		return nil
	})

	p.register("offline", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.Offline = f.Value.String() == "true"
		return nil
	})

	p.register("sync-agent-versions-interval", func(name string, f *pflag.Flag, profile *model.Profile) error {
		duration, err := time.ParseDuration(f.Value.String())
		if err != nil {
			return fmt.Errorf("failed to set sync-agent-versions-interval, must be a valid duration: %s", err.Error())
		}
		if 0 < duration && duration < agent.MinSyncAgentVersionsInterval {
			return fmt.Errorf("%s must be at least %s", f.Name, agent.MinSyncAgentVersionsInterval.String())
		}
		profile.Spec.Server.SyncAgentVersionsInterval = duration
		return nil
	})

	p.register("sessions-secret", func(name string, f *pflag.Flag, profile *model.Profile) error {
		// Try to enforce it as a UUID
		_, err := uuid.Parse(f.Value.String())
		if err != nil {
			return fmt.Errorf("failed to set sessions-secret, must be a UUID")
		}
		profile.Spec.Server.SessionsSecret = f.Value.String()
		return nil
	})

	p.register("trace-type", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.TraceType = f.Value.String()
		return nil
	})

	p.register("otlp-tracing-endpoint", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.OpenTelemetryTracing.Endpoint = f.Value.String()
		return nil
	})

	p.register("otlp-tracing-insecure-tls", func(name string, f *pflag.Flag, profile *model.Profile) error {
		profile.Spec.Server.OpenTelemetryTracing.TLS.Insecure = f.Value.String() == "true"
		return nil
	})

	return p
}

// ----------------------------------------------------------------------

type profileSettingModifier func(name string, f *pflag.Flag, profile *model.Profile) error

var present = struct{}{}

type profileSettingModifiers struct {
	flags    map[string]profileSettingModifier
	errs     error
	modified map[string]struct{}
	visited  map[string]struct{}
}

func newProfileSettingModifiers() *profileSettingModifiers {
	return &profileSettingModifiers{
		flags:    map[string]profileSettingModifier{},
		modified: map[string]struct{}{},
		visited:  map[string]struct{}{},
	}
}

func (p *profileSettingModifiers) register(name string, modifier profileSettingModifier) {
	p.flags[name] = modifier
}

func (p *profileSettingModifiers) messages() []string {
	result := make([]string, 0, len(p.modified))
	for m := range p.modified {
		result = append(result, fmt.Sprintf("%s modified", m))
	}
	return result
}

func (p *profileSettingModifiers) handleFlag(profile *model.Profile) func(f *pflag.Flag) {
	return func(f *pflag.Flag) {
		if f.Changed {
			// only handle each flag once (not Flags and InheritedFlags). otherwise we can get multiple errors for the same
			// flag.
			if _, ok := p.visited[f.Name]; ok {
				return
			}
			if modifier, ok := p.flags[f.Name]; ok {
				if err := modifier(f.Name, f, profile); err == nil {
					p.modified[f.Name] = present
				} else {
					p.errs = multierror.Append(p.errs, err)
				}
			}
			p.visited[f.Name] = present
		}
	}
}
