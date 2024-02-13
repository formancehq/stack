package stack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func waitStackReady(cmd *cobra.Command, client *membershipclient.APIClient, profile *fctl.Profile, stack *membershipclient.Stack) (*membershipclient.Stack, error) {
	var resp *http.Response
	var err error
	var stackRsp *membershipclient.CreateStackResponse

	waitTime := 10 * time.Second
	for {
		stackRsp, resp, err = client.DefaultApi.GetStack(cmd.Context(), stack.OrganizationId, stack.Id).Execute()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("stack %s not found", stack.Id)
		}

		if stack.Status == "READY" {
			return stackRsp.Data, nil
		}

		if waitTime > 2*time.Minute {
			pterm.Warning.Println("Waiting for stack to be ready...")
			return nil, fmt.Errorf("there might a problem with the stack scheduling, retry and if the problem persists please contact the support")
		}
		select {
		case <-time.After(waitTime):
		case <-cmd.Context().Done():
			return nil, cmd.Context().Err()
		}
		waitTime *= 2

	}
}
