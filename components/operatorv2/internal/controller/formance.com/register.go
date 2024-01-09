package formance_com

import (
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
)

func init() {
	reconcilers.Register(
		reconcilers.New[*v1beta1.Stack](ForStack()),
		reconcilers.New[*v1beta1.Topic](ForTopic()),
		reconcilers.New[*v1beta1.TopicQuery](ForTopicQuery()),
		reconcilers.New[*v1beta1.Ledger](ForLedger()),
		reconcilers.New[*v1beta1.HTTPAPI](ForHTTPAPI()),
		reconcilers.New[*v1beta1.Gateway](ForGateway()),
		reconcilers.New[*v1beta1.Auth](ForAuth()),
		reconcilers.New[*v1beta1.Database](ForDatabase()),
		reconcilers.New[*v1beta1.AuthClient](ForAuthClient()),
		reconcilers.New[*v1beta1.Wallets](ForWallets()),
		reconcilers.New[*v1beta1.Orchestration](ForOrchestration()),
		reconcilers.New[*v1beta1.Payments](ForPayments()),
		reconcilers.New[*v1beta1.Reconciliation](ForReconciliation()),
		reconcilers.New[*v1beta1.Webhooks](ForWebhooks()),
		reconcilers.New[*v1beta1.Search](ForSearch()),
		reconcilers.New[*v1beta1.StreamProcessor](ForStreamProcessor()),
		reconcilers.New[*v1beta1.Stream](ForStream()),
		reconcilers.New[*v1beta1.Stargate](ForStargate()),
	)
}
