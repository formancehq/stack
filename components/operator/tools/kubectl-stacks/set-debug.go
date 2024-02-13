package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

func NewSetDebugCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:  "set-debug <stack-name> on|off",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			return setDebug(cmd, client, args[0], args[1] == "on")
		},
	}
}

func setDebug(cmd *cobra.Command, client *rest.RESTClient, name string, b bool) error {
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Setting debug on stack '%s'...\r\n", name)
	content, err := json.Marshal(map[string]any{
		"spec": map[string]any{
			"debug": b,
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
