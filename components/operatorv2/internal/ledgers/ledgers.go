package ledgers

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/utils"
)

//go:embed Caddyfile.ledger.gotpl
var Caddyfile string

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Ledger, error) {
	return utils.GetSingleStackDependencyObject[*v1beta1.LedgerList, *v1beta1.Ledger](ctx, stackName)
}
