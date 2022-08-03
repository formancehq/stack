package cmd

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/numary/auth/pkg/api"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	postgresUriFlag           = "postgres-uri"
	delegatedClientIDFlag     = "delegated-client-id"
	delegatedClientSecretFlag = "delegated-client-secret"
	delegatedIssuerFlag       = "delegated-issuer"
	baseUrlFlag               = "base-url"
	signingKeyFlag            = "signing-key"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		baseUrl := viper.GetString(baseUrlFlag)
		if baseUrl == "" {
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
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}

		app := fx.New(
			fx.Supply(fx.Annotate(cmd.Context(), fx.As(new(context.Context)))),
			api.Module(baseUrl, ":8080"),
			storage.Module(viper.GetString(postgresUriFlag), key),
			delegatedauth.Module(
				delegatedIssuer,
				delegatedClientID,
				delegatedClientSecret,
				fmt.Sprintf("%s/delegatedoidc/callback", baseUrl),
			),
			sharedotlptraces.CLITracesModule(viper.GetViper()),
			fx.Invoke(func() {
				sharedlogging.Infof("App started.")
			}),
			fx.NopLogger,
		)
		err = app.Start(cmd.Context())
		if err != nil {
			return err
		}
		<-app.Done()

		return app.Err()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String(postgresUriFlag, "", "Postgres uri")
	serveCmd.Flags().String(delegatedIssuerFlag, "", "Delegated OIDC issuer")
	serveCmd.Flags().String(delegatedClientIDFlag, "", "Delegated OIDC client id")
	serveCmd.Flags().String(delegatedClientSecretFlag, "", "Delegated OIDC client secret")
	serveCmd.Flags().String(baseUrlFlag, "http://localhost:8080", "Base service url")
	serveCmd.Flags().String(signingKeyFlag, "", "Signing key")

	sharedotlptraces.InitOTLPTracesFlags(serveCmd.Flags())
}
