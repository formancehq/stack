package modules

import (
	"github.com/formancehq/auth/cmd"
	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const AuthIssuer = "http://localhost/api/auth"

var Auth = internal.NewModule("auth").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("auth", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {

				authDir := filepath.Join(os.TempDir(), uuid.NewString())
				if err := os.MkdirAll(authDir, 0o777); err != nil {
					panic(err)
				}
				type configuration struct {
					Clients []auth.StaticClient `yaml:"clients"`
				}
				cfg := &configuration{
					Clients: []auth.StaticClient{{
						ClientOptions: auth.ClientOptions{
							Name:    "global",
							Id:      "global",
							Trusted: true,
						},
						Secrets: []string{"global"},
					}},
				}
				configFile := filepath.Join(authDir, "config.yaml")
				f, err := os.Create(configFile)
				if err != nil {
					panic(err)
				}
				if err := yaml.NewEncoder(f).Encode(cfg); err != nil {
					panic(err)
				}

				if err := os.Setenv("CAOS_OIDC_DEV", "1"); err != nil {
					panic(err)
				}

				args := make([]string, 0)
				args = append(args)

				return []string{
					"serve",
					"--config=" + configFile,
					"--postgres-uri=" + test.GetDatabaseSourceName("auth"),
					"--delegated-client-id=" + internal.OIDCServer().ClientID,
					"--delegated-client-secret=" + internal.OIDCServer().ClientSecret,
					"--delegated-issuer=" + internal.OIDCServer().Issuer(),
					"--base-url=" + AuthIssuer,
					"--listen=0.0.0.0:0",
				}
			}),
	)
