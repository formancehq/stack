package main

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"text/tabwriter"
	"time"
)

func NewListCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			list := &v1beta1.StackList{}
			if err := client.Get().
				Resource("Stacks").
				Do(cmd.Context()).
				Into(list); err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)
			_, _ = fmt.Fprintln(w, "Name\tAge\tReady")
			for _, stack := range list.Items {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\r\n", stack.Name,
					time.Since(stack.CreationTimestamp.Time), func() string {
						if stack.Status.Ready {
							return "yes"
						}
						return "no"
					}())
			}
			return w.Flush()
		},
	}
}
