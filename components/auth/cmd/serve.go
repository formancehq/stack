package cmd

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/otlp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/api"
	"github.com/formancehq/auth/pkg/api/authorization"
	"github.com/formancehq/auth/pkg/delegatedauth"
	"github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	zLogging "github.com/zitadel/logging"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const (
	postgresUriFlag           = "postgres-uri"
	delegatedClientIDFlag     = "delegated-client-id"
	delegatedClientSecretFlag = "delegated-client-secret"
	delegatedIssuerFlag       = "delegated-issuer"
	baseUrlFlag               = "base-url"
	listenFlag                = "listen"
	signingKeyFlag            = "signing-key"
	configFlag                = "config"

	defaultSigningKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAth3atoCXldJgHH9EWnZQMvw5O+vVNKMcvrllEGQsLxvIA5xy
YPnFt2xU7k1dcN5ViBqPiigVHZNeyyHcdVclg26zqjEwYUqH+OPiRFeBn0SwOG+d
ZLpOIJdKt7OjmUG0xN9egq81dbPVPBPckuWqB9XMWmM+dtqydBX4lekj+Q1hFn5E
WXXuAs9aLIc8DzPz8B+oqwLKZ6k6kC5vpj+EaBt8ExVywrWftkWewGWRO7fLw0Fj
7hamaA1ZTYEqCN+MLDLEd6qmtC2cdgVhZM0RG2OnTiq5lGzNFmLXGOsquc35HSQj
OQqcLL+e/72K3giJ1YCqYWAIIJcc/kNKU8HtpwIDAQABAoIBAFY+dSEQbLjq09Er
A/fDJ9+9Sm1yFZnD1Q0NRysoBTSZ93KeWBxMrLFcgCwKP0IASIkX6voGWVmUPMP9
2SVIi99eQX9LpBmu7g2T/cdXmW8PXFSdpu/Yur78ZsnwLH2bfDvvfBZvWuXOsCCv
VznJwWfMe+YiMaafkvsenIaBziNWwUeVGHCWl5f++KGGbWFZjhkRZyjKWfMYflig
EG5e+WaXagCjTah5pUkmvLj3jmB1iGA/Askm8S5QyTt6Z+SIEk+i5T3qCiLFNvzp
7OeSyfbmWWzBYTiSvEoHhaHfdeicUyOpRthc33bb7LnfIWDG3Z+WE0o6U1nR8o7U
t5dsj2ECgYEA7SEuBpd/3wdNVLQSI/RHKKO3sdlymh7yRFf7OAn/UxnSJbSNx4y4
GAEdJD9KwSQlyekLITF+xc0IuyFHOmvuzp1+/LxK/QTY4dcdlwl/r1kmwBbTeR0e
yl9RtulHXmP+Ss/PZgwR081Lk7zlRkh1busyAOmCE4mJW/IvNBze0dsCgYEAxJvy
PcbaLVk497U9cUGznsSbbsyq7JGLkBgTu3eQ/yRgoE7pvagF7dV1gQGuCYjOaYml
U4d95FLPoiE+CE0g2uyouFEsD1UhggTADP33BidUKUcF1ub9VVNcWs4I5LeWPY/X
5vcpOCAkmRZWT5rieAECdIsfRTnePVyn2L7amyUCgYEAqsZAfWLSJm791Eiy383n
CW+OtbjiffhXhbzPIbaheNmZrKnxiYrgcfkrYZVrYtmDlXwOFeOtZwqYhRwcTgi5
PXfTonSAlOPOxibEGqgumrvb2m8V8Z11NU2cbdxnF6Vv17T9qoJ6vEyXZ1iczhcU
68LaiimhEiz1DZDHSgKYvg0CgYEAjVZyQXjXVWxjqKdQ4T9TKhq6hl95rJFA3DiC
zuy4fsKe9/9ixyWoBX7DdxdHDrGbeYErKa4okV/6xdnR51PS/67L55zq6KbRbM+P
ZIeZ8oGJXhchmoj5q0I/DUQ6Xnmf9ueWVQJvTlrFFIxbReTZU12ebzuoIjLkkgYu
34DsVEUCgYEAtHm/aO7/2UJT40PMO+VDvBCEixPtt6j72fLaW8btgVRAnhp9qaWX
Cv6TRZPe2y6Bbgg4Q3FuF0DMqx6ongFKQAWo3DkqNFCGRgjJMQ9JbcfOnGySq4U+
EL/wy5C80pa3jahniqVgO5L6zz0ZLtRIRE7aCtCIu826gctJ1+ShIso=
-----END RSA PRIVATE KEY-----
`
)

func otlpHttpClientModule(debug bool) fx.Option {
	return fx.Provide(func() *http.Client {
		return &http.Client{
			Transport: otlp.NewRoundTripper(http.DefaultTransport, debug, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				str := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
				if len(r.URL.Query()) == 0 {
					return str
				}

				return fmt.Sprintf("%s?%s", str, r.URL.Query().Encode())
			})),
		}
	})
}

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString(baseUrlFlag) == "" {
				return errors.New("base url must be defined")
			}

			delegatedClientID := viper.GetString(delegatedClientIDFlag)
			if delegatedClientID == "" {
				return errors.New("delegated client id must be defined")
			}

			delegatedClientSecret := viper.GetString(delegatedClientSecretFlag)
			if delegatedClientSecret == "" {
				return errors.New("delegated client secret must be defined")
			}

			delegatedIssuer := viper.GetString(delegatedIssuerFlag)
			if delegatedIssuer == "" {
				return errors.New("delegated issuer must be defined")
			}

			signingKey := viper.GetString(signingKeyFlag)
			if signingKey == "" {
				return errors.New("signing key must be defined")
			}

			block, _ := pem.Decode([]byte(signingKey))
			if block == nil {
				return errors.New("invalid signing key, cannot parse as PEM")
			}

			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return err
			}

			if viper.GetString(configFlag) != "" {
				viper.SetConfigFile(viper.GetString(configFlag))
				if err := viper.ReadInConfig(); err != nil {
					return errors.Wrap(err, "reading viper config file")
				}
			}

			type configuration struct {
				Clients []auth.StaticClient `json:"clients" yaml:"clients"`
			}
			o := configuration{}
			if err := viper.Unmarshal(&o); err != nil {
				return errors.Wrap(err, "unmarshal viper config")
			}

			zLogging.SetOutput(cmd.OutOrStdout())

			options := []fx.Option{
				otlpHttpClientModule(viper.GetBool(service.DebugFlag)),
				fx.Supply(fx.Annotate(cmd.Context(), fx.As(new(context.Context)))),
				fx.Supply(delegatedauth.Config{
					Issuer:       delegatedIssuer,
					ClientID:     delegatedClientID,
					ClientSecret: delegatedClientSecret,
					RedirectURL:  fmt.Sprintf("%s/authorize/callback", viper.GetString(baseUrlFlag)),
				}),
				sqlstorage.Module(sqlstorage.KindPostgres, viper.GetString(postgresUriFlag), key, o.Clients...),
				api.Module(viper.GetString(listenFlag), viper.GetString(baseUrlFlag), sharedapi.ServiceInfo{
					Version: Version,
				}),
				oidc.Module(key, viper.GetString(baseUrlFlag), o.Clients...),
				authorization.Module(),
				delegatedauth.Module(),
				fx.Decorate(func(logger logging.Logger) *gorm.Config {
					return &gorm.Config{
						Logger: sqlstorage.NewLogger(logger),
					}
				}),
			}

			options = append(options, otlptraces.CLITracesModule(viper.GetViper()))

			return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
		},
	}

	cmd.Flags().String(postgresUriFlag, "", "Postgres uri")
	cmd.Flags().String(delegatedIssuerFlag, "", "Delegated OIDC issuer")
	cmd.Flags().String(delegatedClientIDFlag, "", "Delegated OIDC client id")
	cmd.Flags().String(delegatedClientSecretFlag, "", "Delegated OIDC client secret")
	cmd.Flags().String(baseUrlFlag, "http://localhost:8080", "Base service url")
	cmd.Flags().String(signingKeyFlag, defaultSigningKey, "Signing key")
	cmd.Flags().String(listenFlag, ":8080", "Listening address")

	cmd.Flags().String(configFlag, "", "Config file name without extension")

	otlptraces.InitOTLPTracesFlags(cmd.Flags())

	return cmd
}
