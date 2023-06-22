package login

import (
	"context"
	"net/url"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

func LogIn(ctx context.Context, dialog Dialog, relyingParty rp.RelyingParty) (*oidc.AccessTokenResponse, error) {
	deviceCode, err := rp.DeviceAuthorization(relyingParty.OAuthConfig().Scopes, relyingParty)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(deviceCode.VerificationURI)
	if err != nil {
		panic(err)
	}
	query := uri.Query()
	query.Set("user_code", deviceCode.UserCode)
	uri.RawQuery = query.Encode()

	dialog.DisplayURIAndCode(deviceCode.VerificationURI, deviceCode.UserCode)

	if err := fctl.Open(uri.String()); err != nil {
		return nil, err
	}

	return rp.DeviceAccessToken(ctx, deviceCode.DeviceCode, time.Duration(deviceCode.Interval)*time.Second, relyingParty)
}
