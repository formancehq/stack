package search

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
	"github.com/formancehq/operator/internal/modules"
	benthosOperator "github.com/formancehq/operator/internal/modules/search/benthos"
	"github.com/formancehq/search/benthos"
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/opensearch-project/opensearch-go"
)

const (
	benthosImage = "public.ecr.aws/h9j1u6h3/jeffail/benthos:4.12.1"
)

type module struct{}

func (s module) Name() string {
	return "search"
}

func (s module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{searchService(ctx), benthosService(ctx)}
			},
			Cron: reindexCron,
		},
		"v0.7.0": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				esClient, err := getOpenSearchClient(config)
				if err != nil {
					return err
				}
				if err := searchengine.CreateIndex(ctx, esClient, stackv1beta3.DefaultESIndex); err != nil {
					return err
				}
				if err := reindexData(ctx, config, esClient); err != nil {
					return err
				}

				return nil
			},
			PostUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				esClient, err := getOpenSearchClient(config)
				if err != nil {
					return err
				}
				if err := reindexData(ctx, config, esClient); err != nil {
					return err
				}

				response, err := esClient.Indices.Delete([]string{config.Stack.Name}, esClient.Indices.Delete.WithContext(ctx))
				if err != nil {
					return err
				}

				if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
					return fmt.Errorf("unexpected status code %d when deleting index: %s", response.StatusCode, config.Stack.Name)
				}
				return nil
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{searchService(ctx), benthosService(ctx)}
			},
			Cron: reindexCron,
		},
	}
}

var Module = &module{}

var _ modules.Module = Module

func init() {
	modules.Register(Module)
}

var CreateOpenSearchClient = func(cfg opensearch.Config) (*opensearch.Client, error) {
	return opensearch.NewClient(cfg)
}

func reindexCron(ctx modules.ReconciliationConfig) []modules.Cron {
	return []modules.Cron{
		{
			Container: modules.Container{
				Command: []string{
					"/bin/sh", "-c",
					fmt.Sprintf("curl http://search-benthos.%s.svc.cluster.local:4195/ledger_reindex_all -X POST -H 'Content-Type: application/json' -d '{}'", ctx.Stack.Name),
				},
				Image: "curlimages/curl:8.2.1",
				Name:  "reindex-ledger",
			},
			Schedule: "* * * * *",
			Suspend:  true,
		},
	}
}

func getOpenSearchClient(ctx modules.ReconciliationConfig) (*opensearch.Client, error) {
	opensearchConfig := opensearch.Config{
		Addresses:            []string{ctx.Configuration.Spec.Services.Search.ElasticSearchConfig.Endpoint()},
		UseResponseCheckOnly: true,
	}

	if ctx.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
		opensearchConfig.Username = ctx.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username
		opensearchConfig.Password = ctx.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password
	}

	return CreateOpenSearchClient(opensearchConfig)
}

func reindexData(ctx context.Context, config modules.ReconciliationConfig, esClient *opensearch.Client) error {
	res, err := esClient.Reindex(bytes.NewBufferString(fmt.Sprintf(`{
	  "source": {
		"index": "%s"
	  },
	  "dest": {
		"index": "%s"
	  },
	  "script": {
		"source": "ctx._source.stack = ctx._index; ctx._id = ctx._index + '-' + ctx._id",
		"lang": "painless"
	  }
	}`, config.Stack.Name, stackv1beta3.DefaultESIndex)),
		esClient.Reindex.WithContext(ctx),
		esClient.Reindex.WithRefresh(true),
	)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code %d when reindexing data", res.StatusCode)
	}
	return nil
}

func searchService(ctx modules.ReconciliationConfig) *modules.Service {
	return &modules.Service{
		ListenEnvVar:       "BIND",
		ExposeHTTP:         modules.DefaultExposeHTTP,
		HasVersionEndpoint: true,
		Annotations:        ctx.Configuration.Spec.Services.Search.Annotations.Service,
		Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
			env := elasticSearchEnvVars(resolveContext.Stack, resolveContext.Configuration, resolveContext.Versions).
				Append(
					modules.Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d%s",
						resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Host,
						resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Port,
						resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.PathPrefix)),
					modules.Env("OPEN_SEARCH_SCHEME", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.Scheme),
					modules.Env("MAPPING_INIT_DISABLED", "true"),
				)
			if resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
				env = env.Append(
					modules.Env("OPEN_SEARCH_USERNAME", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
					modules.Env("OPEN_SEARCH_PASSWORD", resolveContext.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
				)
			}
			if resolveContext.Versions.IsLower("search", "v0.7.0") {
				env = env.Append(modules.Env("ES_INDICES", resolveContext.Stack.Name))
			} else {
				env = env.Append(modules.Env("ES_INDICES", stackv1beta3.DefaultESIndex))
			}
			return modules.Container{
				Env:   env,
				Image: modules.GetImage("search", resolveContext.Versions.Spec.Search),
				Resources: modules.GetResourcesWithDefault(
					resolveContext.Configuration.Spec.Services.Search.SearchResourceProperties,
					modules.ResourceSizeSmall(),
				),
			}
		},
	}
}

