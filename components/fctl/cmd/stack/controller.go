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

func waitStackReady(cmd *cobra.Command, client *membershipclient.APIClient, profile *fctl.Profile, stack *membershipclient.Stack) error {
	waitTime := time.Second
	for {
		stack, resp, err := client.DefaultApi.GetStack(cmd.Context(), stack.OrganizationId, stack.Id).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("stack %s not found", stack.Data.Id)
		}

		if stack.Data.Status == "READY" {
			return nil
		}

		if waitTime > time.Minute {
			pterm.Warning.Println("Waiting for stack to be ready...")
			return fmt.Errorf("there might a problem with the stack scheduling, retry and if the problem persists please contact the support")
		}
		select {
		case <-time.After(waitTime):
		case <-cmd.Context().Done():
			return cmd.Context().Err()
		}
		waitTime *= 2

	}
}
