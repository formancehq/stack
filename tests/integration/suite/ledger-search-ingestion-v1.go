package suite

import (
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
)

var _ = WithModules([]*Module{modules.Search, modules.Ledger}, func() {
	When("sending a v1 created transaction event in the event bus", func() {

	})
})
