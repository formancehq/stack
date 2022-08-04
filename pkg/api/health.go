package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/numary/auth/pkg/delegatedauth"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
)

func delegatedOIDCServerAvailability(issuer delegatedauth.Issuer) sharedhealth.NamedCheck {
	return sharedhealth.NewNamedCheck("Delegated OIDC server", sharedhealth.CheckFn(func(ctx context.Context) error {
		rsp, err := http.Get(fmt.Sprintf("%s/.well-known/openid-configuration", issuer))
		if err != nil {
			return err
		}
		if rsp.Body != nil {
			rsp.Body.Close()
		}
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected status code: %d", rsp.StatusCode)
		}
		return nil
	}))
}
