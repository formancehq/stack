package cmd

import (
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Short:   "Run webhooks server",
		RunE:    serve,
	}
}

func serve(cmd *cobra.Command, _ []string) error {
	retriesSchedule := make([]time.Duration, 0)
	rs := viper.GetString(flag.RetriesSchedule)
	if len(rs) > 2 {
		rs = rs[1 : len(rs)-1]
		ss := strings.Split(rs, ",")
		for _, s := range ss {
			d, err := time.ParseDuration(s)
			if err != nil {
				return errors.Wrap(err, "parsing retries schedule duration")
			}
			if d < time.Second {
				return ErrScheduleInvalid
			}
			retriesSchedule = append(retriesSchedule, d)
		}
	}

	options := []fx.Option{
		fx.Provide(func() server.ServiceInfo {
			return server.ServiceInfo{
				Version: Version,
			}
		}),
		postgres.NewModule(viper.GetString(flag.StoragePostgresConnString)),
		otlp.HttpClientModule(),
		server.StartModule(viper.GetString(flag.Listen)),
	}

	if viper.GetBool(flag.Worker) {
		options = append(options, worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule,
		))
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
