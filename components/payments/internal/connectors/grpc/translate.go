package grpc

import (
	"errors"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/grpc/proto"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TranslateAccount(account models.PSPAccount) *proto.Account {
	return &proto.Account{
		Reference: account.Reference,
		Name: func() *wrapperspb.StringValue {
			if account.Name == nil {
				return nil
			}

			return wrapperspb.String(*account.Name)
		}(),
		CreatedAt: timestamppb.New(account.CreatedAt),
		SyncedAt:  timestamppb.New(time.Now().UTC()),
		DefaultAsset: func() *wrapperspb.StringValue {
			if account.DefaultAsset == nil {
				return nil
			}

			return wrapperspb.String(*account.DefaultAsset)
		}(),
		Metadata: account.Metadata,
		Raw:      account.Raw,
	}
}

func TranslateProtoAccount(account *proto.Account) models.PSPAccount {
	return models.PSPAccount{
		Reference: account.Reference,
		CreatedAt: account.CreatedAt.AsTime(),
		Name: func() *string {
			if account.Name == nil {
				return nil
			}

			return pointer.For(account.Name.GetValue())
		}(),
		DefaultAsset: func() *string {
			if account.DefaultAsset == nil {
				return nil
			}

			return pointer.For(account.DefaultAsset.GetValue())
		}(),
		Metadata: account.Metadata,
		Raw:      account.Raw,
	}
}

func TranslateBankAccount(bankAccount models.BankAccount) *proto.BankAccount {
	return &proto.BankAccount{
		Id:        bankAccount.ID.String(),
		CreatedAt: timestamppb.New(bankAccount.CreatedAt),
		Name:      bankAccount.Name,
		AccountNumber: func() *wrapperspb.StringValue {
			if bankAccount.AccountNumber == nil {
				return nil
			}
			return wrapperspb.String(*bankAccount.AccountNumber)
		}(),
		Iban: func() *wrapperspb.StringValue {
			if bankAccount.IBAN == nil {
				return nil
			}
			return wrapperspb.String(*bankAccount.IBAN)
		}(),
		SwiftBicCode: func() *wrapperspb.StringValue {
			if bankAccount.SwiftBicCode == nil {
				return nil
			}
			return wrapperspb.String(*bankAccount.SwiftBicCode)
		}(),
		Country: func() *wrapperspb.StringValue {
			if bankAccount.Country == nil {
				return nil
			}
			return wrapperspb.String(*bankAccount.Country)
		}(),
		Metadata: bankAccount.Metadata,
	}
}

func TranslateProtoBankAccount(bankAccount *proto.BankAccount) models.BankAccount {
	uuid, err := uuid.Parse(bankAccount.Id)
	if err != nil {
		panic(err)
	}

	return models.BankAccount{
		ID:        uuid,
		CreatedAt: bankAccount.CreatedAt.AsTime(),
		Name:      bankAccount.Name,
		AccountNumber: func() *string {
			if bankAccount.AccountNumber == nil {
				return nil
			}
			return pointer.For(bankAccount.AccountNumber.GetValue())
		}(),
		IBAN: func() *string {
			if bankAccount.Iban == nil {
				return nil
			}
			return pointer.For(bankAccount.Iban.GetValue())
		}(),
		SwiftBicCode: func() *string {
			if bankAccount.SwiftBicCode == nil {
				return nil
			}
			return pointer.For(bankAccount.SwiftBicCode.GetValue())
		}(),
		Country: func() *string {
			if bankAccount.Country == nil {
				return nil
			}
			return pointer.For(bankAccount.Country.GetValue())
		}(),
		Metadata: bankAccount.Metadata,
	}
}

func TranslateBalance(balance models.PSPBalance) *proto.Balance {
	return &proto.Balance{
		AccountReference: balance.AccountReference,
		CreatedAt:        timestamppb.New(balance.CreatedAt),
		Balance: &proto.Monetary{
			Asset:  balance.Asset,
			Amount: []byte(balance.Amount.Text(10)),
		},
	}
}

