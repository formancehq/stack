package send

import (
	"github.com/formancehq/orchestration/internal/schema"
	"github.com/formancehq/orchestration/internal/workflow/stages"
)

type AccountUpdate struct {
	ID       string            `json:"id" validate:"required"`
	Ledger   string            `json:"ledger" validate:"required" spec:"default:main"`
	Metadata map[string]string `json:"metadata" validate:"required"`
}

type Update struct {
	Account *AccountUpdate `json:"account,omitempty"`
}

func (s Update) GetWorkflow() any {
	return RunUpdate
}

func init() {
	schema.RegisterOneOf(&Update{})
	stages.Register("update", Update{})
}
