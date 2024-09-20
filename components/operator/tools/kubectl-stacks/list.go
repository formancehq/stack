package main

import (
	"fmt"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewListCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	ret := &cobra.Command{
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

			stacks := list.Items
			onlyEnabled, err := cmd.Flags().GetBool("only-enabled")
			if err != nil {
				return err
			}
			if onlyEnabled {
				stacks = collectionutils.Filter(stacks, func(stack v1beta1.Stack) bool {
					return !stack.Spec.Disabled
				})
			}

			sort.Slice(stacks, func(i, j int) bool {
				return stacks[i].CreationTimestamp.Before(&stacks[j].CreationTimestamp)
			})

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)
			_, _ = fmt.Fprintln(w, "Name\tCreated at\tReady\tLocked?")
			for _, stack := range stacks {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\r\n", stack.Name,
					stack.CreationTimestamp.Time.Format(time.RFC3339), func() string {
						if stack.Status.Ready {
							return "yes"
						}
						return "no"
					}(),
					func() string {
						if stack.GetAnnotations()[v1beta1.SkipLabel] == "true" {
							return "yes"
						}
						return ""
					}(),
				)
			}
			return w.Flush()
		},
	}

	ret.Flags().Bool("only-enabled", false, "Filter only enabled stacks")

	return ret
}
