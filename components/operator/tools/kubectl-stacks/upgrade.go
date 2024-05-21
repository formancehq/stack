package main

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"time"
)

func NewUpgradeCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	ret := &cobra.Command{
		Use:   "upgrade",
		Short: "Utility to control operator upgrade",
		Long: `
			Upgrading operator with a lot of stacks can stress the k8s cluster.
			This utility help to upgrade stacks smoothly.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			return upgrade(cmd, client)
		},
	}
	ret.Flags().Bool("max", false, "Maximum number of stacks to concurrently upgrade")

	return ret
}

func upgrade(cmd *cobra.Command, client *rest.RESTClient) error {
	stackList, err := lockAllStacks(cmd, client)
	if err != nil {
		return err
	}
	defer func() {
		if err := unlockAllStacks(cmd, client); err != nil {
			panic(err)
		}
	}()
	for _, stack := range stackList.Items {
		pterm.Printfln(`Waiting for stack '%s' to be marked as skipped`, stack.Name)
		if err := waitForStackSkipped(cmd, client, stack.Name, metav1.ConditionTrue); err != nil {
			return err
		}
	}

	pterm.Println(`
All stacks has been locked, you can now upgrade the operator.
Press enter when you are ready`)
	result, err := pterm.DefaultInteractiveConfirm.Show()
	if err != nil {
		return err
	}
	if result {
		pterm.Println(`Will now upgrade stacks one by one`)
		for _, stack := range stackList.Items {
			pterm.Printfln(`Upgrade stack '%s'`, stack.Name)
			if err := unlockStack(cmd, client, stack.Name); err != nil {
				return err
			}

			pterm.Printfln(`Waiting for stack '%s' to be unmarked as skipped`, stack.Name)
			if err := waitForStackSkipped(cmd, client, stack.Name, metav1.ConditionFalse); err != nil {
				return err
			}
		}
	}

	return nil
}

func waitForStackSkipped(cmd *cobra.Command, client *rest.RESTClient, name string, status metav1.ConditionStatus) error {
	for {
		stack := &v1beta1.Stack{}
		if err := client.Get().
			Resource("Stacks").
			Name(name).
			Do(cmd.Context()).
			Into(stack); err != nil {
			return err
		}
		skippedCondition := stack.Status.Conditions.Get("Skipped")
		if status == metav1.ConditionTrue {
			if skippedCondition != nil && skippedCondition.Status == status {
				return nil
			}
		} else {
			if skippedCondition == nil || skippedCondition.Status == status {
				return nil
			}
		}

		select {
		case <-cmd.Context().Done():
			return cmd.Context().Err()
		case <-time.After(time.Second):
		}
	}
}
