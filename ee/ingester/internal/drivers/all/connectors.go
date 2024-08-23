package all

import (
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/clickhouse"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/elasticsearch"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/http"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/noop"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/stdout"
)

func Register(connectorRegistry *drivers.Registry) {
	connectorRegistry.RegisterConnector("elasticsearch", elasticsearch.NewConnector)
	connectorRegistry.RegisterConnector("clickhouse", clickhouse.NewConnector)
	connectorRegistry.RegisterConnector("stdout", stdout.NewConnector)
	connectorRegistry.RegisterConnector("http", http.NewConnector)
	connectorRegistry.RegisterConnector("noop", noop.NewConnector)
}
