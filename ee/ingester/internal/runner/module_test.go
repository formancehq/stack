package runner

import (
	"net/http"
	"testing"
	"text/template"
	"time"

	"github.com/formancehq/stack/ee/ingester/internal/modules"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"

	"github.com/formancehq/stack/ee/ingester/internal/httpclient"

	"go.uber.org/mock/gomock"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)

	driversStore := drivers.NewMockStore(ctrl)

	runnerStore := NewMockStore(ctrl)
	runnerStore.EXPECT().ListEnabledPipelines(gomock.Any()).Return(nil, nil)

	var (
		runner           *Runner
		connectorFactory drivers.Factory
	)
	app := fxtest.New(t,
		fx.Supply(drivers.NewServiceConfig("", testing.Verbose())),
		fx.Supply(http.DefaultClient),
		fx.Provide(httpclient.NewStackAuthenticatedClientFromHTTPClient),
		fx.Supply(gochannel.Config{}),
		fx.Supply(fx.Annotate(logging.Testing(), fx.As(new(logging.Logger)))),
		fx.Supply(fx.Annotate(watermill.NopLogger{}, fx.As(new(watermill.LoggerAdapter)))),
		fx.Supply(fx.Annotate(runnerStore, fx.As(new(Store)))),
		fx.Supply(fx.Annotate(driversStore, fx.As(new(drivers.Store)))),
		fx.Provide(fx.Annotate(gochannel.NewGoChannel, fx.As(new(message.Subscriber)))),
		NewModule("goo", modules.PullConfiguration{
			ModuleURLTpl: template.Must(template.New("").Parse("http://localhost")),
			PullPageSize: 100,
		}),
		fx.Populate(&runner, &connectorFactory),
	)
	require.NoError(t, app.Start(ctx))
	require.Eventually(t, runner.IsReady, time.Second, 20*time.Millisecond)
	require.IsType(t, &drivers.DriverFactoryWithBatching{}, connectorFactory)
	require.NoError(t, app.Stop(ctx))
}
