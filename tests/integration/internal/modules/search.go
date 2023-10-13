package modules

import (
	"fmt"
	"github.com/formancehq/search/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Search = internal.NewModule("search").
	WithServices(
		internal.NewCommandService("search", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--open-search-service=" + internal.GetOpenSearchUrl(),
					"--open-search-scheme=http",
					"--open-search-username=admin",
					"--open-search-password=admin",
					"--bind=0.0.0.0:0",
					"--stack=" + test.ID(),
					fmt.Sprintf("--es-indices=%s", test.ID()),
					//"--mapping-init-disabled",
				}
			}),
		internal.NewDockerService("jeffail/benthos", "4.11").
			WithEntrypoint([]string{
				"/benthos",
				"-c", "/config/config.yml",
				"-t", "/config/templates/*.yaml",
				"-r", "/config/resources/*.yaml",
				"streams", "/config/streams/*.yaml",
			}).
			WithEnv(func(test *internal.Test) []string {
				return []string{
					fmt.Sprintf("OPENSEARCH_URL=http://%s:9200", internal.GetDockerEndpoint()),
					"BASIC_AUTH_ENABLED=true",
					"BASIC_AUTH_USERNAME=admin",
					"BASIC_AUTH_PASSWORD=admin",
					fmt.Sprintf("OPENSEARCH_INDEX=%s", test.ID()),
					fmt.Sprintf("NATS_URL=nats://%s:4222", internal.GetDockerEndpoint()),
					fmt.Sprintf("STACK=%s", test.ID()),
					fmt.Sprintf("TOPIC_PREFIX=%s-", test.ID()),
				}
			}).
			WithMounts(func(test *internal.Test) []string {
				return []string{
					test.Workdir() + "/../../../components/search/benthos:/config",
				}
			}),
	)
