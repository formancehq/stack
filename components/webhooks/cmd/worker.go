package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Run webhooks worker",
		RunE:  runWorker,
	}
}

func runWorker(cmd *cobra.Command, _ []string) error {
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

	return service.New(
		cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		postgres.NewModule(viper.GetString(flag.StoragePostgresConnString)),
		fx.Provide(worker.NewWorkerHandler),
		fx.Invoke(func(lc fx.Lifecycle, h http.Handler) {
			lc.Append(httpserver.NewHook(h, httpserver.WithAddress(viper.GetString(flag.Listen))))
		}),
		worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule,
		),
	).Run(cmd.Context())
}
