package main

import (
	"encoding/json"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewSettingsCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	ret := &cobra.Command{
		Use:  "settings",
	}

	ret.AddCommand(
		NewAddSettingsCommand(configFlags),
		NewRmSettingsCommand(configFlags),
	)

	return ret
}

func  NewAddSettingsCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	ret := &cobra.Command{
		Use: "add <name> <key> <value>",
		Short: "Create a new settings",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			stacks, err := cmd.Flags().GetStringSlice("stacks")
			if err != nil {
				return err
			}

			settings := v1beta1.Settings{
				ObjectMeta: v1.ObjectMeta{
					Name: args[0],
				},
				Spec: v1beta1.SettingsSpec{
					Key: args[1],
					Value: args[2],
					Stacks: stacks,
				},
			}
			settings.SetGroupVersionKind(v1beta1.GroupVersion.WithKind("Settings"))

			data, err := json.Marshal(settings)
			if err != nil {
				return err
			}

			return client.Post().
				Name(args[0]).
				Resource("Settings").
				Body(data).
				Do(cmd.Context()).
				Error()
		},
	}

	ret.Flags().StringSlice("stacks", []string{"*"}, "Select stacks on which the setting is applied")

	return ret
}

func NewRmSettingsCommand(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use: "rm",
		Short: "Delete a setting",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getRestClient(configFlags)
			if err != nil {
				return err
			}

			return client.Delete().
				Name(args[0]).
				Resource("Settings").
				Do(cmd.Context()).
				Error()
		},
	}
}