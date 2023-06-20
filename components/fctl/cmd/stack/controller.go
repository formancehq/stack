package stack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func waitStackReady(cmd *cobra.Command, profile *fctl.Profile, stack *membershipclient.Stack) error {
	baseUrlStr := profile.ServicesBaseUrl(stack).String()
	authServerUrl := fmt.Sprintf("%s/api/auth", baseUrlStr)
	for {
		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet,
			fmt.Sprintf(authServerUrl+"/.well-known/openid-configuration"), nil)
		if err != nil {
			return err
		}
		rsp, err := fctl.GetHttpClient(cmd, map[string][]string{}).Do(req)
		if err == nil && rsp.StatusCode == http.StatusOK {
			break
		}
		select {
		case <-cmd.Context().Done():
			return cmd.Context().Err()
		case <-time.After(time.Second):
		}
	}
	return nil
}
