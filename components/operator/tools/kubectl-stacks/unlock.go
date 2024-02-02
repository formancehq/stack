package main

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

func NewUnlockCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:  "unlock [<stack-name>]",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return unlockAllStacks(cmd, client)
			} else {
				return unlockStack(cmd, client, args[0])
			}
		},
	}
}

func unlockAllStacks(cmd *cobra.Command, client *rest.RESTClient) error {
	list := &v1beta1.StackList{}
	if err := client.Get().
		Resource("Stacks").
		Do(cmd.Context()).
		Into(list); err != nil {
		return err
	}

	for _, stack := range list.Items {
		if err := unlockStack(cmd, client, stack.Name); err != nil {
			return err
		}
	}
	return nil
}

func unlockStack(cmd *cobra.Command, client *rest.RESTClient, name string) error {
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Unlocking stack '%s'...\r\n", name)
	content, err := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"annotations": map[string]any{
				v1beta1.SkipLabel: nil,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return client.Patch(types.MergePatchType).
		Resource("Stacks").
		Name(name).
		Body(content).
		Do(cmd.Context()).
		Error()
}
