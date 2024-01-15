package formance_com

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/reconcilers"
)

func init() {
	reconcilers.Register(
		reconcilers.New[*v1beta1.Stack](ForStack()),
		reconcilers.NewStackDependency[*v1beta1.BrokerTopic](ForBrokerTopic()),
		reconcilers.NewStackDependency[*v1beta1.BrokerTopicConsumer](ForBrokerTopicConsumer()),
		reconcilers.NewStackDependency[*v1beta1.Ledger](ForLedger()),
		reconcilers.NewStackDependency[*v1beta1.HTTPAPI](ForHTTPAPI()),
		reconcilers.NewStackDependency[*v1beta1.Gateway](ForGateway()),
		reconcilers.NewStackDependency[*v1beta1.Auth](ForAuth()),
		reconcilers.NewStackDependency[*v1beta1.Database](ForDatabase()),
		reconcilers.NewStackDependency[*v1beta1.AuthClient](ForAuthClient()),
		reconcilers.NewStackDependency[*v1beta1.Wallets](ForWallets()),
		reconcilers.NewStackDependency[*v1beta1.Orchestration](ForOrchestration()),
		reconcilers.NewStackDependency[*v1beta1.Payments](ForPayments()),
		reconcilers.NewStackDependency[*v1beta1.Reconciliation](ForReconciliation()),
		reconcilers.NewStackDependency[*v1beta1.Webhooks](ForWebhooks()),
		reconcilers.NewStackDependency[*v1beta1.Search](ForSearch()),
		reconcilers.NewStackDependency[*v1beta1.StreamProcessor](ForStreamProcessor()),
		reconcilers.NewStackDependency[*v1beta1.Stream](ForStream()),
		reconcilers.NewStackDependency[*v1beta1.Stargate](ForStargate()),
	)
}
