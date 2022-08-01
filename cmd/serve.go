package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	auth "github.com/numary/auth/pkg"
	"github.com/numary/auth/pkg/api"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	postgresUriFlag           = "postgres-uri"
	delegatedClientIDFlag     = "delegated-client-id"
	delegatedClientSecretFlag = "delegated-client-secret"
	delegatedIssuerFlag       = "delegated-issuer"
	baseUrlFlag               = "base-url"
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

		app := fx.New(
			fx.Supply(fx.Annotate(cmd.Context(), fx.As(new(context.Context)))),
			api.Module(baseUrl, ":8080"),
			storage.Module(viper.GetString(postgresUriFlag)),
			delegatedauth.Module(
				delegatedIssuer,
				delegatedClientID,
				delegatedClientSecret,
				fmt.Sprintf("%s/authorize/callback", baseUrl),
			),
			fx.Invoke(func() {
				sharedlogging.Infof("App started.")
			}),
			fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return db.
							WithContext(ctx).
							Clauses(clause.OnConflict{
								Columns: []clause.Column{{Name: "id"}},
								DoUpdates: clause.Assignments(map[string]interface{}{
									"grant_types": auth.Array[oidc.GrantType]{
										oidc.GrantTypeCode,
										oidc.GrantTypeRefreshToken,
										oidc.GrantTypeClientCredentials,
									},
									"post_logout_redirect_uri": `["http://localhost:3000/"]`,
									"scopes":                   fmt.Sprintf(`["%s"]`, strings.Join(auth.Scopes, `", "`)),
									"access_token_type":        op.AccessTokenTypeJWT,
								}),
							}).
							Create(&auth.Client{
								Id:     "demo",
								Secret: "1234",
								RedirectURIs: auth.Array[string]{
									"http://localhost:3000/auth-callback",
								},
								ApplicationType: op.ApplicationTypeWeb,
								AuthMethod:      oidc.AuthMethodNone,
								ResponseTypes:   []oidc.ResponseType{oidc.ResponseTypeCode},
								GrantTypes: []oidc.GrantType{
									oidc.GrantTypeCode,
									oidc.GrantTypeRefreshToken,
									oidc.GrantTypeClientCredentials,
								},
								AccessTokenType:       op.AccessTokenTypeJWT,
								PostLogoutRedirectUri: auth.Array[string]{"http://localhost:3000/"},
								Scopes:                auth.Scopes,
							}).Error
					},
				})
			}),
			fx.NopLogger,
		)
		err := app.Start(cmd.Context())
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
}
