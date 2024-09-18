package activities

import (
	"errors"

	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/storage"
	temporalworker "github.com/formancehq/stack/libs/go-libs/temporal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	storage storage.Storage

	plugins plugins.Plugins
}

func (a Activities) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "PluginInstallConnector",
			Func: a.PluginInstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextAccounts",
			Func: a.PluginFetchNextAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextBalances",
			Func: a.PluginFetchNextBalances,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextExternalAccounts",
			Func: a.PluginFetchNextExternalAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextPayments",
			Func: a.PluginFetchNextPayments,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextOthers",
			Func: a.PluginFetchNextOthers,
		}).
		Append(temporalworker.Definition{
			Name: "PluginCreateBankAccount",
			Func: a.PluginCreateBankAccount,
		}).
		Append(temporalworker.Definition{
			Name: "PluginCreateWebhooks",
			Func: a.PluginCreateWebhooks,
		}).
		Append(temporalworker.Definition{
			Name: "PluginTranslateWebhook",
			Func: a.PluginTranslateWebhook,
		}).
		Append(temporalworker.Definition{
			Name: "StorageAccountsStore",
			Func: a.StorageAccountsStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageAccountsDelete",
			Func: a.StorageAccountsDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StoragePaymentsStore",
			Func: a.StoragePaymentsStore,
		}).
		Append(temporalworker.Definition{
			Name: "StoragePaymentsDelete",
			Func: a.StoragePaymentsDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStatesGet",
			Func: a.StorageStatesGet,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStatesStore",
			Func: a.StorageStatesStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStatesDelete",
			Func: a.StorageStatesDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageTasksTreeStore",
			Func: a.StorageTasksTreeStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageTasksTreeDelete",
			Func: a.StorageTasksTreeDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageConnectorsStore",
			Func: a.StorageConnectorsStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageConnectorsDelete",
			Func: a.StorageConnectorsDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageSchedulesStore",
			Func: a.StorageSchedulesStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageSchedulesList",
			Func: a.StorageSchedulesList,
		}).
		Append(temporalworker.Definition{
			Name: "StorageSchedulesDelete",
			Func: a.StorageSchedulesDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageInstancesStore",
			Func: a.StorageInstancesStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageInstancesUpdate",
			Func: a.StorageInstancesUpdate,
		}).
		Append(temporalworker.Definition{
			Name: "StorageInstancesDelete",
			Func: a.StorageInstancesDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageBankAccountsDeleteRelatedAccounts",
			Func: a.StorageBankAccountsDeleteRelatedAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "StorageBankAccountsAddRelatedAccount",
			Func: a.StorageBankAccountsAddRelatedAccount,
		}).
		Append(temporalworker.Definition{
			Name: "StorageBankAccountsGet",
			Func: a.StorageBankAccountsGet,
		}).
		Append(temporalworker.Definition{
			Name: "StorageBalancesDelete",
			Func: a.StorageBalancesDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageBalancesStore",
			Func: a.StorageBalancesStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageWebhooksConfigsStore",
			Func: a.StorageWebhooksConfigsStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageWebhooksConfigsDelete",
			Func: a.StorageWebhooksConfigsDelete,
		}).
		Append(temporalworker.Definition{
			Name: "StorageWebhooksStore",
			Func: a.StorageWebhooksStore,
		}).
		Append(temporalworker.Definition{
			Name: "StorageWebhooksDelete",
			Func: a.StorageWebhooksDelete,
		})
}

func New(storage storage.Storage, plugins plugins.Plugins) Activities {
	return Activities{
		storage: storage,
		plugins: plugins,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, args ...any) error {
	if err := workflow.ExecuteActivity(ctx, activity, args...).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		return err
	}
	return nil
}