func TranslateProtoBalance(balance *proto.Balance) (models.PSPBalance, error) {
	amount, ok := big.NewInt(0).SetString(string(balance.Balance.Amount), 10)
	if !ok {
		return models.PSPBalance{}, errors.New("failed to parse amount")
	}

	return models.PSPBalance{
		AccountReference: balance.AccountReference,
		CreatedAt:        balance.CreatedAt.AsTime(),
		Amount:           amount,
		Asset:            balance.Balance.Asset,
	}, nil
}

func TranslatePayment(payment models.PSPPayment) *proto.Payment {
	return &proto.Payment{
		Reference:   payment.Reference,
		CreatedAt:   timestamppb.New(payment.CreatedAt),
		SyncedAt:    timestamppb.New(time.Now().UTC()),
		PaymentType: proto.PaymentType(payment.Type),
		Amount: &proto.Monetary{
			Asset:  payment.Asset,
			Amount: []byte(payment.Amount.Text(10)),
		},
		Scheme: proto.PaymentScheme(payment.Scheme),
		Status: proto.PaymentStatus(payment.Status),
		SourceAccountReference: func() *wrapperspb.StringValue {
			if payment.SourceAccountReference == nil {
				return nil
			}

			return wrapperspb.String(*payment.SourceAccountReference)
		}(),
		DestinationAccountReference: func() *wrapperspb.StringValue {
			if payment.DestinationAccountReference == nil {
				return nil
			}

			return wrapperspb.String(*payment.DestinationAccountReference)
		}(),
		Metadata: payment.Metadata,
		Raw:      payment.Raw,
	}
}

func TranslateProtoPayment(payment *proto.Payment) (models.PSPPayment, error) {
	amount, ok := big.NewInt(0).SetString(string(payment.Amount.Amount), 10)
	if !ok {
		return models.PSPPayment{}, errors.New("failed to parse amount")
	}
	return models.PSPPayment{
		Reference: payment.Reference,
		CreatedAt: payment.CreatedAt.AsTime(),
		Type:      models.PaymentType(payment.PaymentType),
		Amount:    amount,
		Asset:     payment.Amount.Asset,
		Scheme:    models.PaymentScheme(payment.Scheme),
		Status:    models.PaymentStatus(payment.Status),
		SourceAccountReference: func() *string {
			if payment.SourceAccountReference == nil {
				return nil
			}

			return pointer.For(payment.SourceAccountReference.GetValue())
		}(),
		DestinationAccountReference: func() *string {
			if payment.DestinationAccountReference == nil {
				return nil
			}

			return pointer.For(payment.DestinationAccountReference.GetValue())
		}(),
		Metadata: payment.Metadata,
		Raw:      payment.Raw,
	}, nil
}

func TranslateTask(taskTree models.TaskTree) *proto.TaskTree {
	res := proto.TaskTree{
		NextTasks:    []*proto.TaskTree{},
		Name:         taskTree.Name,
		Periodically: taskTree.Periodically,
		Task:         nil,
	}

	switch taskTree.TaskType {
	case models.TASK_FETCH_ACCOUNTS:
		res.Task = &proto.TaskTree_FetchAccounts_{
			FetchAccounts: &proto.TaskTree_FetchAccounts{},
		}
	case models.TASK_FETCH_EXTERNAL_ACCOUNTS:
		res.Task = &proto.TaskTree_FetchExternalAccounts_{
			FetchExternalAccounts: &proto.TaskTree_FetchExternalAccounts{},
		}
	case models.TASK_FETCH_PAYMENTS:
		res.Task = &proto.TaskTree_FetchPayments_{
			FetchPayments: &proto.TaskTree_FetchPayments{},
		}
	case models.TASK_FETCH_OTHERS:
		res.Task = &proto.TaskTree_FetchOthers_{
			FetchOthers: &proto.TaskTree_FetchOthers{},
		}
	case models.TASK_FETCH_BALANCES:
		res.Task = &proto.TaskTree_FetchBalances_{
			FetchBalances: &proto.TaskTree_FetchBalances{},
		}
	case models.TASK_CREATE_WEBHOOKS:
		res.Task = &proto.TaskTree_CreateWebhooks_{
			CreateWebhooks: &proto.TaskTree_CreateWebhooks{},
		}
	default:
		// TODO(polo): better error handling
		panic("unknown task type")
	}

	for _, nextTask := range taskTree.NextTasks {
		res.NextTasks = append(res.NextTasks, TranslateTask(nextTask))
	}

	return &res
}

