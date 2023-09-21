package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path/filepath"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"gopkg.in/yaml.v3"
)

func init() {
	modules.Register("auth", modules.Module{
		Postgres: func(ctx modules.Context) stackv1beta3.PostgresConfig {
			return ctx.Configuration.Spec.Services.Auth.Postgres
		},
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
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
		},
	})
}

func resolveAuthContainer(resolveContext modules.ContainerResolutionContext) modules.Container {
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
		Resources: getResourcesWithDefault(
			resolveContext.Configuration.Spec.Services.Auth.ResourceProperties,
			modules.ResourceSizeSmall(),
		),
	}
}

func resolveAuthSecrets(resolveContext modules.ServiceInstallContext) modules.Secrets {
	return modules.Secrets{
		"secret": modules.Secret{
			Data: map[string][]byte{
				"signingKey": []byte(RSAKeyGenerator()),
			},
		},
	}
}

func resolveAuthConfigs(resolveContext modules.ServiceInstallContext) modules.Configs {
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
