package cmd

import (
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewUpdateMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-mapping",
		Short: "Update ElasticSearch mapping",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := newConfig(cmd)
			if err != nil {
				return err
			}

			client, err := newOpensearchClient(cmd, config)
			if err != nil {
				return err
			}

			esIndex, _ := cmd.Flags().GetString(esIndicesFlag)
			if esIndex == "" {
				return errors.New("es index not defined")
			}

			return searchengine.UpdateMapping(cmd.Context(), client, esIndex)
		},
	}

	cmd.Flags().Bool(awsIAMEnabledFlag, false, "Enable AWS IAM")
	iam.AddFlags(cmd.Flags())
	service.AddFlags(cmd.Flags())

	return cmd
}
