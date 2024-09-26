package suite

import (
	"fmt"
	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/collectionutils"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

var _ = WithModules([]*Module{modules.Auth}, func() {
	With("the default static client", func() {
		var (
			httpClient *http.Client
			sdkClient  *formance.Formance
		)
		BeforeEach(func() {
			config := clientcredentials.Config{
				ClientID:     "global",
				ClientSecret: "global",
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", GatewayURL()),
			}
			httpClient = config.Client(TestContext())
			sdkClient = Client(formance.WithClient(httpClient))
		})
		When("creating a new brand client", func() {
			var (
				createClientResponse *operations.CreateClientResponse
				createSecretResponse *operations.CreateSecretResponse
				err                  error
			)
			BeforeEach(func() {
				createClientResponse, err = sdkClient.Auth.V1.CreateClient(TestContext(), &shared.CreateClientRequest{
					Name:   "client1",
					Scopes: []string{"scope1"},
				})
				Expect(err).To(Succeed())
				Expect(createClientResponse.StatusCode).To(Equal(http.StatusCreated))

				createSecretResponse, err = sdkClient.Auth.V1.CreateSecret(TestContext(), operations.CreateSecretRequest{
					CreateSecretRequest: &shared.CreateSecretRequest{
						Name: "secret1",
					},
					ClientID: createClientResponse.CreateClientResponse.Data.ID,
				})
				Expect(err).To(Succeed())
			})
			When("requiring an access token using client credentials flow", func() {
				var (
					token *oauth2.Token
					err   error
				)
				BeforeEach(func() {
					config := clientcredentials.Config{
						ClientID:     createClientResponse.CreateClientResponse.Data.ID,
						ClientSecret: createSecretResponse.CreateSecretResponse.Data.Clear,
						TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", GatewayURL()),
						Scopes:       []string{"scope1", "scope2"},
					}
					token, err = config.Token(TestContext())
					Expect(err).To(BeNil())
					Expect(token).NotTo(BeNil())
				})
				It("should be ok", func() {
					accessTokenClaims := &oidc.AccessTokenClaims{}
					_, err = oidc.ParseToken(token.AccessToken, accessTokenClaims)
					Expect(err).To(Succeed())

					Expect(accessTokenClaims.Scopes).To(HaveLen(1))
					Expect(collectionutils.Contains(accessTokenClaims.Scopes, "scope1")).To(BeTrue())
				})
			})
		})
	})
})
