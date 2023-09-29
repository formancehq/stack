package all

import (
	_ "github.com/formancehq/operator/internal/modules/auth"
	_ "github.com/formancehq/operator/internal/modules/control"
	_ "github.com/formancehq/operator/internal/modules/gateway"
	_ "github.com/formancehq/operator/internal/modules/ledger"
	_ "github.com/formancehq/operator/internal/modules/orchestration"
	_ "github.com/formancehq/operator/internal/modules/payments"
	_ "github.com/formancehq/operator/internal/modules/search"
	_ "github.com/formancehq/operator/internal/modules/stargate"
	_ "github.com/formancehq/operator/internal/modules/wallets"
	_ "github.com/formancehq/operator/internal/modules/webhooks"
)
