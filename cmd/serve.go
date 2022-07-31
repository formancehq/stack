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
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	issuerFlag      = "issuer"
	postgresUriFlag = "postgres-uri"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		issuer := viper.GetString(issuerFlag)
		if issuer == "" {
			return errors.New("issuer has to be defined")
		}
		app := fx.New(
			fx.Supply(fx.Annotate(cmd.Context(), fx.As(new(context.Context)))),
			fx.Supply(&delegatedauth.OAuth2Config{
				//TODO: Make configurable
				ClientID:     "gateway",
				ClientSecret: "ZXhhbXBsZS1hcHAtc2VjcmV0",
				Endpoint: oauth2.Endpoint{
					TokenURL: "http://dex:5556/dex/token",
				},
				RedirectURL: "http://127.0.0.1:8080/authorize/callback",
			}),
			api.Module(issuer, ":8080"),
			storage.Module(viper.GetString(postgresUriFlag)),
			delegatedauth.Module("http://127.0.0.1:5556/dex"),
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
								AccessTokenType:       op.AccessTokenTypeBearer,
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
	serveCmd.Flags().StringP(issuerFlag, "i", "http://localhost:8080", "OIDC issuer")
	serveCmd.Flags().String(postgresUriFlag, "", "Postgres uri")
}