func TranslateProtoTask(task *proto.TaskTree) models.TaskTree {
	res := models.TaskTree{
		TaskType:     0,
		Name:         task.Name,
		Periodically: task.Periodically,
		NextTasks:    []models.TaskTree{},
	}

	switch task.Task.(type) {
	case *proto.TaskTree_FetchAccounts_:
		res.TaskType = models.TASK_FETCH_ACCOUNTS
		res.TaskTreeFetchAccounts = &models.TaskTreeFetchAccounts{}
	case *proto.TaskTree_FetchExternalAccounts_:
		res.TaskType = models.TASK_FETCH_EXTERNAL_ACCOUNTS
		res.TaskTreeFetchExternalAccounts = &models.TaskTreeFetchExternalAccounts{}
	case *proto.TaskTree_FetchPayments_:
		res.TaskType = models.TASK_FETCH_PAYMENTS
		res.TaskTreeFetchPayments = &models.TaskTreeFetchPayments{}
	case *proto.TaskTree_FetchOthers_:
		res.TaskType = models.TASK_FETCH_OTHERS
		res.TaskTreeFetchOther = &models.TaskTreeFetchOther{}
	case *proto.TaskTree_FetchBalances_:
		res.TaskType = models.TASK_FETCH_BALANCES
		res.TaskTreeFetchBalances = &models.TaskTreeFetchBalances{}
	case *proto.TaskTree_CreateWebhooks_:
		res.TaskType = models.TASK_CREATE_WEBHOOKS
		res.TaskTreeCreateWebhooks = &models.TaskTreeCreateWebhooks{}
	default:
		panic("unknown task type")
	}

	for _, nextTask := range task.NextTasks {
		res.NextTasks = append(res.NextTasks, TranslateProtoTask(nextTask))
	}

	return res
}

func TranslateWorkflow(workflows models.Tasks) *proto.Workflow {
	res := proto.Workflow{}

	for _, task := range workflows {
		res.Tasks = append(res.Tasks, TranslateTask(task))
	}

	return &res
}

func TranslateProtoWorkflow(workflow *proto.Workflow) models.Tasks {
	res := models.Tasks{}

	for _, task := range workflow.Tasks {
		res = append(res, TranslateProtoTask(task))
	}

	return res
}

func TranslateWebhook(from models.PSPWebhook) *proto.Webhook {
	headers := make(map[string]*proto.Webhook_Values)
	for k, v := range from.Headers {
		headers[k] = &proto.Webhook_Values{Values: v}
	}

	queryValues := make(map[string]*proto.Webhook_Values)
	for k, v := range from.QueryValues {
		queryValues[k] = &proto.Webhook_Values{Values: v}
	}

	return &proto.Webhook{
		Headers:     headers,
		QueryValues: queryValues,
		Body:        from.Body,
	}
}

func TranslateProtoWebhook(from *proto.Webhook) models.PSPWebhook {
	headers := make(map[string][]string)
	for k, v := range from.Headers {
		headers[k] = v.Values
	}

	queryValues := make(map[string][]string)
	for k, v := range from.QueryValues {
		queryValues[k] = v.Values
	}

	return models.PSPWebhook{
		QueryValues: queryValues,
		Headers:     headers,
		Body:        from.Body,
	}
}
