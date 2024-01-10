package stack_formance_com

import (
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	"github.com/formancehq/operator/internal/reconcilers"
)

func init() {
	reconcilers.Register(
		reconcilers.New[*v1beta3.Stack](ForStack()),
	)
}
