package cmd

import (
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewInitMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-mapping",
		Short: "Init ElasticSearch mapping",
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

			return searchengine.CreateIndex(cmd.Context(), client, esIndex)
		},
	}
	cmd.Flags().Bool(awsIAMEnabledFlag, false, "Enable AWS IAM")
	cmd.Flags().String(esIndicesFlag, "", "ES index to look")
	cmd.Flags().String(openSearchServiceFlag, "", "Open search service hostname")
	cmd.Flags().String(openSearchSchemeFlag, "https", "OpenSearch scheme")
	cmd.Flags().String(openSearchUsernameFlag, "", "OpenSearch username")
	cmd.Flags().String(openSearchPasswordFlag, "", "OpenSearch password")
	cmd.Flags().String(stackFlag, "", "Stack id")

	iam.AddFlags(cmd.Flags())
	service.AddFlags(cmd.Flags())

	return cmd
}
