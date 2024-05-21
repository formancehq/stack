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

func NewLockCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:  "lock [<stack-name>]",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				_, err := lockAllStacks(cmd, client)
				return err
			} else {
				return lockStack(cmd, client, args[0])
			}
		},
	}
}

func lockAllStacks(cmd *cobra.Command, client *rest.RESTClient) (*v1beta1.StackList, error) {
	list := &v1beta1.StackList{}
	if err := client.Get().
		Resource("Stacks").
		Do(cmd.Context()).
		Into(list); err != nil {
		return nil, err
	}

	for _, stack := range list.Items {
		if err := lockStack(cmd, client, stack.Name); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func lockStack(cmd *cobra.Command, client *rest.RESTClient, name string) error {

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Locking stack '%s'...\r\n", name)

	content, err := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"annotations": map[string]any{
				v1beta1.SkipLabel: "true",
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
