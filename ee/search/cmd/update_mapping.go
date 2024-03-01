package cmd

import (
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUpdateMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-mapping",
		Short: "Update ElasticSearch mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			openSearchServiceHost := viper.GetString(openSearchServiceFlag)
			if openSearchServiceHost == "" {
				exitWithError(cmd.Context(), "missing open search service host")
			}

			config, err := newConfig(openSearchServiceHost)
			if err != nil {
				return err
			}

			client, err := newOpensearchClient(config)
			if err != nil {
				return err
			}

			esIndex := viper.GetString(esIndicesFlag)
			if esIndex == "" {
				return errors.New("es index not defined")
			}

			return searchengine.UpdateMapping(cmd.Context(), client, esIndex)
		},
	}

	cmd.Flags().Bool(awsIAMEnabledFlag, false, "Enable AWS IAM")
	iam.InitFlags(cmd.Flags())

	return cmd
}
