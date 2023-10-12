package search

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	benthosOperator "github.com/formancehq/operator/internal/modules/search/benthos"
	"github.com/formancehq/search/benthos"
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
			PreUpgrade: func(ctx context.Context, jobRunner modules.JobRunner, config modules.ReconciliationConfig) (bool, error) {
				return jobRunner.RunJob(ctx, "create-index-mapping", nil, initMappingJob(config))
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{searchService(ctx), benthosService(ctx)}
			},
			Cron: reindexCron,
		},
		"v0.7.0": {
			PreUpgrade: func(ctx context.Context, jobRunner modules.JobRunner, config modules.ReconciliationConfig) (bool, error) {
				ok, err := jobRunner.RunJob(ctx, "create-index-mapping", nil, initMappingJob(config))
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}

				ok, err = jobRunner.RunJob(ctx, "reindex-data", nil, reindexDataJob(config))
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}

				return true, nil
			},
			PostUpgrade: func(ctx context.Context, jobRunner modules.JobRunner, config modules.ReconciliationConfig) (bool, error) {

				ok, err := jobRunner.RunJob(ctx, "reindex-data", nil, reindexDataJob(config))
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}

				ok, err = jobRunner.RunJob(ctx, "delete-old-index", nil, deleteIndexJob(config, config.Stack.Name))
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}

				return true, nil
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

func initMappingJob(config modules.ReconciliationConfig) func(t *batchv1.Job) {
	return func(t *batchv1.Job) {
		imageVersion := config.Versions.Spec.Search

		// notes(gfyrag): this is the first version where the command is available
		// this code will evolved if the mapping change, but it is not planned today
		if config.Versions.IsLower("search", "v0.8.0") {
			imageVersion = "v0.8.0"
		}
		t.Spec = batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{{
						Name:  "init-mapping",
						Image: "ghcr.io/formancehq/search:" + imageVersion,
						Args:  []string{"init-mapping"},
						Env:   searchEnvVars(config).ToCoreEnv(),
					}},
				},
			},
		}
	}
}

func reindexDataJob(config modules.ReconciliationConfig) func(t *batchv1.Job) {
	return func(t *batchv1.Job) {
		t.Spec = batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{{
						Name:  "search-reindex-data",
						Image: "curlimages/curl",
						Command: modules.ShellCommand(`
							set -x
							curl -X POST \
								-H 'Content-Type: application/json' \
								-d '{ "source": { "index": "%s" }, "dest": { "index": "%s" }, "script": { "source": "ctx._source.stack = ctx._index; ctx._id = ctx._index + '"'"'-'"'"' + ctx._id", "lang": "painless" }}' \
								${OPEN_SEARCH_SCHEME}://${OPEN_SEARCH_SERVICE}/_reindex?refresh`, config.Stack.Name, stackv1beta3.DefaultESIndex),
						Env: searchEnvVars(config).ToCoreEnv(),
					}},
				},
			},
		}
	}
}

func deleteIndexJob(config modules.ReconciliationConfig, name string) func(t *batchv1.Job) {
	return func(t *batchv1.Job) {
		t.Spec = batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{{
						Name:  "delete-index",
						Image: "curlimages/curl",
						Command: modules.ShellCommand(`
							curl -X DELETE -H 'Content-Type: application/json' ${OPEN_SEARCH_SCHEME}://${OPEN_SEARCH_SERVICE}/%s`, name),
						Env: searchEnvVars(config).ToCoreEnv(),
					}},
				},
			},
		}
	}
}

func reindexCron(ctx modules.ReconciliationConfig) []modules.Cron {
	return []modules.Cron{
		{
			Container: modules.Container{
				Command: modules.ShellCommand(`
					curl http://search-benthos.%s.svc.cluster.local:4195/ledger_reindex_all -X POST -H 'Content-Type: application/json' -d '{}'`, ctx.Stack.Name),
				Image: "curlimages/curl:8.2.1",
				Name:  "reindex-ledger",
			},
			Schedule: "* * * * *",
			Suspend:  true,
		},
	}
}

func searchEnvVars(rc modules.ReconciliationConfig) modules.ContainerEnv {
	env := elasticSearchEnvVars(rc.Stack, rc.Configuration, rc.Versions).
		Append(
			modules.Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d%s",
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Host,
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Port,
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.PathPrefix)),
			modules.Env("OPEN_SEARCH_SCHEME", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Scheme),
			modules.Env("MAPPING_INIT_DISABLED", "true"),
		)
	if rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
		env = env.Append(
			modules.Env("OPEN_SEARCH_USERNAME", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
			modules.Env("OPEN_SEARCH_PASSWORD", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
		)
	}
	if rc.Versions.IsLower("search", "v0.7.0") {
		env = env.Append(modules.Env("ES_INDICES", rc.Stack.Name))
	} else {
		env = env.Append(modules.Env("ES_INDICES", stackv1beta3.DefaultESIndex))
	}

	return env
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
	return &modules.Service{
		Name: "benthos",
		Port: 4195,
		ExposeHTTP: &modules.ExposeHTTP{
			Name: "benthos",
		},
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
		)
		if configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName == "" {
			ret = ret.Append(
				modules.Env("BASIC_AUTH_USERNAME", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
				modules.Env("BASIC_AUTH_PASSWORD", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
			)
		} else {
			ret = ret.Append(
				modules.EnvFromSecret("BASIC_AUTH_USERNAME", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName, "username"),
				modules.EnvFromSecret("BASIC_AUTH_PASSWORD", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName, "password"),
			)
		}
	}
	return ret
}
