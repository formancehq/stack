package atlar

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

const (
	taskNameFetchAccounts = "fetch_accounts"
	taskNameFetchPayments = "fetch_payments"
	// taskNameInitiatePayment   = "initiate_payment"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name    string `json:"name" yaml:"name" bson:"name"`
	Key     string `json:"key" yaml:"key" bson:"key"`
	Main    bool   `json:"main,omitempty" yaml:"main" bson:"main"`
	Account string `json:"account,omitempty" yaml:"account" bson:"account"`
}

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	transport := httptransport.New(
		config.BaseUrl.Host,
		config.BaseUrl.Path,
		[]string{config.BaseUrl.Scheme},
	)
	basicAuth := httptransport.BasicAuth(config.AccessKey, config.Secret)
	transport.DefaultAuthentication = basicAuth
	client := atlar_client.New(transport, strfmt.Default)

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Main {
			return MainTask(logger)
		}

		switch taskDescriptor.Key {
		case taskNameFetchAccounts:
			return FetchAccountsTask(config, client)
		case taskNameFetchPayments:
			return FetchPaymentsTask(config, taskDescriptor.Account, client)
		default:
			return nil
		}
	}
}
