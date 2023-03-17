package bus

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

const (
	RevertedTransactionTopic  = "ledger.reverted_transaction"
	SavedMetadataTopic        = "ledger.saved_metadata"
	CommittedTransactionTopic = "ledger.committed_transaction"
)

type ledgerMonitor struct {
	publisher message.Publisher
}

var _ ledger.Monitor = &ledgerMonitor{}

func newLedgerMonitor(publisher message.Publisher) *ledgerMonitor {
	m := &ledgerMonitor{
		publisher: publisher,
	}
	return m
}

func LedgerMonitorModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				newLedgerMonitor,
				fx.ParamTags(``, `group:"monitorOptions"`),
			),
		),
		ledger.ProvideResolverOption(func(monitor *ledgerMonitor) ledger.ResolveOptionFn {
			return ledger.WithMonitor(monitor)
		}),
	)
}

func (l *ledgerMonitor) CommittedTransactions(ctx context.Context, ledger string, txs ...core.ExpandedTransaction) {
	committedTransactions, err := newEventCommittedTransactions(ledger, txs...)
	if err != nil {
		// TODO(polo): add metrics
		logging.FromContext(ctx).Errorf("building committed transactions event: %s", err)
	}

	l.publish(ctx, CommittedTransactionTopic, committedTransactions)
}

func (l *ledgerMonitor) SavedMetadata(ctx context.Context, ledger, targetID string, targetType core.TargetType, metadata core.Metadata) {
	savedMetadata, err := newEventSavedMetadata(ledger, targetID, targetType, metadata)
	if err != nil {
		// TODO(polo): add metrics
		logging.FromContext(ctx).Errorf("building saved metadata event: %s", err)
	}

	l.publish(ctx, SavedMetadataTopic, savedMetadata)
}

func (l *ledgerMonitor) RevertedTransaction(ctx context.Context, ledger string, reverted, revert *core.ExpandedTransaction) {
	revertedTransactions, err := newEventRevertedTransaction(ledger, reverted, revert)
	if err != nil {
		// TODO(polo): add metrics
		logging.FromContext(ctx).Errorf("building reverted transaction event: %s", err)
	}
	l.publish(ctx, RevertedTransactionTopic, revertedTransactions)
}

func (l *ledgerMonitor) publish(ctx context.Context, topic string, ev proto.Message) {
	if err := l.publisher.Publish(topic, publish.NewMessage(ctx, ev)); err != nil {
		logging.FromContext(ctx).Errorf("publishing message: %s", err)
		return
	}
}
