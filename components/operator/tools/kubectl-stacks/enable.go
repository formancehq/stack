package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

func NewEnableCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:  "enable <stack-name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			return enable(cmd, client, args[0])
		},
	}
}

func enable(cmd *cobra.Command, client *rest.RESTClient, name string) error {
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Enable stack '%s'...\r\n", name)
	content, err := json.Marshal(map[string]any{
		"spec": map[string]any{
			"disabled": false,
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
