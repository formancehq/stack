package main

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"os"
)

// NewRootCommand provides a cobra command wrapping NamespaceOptions
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
	}

	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(cmd.Flags())
	cmd.AddCommand(
		NewLockCommand(configFlags),
		NewUnlockCommand(configFlags),
		NewListCommand(configFlags),
		NewSetDebugCommand(configFlags),
	)

	return cmd
}

func getRestClient(configFlags *genericclioptions.ConfigFlags) (*rest.RESTClient, error) {
	restConfig, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	restConfig.GroupVersion = &v1beta1.GroupVersion
	restConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	restConfig.APIPath = "/apis"

	return rest.RESTClientFor(restConfig)
}

func main() {
	flags := pflag.NewFlagSet("kubectl-stacks", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
}
