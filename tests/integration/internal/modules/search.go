package modules

import (
	"fmt"
	"github.com/formancehq/search/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
	"github.com/ory/dockertest/v3"
	"net/http"
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
		internal.NewDockerService("public.ecr.aws/formance-internal/jeffail/benthos", "v4.23.0-es").
			WithEntrypoint([]string{
				"/benthos",
				"-c", "/config/config.yml",
				"-t", "/config/templates/*.yaml",
				"-r", "/config/resources/*.yaml",
				"streams", "/config/streams/ledger/v2.0.0/*.yaml", "/config/streams/payments/v0.0.0/*.yaml",
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
			}).
			WithHealthCheck(func(test *internal.Test, resource *dockertest.Resource) bool {
				rsp, err := http.Get("http://localhost:" + resource.GetPort("4195/tcp") + "/ping")
				if err != nil {
					return false
				}

				return rsp.StatusCode == http.StatusOK
			}),
	)
