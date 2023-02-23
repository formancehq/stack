package handlers

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
	"github.com/formancehq/operator/internal/modules"
	"github.com/formancehq/search/pkg/searchengine"
)

const (
	benthosImage = "jeffail/benthos:4.10.0"
)

func init() {
	modules.Register("search", modules.Module{
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{
				modules.Service{
					Port:               8080,
					HasVersionEndpoint: true,
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						env := elasticSearchEnvVars(resolveContext.Stack, resolveContext.Configuration).
							Append(
								modules.Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d%s",
									resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Host,
									resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Port,
									resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.PathPrefix)),
								modules.Env("OPEN_SEARCH_SCHEME", resolveContext.Stack.Spec.Scheme),
								modules.Env("ES_INDICES", resolveContext.Stack.Name),
								modules.Env("MAPPING_INIT_DISABLED", "true"),
							)
						if resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
							env = env.Append(
								modules.Env("OPEN_SEARCH_USERNAME", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
								modules.Env("OPEN_SEARCH_PASSWORD", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
							)
						}
						return modules.Container{
							Env:   env,
							Image: modules.GetImage("search", resolveContext.Versions.Spec.Search),
						}
					},
				},
				modules.Service{
					Name: "benthos",
					Port: 4195,
					Configs: func(resolveContext modules.InstallContext) modules.Configs {
						ret := modules.Configs{}
						for _, dir := range []string{"templates", "streams", "resources", "global"} {
							data := make(map[string]string)
							copyDir(benthosConfigDir, "benthos/"+dir, "benthos/"+dir, &data)
							ret[dir] = modules.Config{
								Mount: true,
								Data:  data,
							}
						}
						return ret
					},
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						env := elasticSearchEnvVars(resolveContext.Stack, resolveContext.Configuration).
							// Postgres config of the publishers (ledger and payments)
							// To be able to reindex data
							Append(modules.DefaultPostgresEnvVarsWithPrefix(
								resolveContext.Configuration.Spec.Services.Ledger.Postgres,
								resolveContext.Stack.GetServiceName("ledger"),
								"LEDGER_",
							)...)

						if resolveContext.Configuration.Spec.Broker.Kafka != nil {
							env = env.Append(modules.Env("BROKER", "kafka"))
						} else {
							env = env.Append(modules.Env("BROKER", "nats"))
						}

						return modules.Container{
							Env:   env,
							Image: benthosImage,
							Command: []string{
								"/benthos",
								"-r", resolveContext.GetConfig("resources").GetMountPath() + "/*.yaml",
								"-t", resolveContext.GetConfig("templates").GetMountPath() + "/*.yaml",
								"-c", resolveContext.GetConfig("global").GetMountPath() + "/config.yaml",
								"--log.level", "trace", "streams",
								resolveContext.GetConfig("streams").GetMountPath() + "/*.yaml",
							},
							Liveness:             modules.LivenessDisable,
							DisableRollingUpdate: true,
						}
					},
					InitContainer: func(resolveContext modules.ContainerResolutionContext) []modules.Container {
						env := modules.ContainerEnv{
							modules.Env("OPEN_SEARCH_HOST", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Host),
							modules.Env("OPEN_SEARCH_PORT", fmt.Sprint(resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Port)),
							modules.Env("OPEN_SEARCH_PATH_PREFIX", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.PathPrefix),
							modules.Env("OPEN_SEARCH_SCHEME", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Scheme),
							modules.Env("OPEN_SEARCH_SERVICE", controllerutils.ComputeEnvVar("", "%s:%s%s",
								"OPEN_SEARCH_HOST",
								"OPEN_SEARCH_PORT",
								"OPEN_SEARCH_PATH_PREFIX",
							)),
						}
						if resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
							env = env.Append(
								modules.Env("OPEN_SEARCH_USERNAME", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
								modules.Env("OPEN_SEARCH_PASSWORD", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
							)
						}

						credentialsStr := ""
						if resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
							credentialsStr = "-u ${OPEN_SEARCH_USERNAME}:${OPEN_SEARCH_PASSWORD} "
						}

						mapping, err := searchengine.GetIndexDefinition()
						if err != nil {
							panic(err)
						}

						var args []string
						if resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.UseZinc {
							if err != nil {
								panic(err)
							}
							args = []string{
								"-c", fmt.Sprintf("curl -H 'Content-Type: application/json' "+
									"-X POST -v -d '%s' "+
									credentialsStr+
									"${OPEN_SEARCH_SCHEME}://${OPEN_SEARCH_SERVICE}/index", string(mapping)),
							}
						} else {
							args = []string{
								"-c", fmt.Sprintf("curl -H 'Content-Type: application/json' "+
									"-X PUT -v -d '%s' "+
									credentialsStr+
									"${OPEN_SEARCH_SCHEME}://${OPEN_SEARCH_SERVICE}/%s", string(mapping), resolveContext.Stack.Name),
							}
						}

						return []modules.Container{{
							Command: []string{"sh"},
							Env:     env,
							Image:   "curlimages/curl:7.86.0",
							Name:    "init-mapping",
							Args:    args,
						}}
					},
				},
			}
		},
	})
}

