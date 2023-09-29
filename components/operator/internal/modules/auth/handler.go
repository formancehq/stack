package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path/filepath"

	"github.com/formancehq/operator/internal/modules/control"
	"github.com/formancehq/operator/internal/modules/orchestration"
	"github.com/formancehq/operator/internal/modules/wallets"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"gopkg.in/yaml.v3"
)

type module struct{}

func (a module) DependsOn() []modules.Module {
	return []modules.Module{
		control.Module,
		orchestration.Module,
		wallets.Module,
	}
}

func (a module) Name() string {
	return "auth"
}

func (a module) Postgres(ctx modules.ReconciliationConfig) stackv1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Auth.Postgres
}

func (a module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{{
					Secured:                 true,
					ListenEnvVar:            "LISTEN",
					ExposeHTTP:              true,
					Configs:                 resolveAuthConfigs,
					Secrets:                 resolveAuthSecrets,
					Container:               resolveAuthContainer,
					InjectPostgresVariables: true,
					HasVersionEndpoint:      true,
					Annotations:             ctx.Configuration.Spec.Services.Auth.Annotations.Service,
				}}
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.PostgresAwareModule = Module
var _ modules.DependsOnAwareModule = Module

func init() {
	modules.Register(Module)
}

func resolveAuthContainer(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
	env := modules.ContainerEnv{
		modules.Env("CONFIG", filepath.Join(resolveContext.GetConfig("config").GetMountPath(), "config.yaml")),
		modules.Env("DELEGATED_CLIENT_SECRET", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientSecret),
		modules.Env("DELEGATED_CLIENT_ID", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientID),
		modules.Env("DELEGATED_ISSUER", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.Issuer),
		modules.Env("BASE_URL", fmt.Sprintf("%s://%s/api/auth", resolveContext.Stack.Spec.Scheme, resolveContext.Stack.Spec.Host)),
		modules.EnvFromSecret("SIGNING_KEY", resolveContext.GetSecret("secret").GetName(), "signingKey"),
	}
	if resolveContext.Stack.Spec.Dev {
		env = env.Append(modules.Env("CAOS_OIDC_DEV", "1"))
	}
	return modules.Container{
		Args:  []string{"serve"},
		Env:   env,
		Image: modules.GetImage("auth", resolveContext.Versions.Spec.Auth),
		Resources: modules.GetResourcesWithDefault(
			resolveContext.Configuration.Spec.Services.Auth.ResourceProperties,
			modules.ResourceSizeSmall(),
		),
	}
}

func resolveAuthSecrets(resolveContext modules.ServiceInstallConfiguration) modules.Secrets {
	return modules.Secrets{
		"secret": modules.Secret{
			Data: map[string][]byte{
				"signingKey": []byte(RSAKeyGenerator()),
			},
		},
	}
}

func resolveAuthConfigs(resolveContext modules.ServiceInstallConfiguration) modules.Configs {
	yaml, err := yaml.Marshal(struct {
		Clients []stackv1beta3.StaticClient `yaml:"clients"`
	}{
		Clients: resolveContext.Stack.GetStaticClients(resolveContext.Configuration),
	})
	if err != nil {
		panic(err)
	}
	return modules.Configs{
		"config": modules.Config{
			Data: map[string]string{
				"config.yaml": string(yaml),
			},
			Mount: true,
		},
	}
}

var (
	RSAKeyGenerator = func() string {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		var privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
		buf := bytes.NewBufferString("")
		err = pem.Encode(buf, privateKeyBlock)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
)
