package fctl

import (
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
)

func GetAuthRelyingParty(httpClient *http.Client, membershipURI string) (rp.RelyingParty, error) {
	return rp.NewRelyingPartyOIDC(membershipURI, AuthClient, "",
		"", []string{"openid", "email", "offline_access", "supertoken", "accesses"}, rp.WithHTTPClient(httpClient))
}
