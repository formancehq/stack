package workflow

import (
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type CreateBankAccount struct {
	ConnectorID   models.ConnectorID
	BankAccountID uuid.UUID
}

func (w Workflow) runCreateBankAccount(
	ctx workflow.Context,
	createBankAccount CreateBankAccount,
) (*models.BankAccount, error) {
	bankAccount, err := activities.StorageBankAccountsGet(
		infiniteRetryContext(ctx),
		createBankAccount.BankAccountID,
		true,
	)
	if err != nil {
		return nil, err
	}

	createBAResponse, err := activities.PluginCreateBankAccount(
		infiniteRetryContext(ctx),
		createBankAccount.ConnectorID,
		models.CreateBankAccountRequest{
			BankAccount: *bankAccount,
		},
	)
	if err != nil {
		return nil, err
	}

	account := models.FromPSPAccount(
		createBAResponse.RelatedAccount,
		models.ACCOUNT_TYPE_EXTERNAL,
		createBankAccount.ConnectorID,
	)

	err = activities.StorageAccountsStore(
		infiniteRetryContext(ctx),
		[]models.Account{account},
	)
	if err != nil {
		return nil, err
	}

	relatedAccount := models.BankAccountRelatedAccount{
		BankAccountID: createBankAccount.BankAccountID,
		AccountID:     account.ID,
		ConnectorID:   createBankAccount.ConnectorID,
		CreatedAt:     createBAResponse.RelatedAccount.CreatedAt,
	}

	err = activities.StorageBankAccountsAddRelatedAccount(
		infiniteRetryContext(ctx),
		relatedAccount,
	)
	if err != nil {
		return nil, err
	}

	bankAccount.RelatedAccounts = append(bankAccount.RelatedAccounts, relatedAccount)

	return bankAccount, nil
}

var RunCreateBankAccount any

func init() {
	RunCreateBankAccount = Workflow{}.runCreateBankAccount
}