func copyDir(f fs.FS, root, path string, ret *map[string]string) {
	dirEntries, err := fs.ReadDir(f, path)
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		dirEntryPath := filepath.Join(path, dirEntry.Name())
		if dirEntry.IsDir() {
			copyDir(f, root, dirEntryPath, ret)
		} else {
			fileContent, err := fs.ReadFile(f, dirEntryPath)
			if err != nil {
				panic(err)
			}
			sanitizedPath := strings.TrimPrefix(dirEntryPath, root)
			sanitizedPath = strings.TrimPrefix(sanitizedPath, "/")
			(*ret)[sanitizedPath] = string(fileContent)
		}
	}
}

func elasticSearchEnvVars(stack *stackv1beta3.Stack, configuration *stackv1beta3.Configuration) modules.ContainerEnv {
	ret := modules.ContainerEnv{
		modules.Env("OPENSEARCH_URL", configuration.Spec.Services.Search.ElasticSearchConfig.Endpoint()),
		modules.Env("OPENSEARCH_INDEX", stack.Name),
		modules.Env("OPENSEARCH_BATCHING_COUNT", fmt.Sprint(configuration.Spec.Services.Search.Batching.Count)),
		modules.Env("OPENSEARCH_BATCHING_PERIOD", configuration.Spec.Services.Search.Batching.Period),
		modules.Env("TOPIC_PREFIX", stack.Name+"-"),
	}
	if configuration.Spec.Broker.Kafka != nil {
		ret = ret.Append(
			modules.Env("KAFKA_ADDRESS", strings.Join(configuration.Spec.Broker.Kafka.Brokers, ",")),
		)
		if configuration.Spec.Broker.Kafka.TLS {
			ret = ret.Append(
				modules.Env("KAFKA_TLS_ENABLED", "true"),
			)
		}
		if configuration.Spec.Broker.Kafka.SASL != nil {
			ret = ret.Append(
				modules.Env("KAFKA_SASL_USERNAME", configuration.Spec.Broker.Kafka.SASL.Username),
				modules.Env("KAFKA_SASL_PASSWORD", configuration.Spec.Broker.Kafka.SASL.Password),
				modules.Env("KAFKA_SASL_MECHANISM", configuration.Spec.Broker.Kafka.SASL.Mechanism),
			)
		}
	}
	if configuration.Spec.Broker.Nats != nil {
		ret = ret.Append(
			modules.Env("NATS_URL", configuration.Spec.Broker.Nats.URL),
		)
	}
	if configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
		ret = ret.Append(
			modules.Env("BASIC_AUTH_ENABLED", "true"),
			modules.Env("BASIC_AUTH_USERNAME", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
			modules.Env("BASIC_AUTH_PASSWORD", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
		)
	}
	return ret
}
