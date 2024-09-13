package bankingcircle

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) createBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	// We can't create bank accounts in Banking Circle since they do not store
	// the bank account information. We just have to return the related formance
	// account in order to use it in the future.
	raw, err := json.Marshal(req.BankAccount)
	if err != nil {
		return models.CreateBankAccountResponse{}, err
	}

	return models.CreateBankAccountResponse{
		RelatedAccount: models.PSPAccount{
			Reference: req.BankAccount.ID.String(),
			CreatedAt: req.BankAccount.CreatedAt,
			Name:      &req.BankAccount.Name,
			Raw:       raw,
		},
	}, nil
}
