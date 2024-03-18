package stack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func waitStackReady(cmd *cobra.Command, client *fctl.MembershipClient, organizationId, stackId string) (*membershipclient.Stack, error) {
	var resp *http.Response
	var err error
	var stackRsp *membershipclient.CreateStackResponse

	waitTime := 2 * time.Second
	sum := 2 * time.Second

	// Hack to ignore first Status
	select {
	case <-cmd.Context().Done():
		return nil, cmd.Context().Err()
	case <-time.After(waitTime):
	}

	for {
		err = client.RefreshIfNeeded(cmd)
		if err != nil {
			return nil, err
		}

		stackRsp, resp, err = client.DefaultApi.GetStack(cmd.Context(), organizationId, stackId).Execute()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("stack %s not found", stackId)
		}

		if stackRsp.Data.Status == "READY" {
			return stackRsp.Data, nil
		}

		if sum > 10*time.Minute {
			pterm.Warning.Printf("You can check fctl stack show %s --organization %s to see the status of the stack", stackId, organizationId)
			problem := fmt.Errorf("there might a problem with the stack scheduling, if the problem persists, please contact the support")

			err = internal.PrintStackInformation(cmd.OutOrStdout(), client.GetProfile(), stackRsp.Data, nil)
			if err != nil {
				return nil, problem
			}

			return nil, problem
		}

		sum += waitTime
		select {
		case <-time.After(waitTime):
		case <-cmd.Context().Done():
			return nil, cmd.Context().Err()
		}
	}
}