func benthosService(ctx modules.ReconciliationConfig) *modules.Service {
	ret := &modules.Service{
		Name:        "benthos",
		Port:        4195,
		ExposeHTTP:  modules.DefaultExposeHTTP,
		Liveness:    modules.LivenessDisable,
		Annotations: ctx.Configuration.Spec.Services.Search.Annotations.Service,
		Configs: func(resolveContext modules.ServiceInstallConfiguration) modules.Configs {
			ret := modules.Configs{}

			type directory struct {
				name string
				fs   embed.FS
			}

			directories := []directory{
				{
					name: "templates",
					fs:   benthos.Templates,
				},
				{
					name: "resources",
					fs:   benthos.Resources,
				},
				{
					name: "streams",
					fs:   benthos.Streams,
				},
			}
			if ctx.Configuration.Spec.Monitoring != nil {
				directories = append(directories, directory{
					name: "global",
					fs:   benthosOperator.Global,
				})
			}

			for _, x := range directories {
				data := make(map[string]string)
				copyDir(x.fs, x.name, x.name, &data)
				ret[x.name] = modules.Config{
					Mount: true,
					Data:  data,
				}
			}
			return ret
		},
		Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
			env := elasticSearchEnvVars(resolveContext.Stack, resolveContext.Configuration, resolveContext.Versions).
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

			cmd := []string{
				"/benthos",
				"-r", resolveContext.GetConfig("resources").GetMountPath() + "/*.yaml",
				"-t", resolveContext.GetConfig("templates").GetMountPath() + "/*.yaml",
			}
			if ctx.Configuration.Spec.Monitoring != nil {
				cmd = append(cmd, "-c", resolveContext.GetConfig("global").GetMountPath()+"/config.yaml")
			}
			cmd = append(cmd,
				"--log.level", "trace", "streams",
				resolveContext.GetConfig("streams").GetMountPath()+"/*.yaml")

			return modules.Container{
				Env:                  env,
				Image:                benthosImage,
				Command:              cmd,
				DisableRollingUpdate: true,
				Resources: modules.GetResourcesWithDefault(
					resolveContext.Configuration.Spec.Services.Search.BenthosResourceProperties,
					modules.ResourceSizeSmall(),
				),
			}
		},
	}

	if ctx.Versions.IsLower("search", "v0.7.0") {
		ret.InitContainer = initContainerCreateIndex()
	}

	return ret
}

func initContainerCreateIndex() func(resolveContext modules.ContainerResolutionConfiguration) []modules.Container {
	return func(resolveContext modules.ContainerResolutionConfiguration) []modules.Container {
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
	}
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

func elasticSearchEnvVars(stack *stackv1beta3.Stack, configuration *stackv1beta3.Configuration, versions *stackv1beta3.Versions) modules.ContainerEnv {
	ret := modules.ContainerEnv{
		modules.Env("OPENSEARCH_URL", configuration.Spec.Services.Search.ElasticSearchConfig.Endpoint()),
		modules.Env("OPENSEARCH_BATCHING_COUNT", fmt.Sprint(configuration.Spec.Services.Search.Batching.Count)),
		modules.Env("OPENSEARCH_BATCHING_PERIOD", configuration.Spec.Services.Search.Batching.Period),
		modules.Env("TOPIC_PREFIX", stack.Name+"-"),
	}
	if versions.IsLower("search", "v0.7.0") {
		ret = append(ret, modules.Env("OPENSEARCH_INDEX", stack.Name))
	} else {
		ret = append(ret, modules.Env("OPENSEARCH_INDEX", stackv1beta3.DefaultESIndex))
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
