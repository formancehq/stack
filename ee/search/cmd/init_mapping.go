package cmd

import (
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInitMapping() *cobra.Command {
	return &cobra.Command{
		Use:   "init-mapping",
		Short: "Init ElasticSearch mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			openSearchServiceHost := viper.GetString(openSearchServiceFlag)
			if openSearchServiceHost == "" {
				exitWithError(cmd.Context(), "missing open search service host")
			}

			client, err := newOpensearchClient(openSearchServiceHost)
			if err != nil {
				return err
			}

			esIndex := viper.GetString(esIndicesFlag)
			if esIndex == "" {
				return errors.New("es index not defined")
			}

			return searchengine.CreateIndex(cmd.Context(), client, esIndex)
		},
	}
}
