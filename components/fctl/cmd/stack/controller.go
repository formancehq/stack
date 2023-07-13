package stack

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
)

func waitStackReady(ctx context.Context, flags *flag.FlagSet, profile *fctl.Profile, stack *membershipclient.Stack) error {
	baseUrlStr := profile.ServicesBaseUrl(stack).String()
	authServerUrl := fmt.Sprintf("%s/api/auth", baseUrlStr)
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			fmt.Sprintf(authServerUrl+"/.well-known/openid-configuration"), nil)
		if err != nil {
			return err
		}
		rsp, err := fctl.GetHttpClient(flags, map[string][]string{}).Do(req)
		if err == nil && rsp.StatusCode == http.StatusOK {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return nil
}
