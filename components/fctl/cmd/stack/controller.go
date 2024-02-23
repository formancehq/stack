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

func waitStackReady(cmd *cobra.Command, client *membershipclient.APIClient, profile *fctl.Profile, stackId, organizationId string) (*membershipclient.Stack, error) {
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

		if sum > 2*time.Minute {
			pterm.Warning.Println("Waiting for stack to be ready...")
			return nil, fmt.Errorf("there might a problem with the stack scheduling, retry and if the problem persists please contact the support")
		}

		sum += waitTime
		select {
		case <-time.After(waitTime):
		case <-cmd.Context().Done():
			return nil, cmd.Context().Err()
		}
	}
}
