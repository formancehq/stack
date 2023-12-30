package ledgers

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
)

//go:embed Caddyfile.ledger.gotpl
var Caddyfile string

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Ledger, error) {
	ret, err := stacks.GetSingleStackDependencyObject(ctx, stackName, &v1beta1.LedgerList{})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}

	return ret.(*v1beta1.Ledger), nil
}
