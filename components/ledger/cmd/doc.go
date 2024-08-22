package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "doc",
	}
	cmd.AddCommand(NewDocFlagCommand())
	return cmd
}

func NewDocFlagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "flags",
		Run: func(cmd *cobra.Command, args []string) {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			defer func(w *tabwriter.Writer) {
				if err := w.Flush(); err != nil {
					panic(err)
				}
			}(w)

			allKeys := make([]string, 0)
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				allKeys = append(allKeys, f.Name)
			})
			sort.Strings(allKeys)

			if _, err := fmt.Fprintf(w,
				"\tFlag\tEnv var\tDefault value\tDescription\t\r\n"); err != nil {
				panic(err)
			}
			if _, err := fmt.Fprintf(w,
				"\t-\t-\t-\t-\t\r\n"); err != nil {
				panic(err)
			}
			for _, key := range allKeys {
				asEnvVar := strings.ToUpper(strings.Replace(key, "-", "_", -1))
				flag := cmd.Parent().Parent().PersistentFlags().Lookup(key)
				if flag == nil {
					continue
				}
				if _, err := fmt.Fprintf(w,
					"\t--%s\t%s\t%s\t%s\t\r\n", key, asEnvVar, flag.DefValue, flag.Usage); err != nil {
					panic(err)
				}
			}
		},
	}
	return cmd
}
