package suite

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/formancehq/go-libs/collectionutils"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
)

type claims struct {
	jwt.StandardClaims
	Scopes string `json:"scope"`
}

func forgeSecurityToken(scopes ...string) string {
	claims := claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  modules.AuthIssuer,
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			Issuer:    OIDCServer().Issuer(),
		},
		Scopes: strings.Join(scopes, " "),
	}
	signJWT, err := OIDCServer().Keypair.SignJWT(claims)
	Expect(err).To(BeNil())

	return signJWT
}

func exchangeSecurityToken(securityToken string, scopes ...string) *oauth2.Token {
	scopes = append(scopes, "email")
	form := url.Values{
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  []string{securityToken},
		"scope":      []string{strings.Join(scopes, " ")},
	}

	req, err := http.NewRequestWithContext(TestContext(), http.MethodPost, GatewayURL()+"/api/auth/oauth/token",
		bytes.NewBufferString(form.Encode()))
	Expect(err).To(BeNil())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err := HTTPClient().Do(req)
	Expect(err).To(BeNil())
	Expect(ret.StatusCode).To(Equal(http.StatusOK))

	stackToken := &oauth2.Token{}
	Expect(json.NewDecoder(ret.Body).Decode(stackToken)).To(Succeed())

	return stackToken
}

var _ = WithModules([]*Module{modules.Auth}, func() {
	var (
		securityToken string
	)
	BeforeEach(func() {
		securityToken = forgeSecurityToken("openid scope1")
	})
	When("exchanging security token against an access token", func() {
		var (
			token *oauth2.Token
		)
		BeforeEach(func() {
			token = exchangeSecurityToken(securityToken, "other_scope1 other_scope2")
		})
		It("should be ok, even if wrong scope are asked", func() {
			accessTokenClaims := &oidc.AccessTokenClaims{}
			_, err := oidc.ParseToken(token.AccessToken, accessTokenClaims)
			Expect(err).To(Succeed())

			Expect(accessTokenClaims.Scopes).To(HaveLen(2))
			Expect(collectionutils.Contains(accessTokenClaims.Scopes, "scope1")).To(BeTrue())
			Expect(collectionutils.Contains(accessTokenClaims.Scopes, "openid")).To(BeTrue())
		})
	})
})
