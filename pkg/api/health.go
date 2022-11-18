package api

import (
	"context"
	"net/http"

	sharedhealth "github.com/formancehq/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/client"
	"github.com/zitadel/oidc/pkg/client/rp"
)

func delegatedOIDCServerAvailability(rp rp.RelyingParty) sharedhealth.NamedCheck {
	return sharedhealth.NewNamedCheck("Delegated OIDC server", sharedhealth.CheckFn(func(ctx context.Context) error {
		_, err := client.Discover(rp.Issuer(), http.DefaultClient)
		return err
	}))
}
