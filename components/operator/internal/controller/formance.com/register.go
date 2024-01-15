package formance_com

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/reconcilers"
)

func init() {
	reconcilers.Register(
		reconcilers.New[*v1beta1.Stack](ForStack()),
		reconcilers.ForModule[*v1beta1.Ledger](ForLedger()),
		reconcilers.ForModule[*v1beta1.Auth](ForAuth()),
		reconcilers.ForModule[*v1beta1.Wallets](ForWallets()),
		reconcilers.ForModule[*v1beta1.Orchestration](ForOrchestration()),
		reconcilers.ForModule[*v1beta1.Payments](ForPayments()),
		reconcilers.ForModule[*v1beta1.Reconciliation](ForReconciliation()),
		reconcilers.ForModule[*v1beta1.Webhooks](ForWebhooks()),
		reconcilers.ForModule[*v1beta1.Search](ForSearch()),
		reconcilers.ForStackDependency[*v1beta1.BrokerTopic](ForBrokerTopic()),
		reconcilers.ForStackDependency[*v1beta1.BrokerTopicConsumer](ForBrokerTopicConsumer()),
		reconcilers.ForStackDependency[*v1beta1.HTTPAPI](ForHTTPAPI()),
		reconcilers.ForStackDependency[*v1beta1.Gateway](ForGateway()),
		reconcilers.ForStackDependency[*v1beta1.Database](ForDatabase()),
		reconcilers.ForStackDependency[*v1beta1.AuthClient](ForAuthClient()),
		reconcilers.ForStackDependency[*v1beta1.StreamProcessor](ForStreamProcessor()),
		reconcilers.ForStackDependency[*v1beta1.Stream](ForStream()),
		reconcilers.ForStackDependency[*v1beta1.Stargate](ForStargate()),
	)
}
