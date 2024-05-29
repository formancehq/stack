package worker

import (
	"context"
	"fmt"
	"math/big"
	"runtime/debug"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/reconciliation/internal/client"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type handlerEvent func() error
type handlerEvents []handlerEvent

type Listener struct {
	logger logging.Logger

	store  Store
	client client.SDKSearchFormance

	ruleTreeCache map[uuid.UUID]*RuleTree
}

func NewListener(
	logger logging.Logger,
	store Store,
	client client.SDKSearchFormance,
) *Listener {
	return &Listener{
		logger:        logger,
		store:         store,
		client:        client,
		ruleTreeCache: make(map[uuid.UUID]*RuleTree),
	}
}

func (l *Listener) handleMessage(msg *message.Message) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
		}
	}()

	var event *publish.EventMessage
	span, event, err := publish.UnmarshalMessage(msg)
	if err != nil {
		logging.FromContext(msg.Context()).Error(err.Error())
		return err
	}

	switch event.Type {
	case EventTypeCommittedTransactions, EventTypeRevertedTransaction, EventTypeSavedPayments:
	default:
		// We don't care about other events
		return nil
	}

	ctx := msg.Context()

	policies, err := l.getPolicies(ctx)
	if err != nil {
		span.RecordError(err)
		span.End()
		return err
	}

	eg, ctxGroup := errgroup.WithContext(ctx)
	for _, policy := range policies {
		policy := policy
		eg.Go(func() error {
			return l.handlerEventByPolicy(ctxGroup, policy, event)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (l *Listener) handlerEventByPolicy(
	ctx context.Context,
	policy models.Policy,
	event *publish.EventMessage,
) error {
	rules, err := l.getRulesTree(ctx, policy.Rule)
	if err != nil {
		return err
	}

	fns := make(handlerEvents, 0)
	switch event.Type {
	case EventTypeCommittedTransactions:
		ev, ok := event.Payload.(committedTransactions)
		if !ok {
			return fmt.Errorf("invalid event payload type")
		}

		for _, tx := range ev.Transactions {
			inserted, err := l.store.CreateTransactionBasedReconciliation(ctx, &models.ReconciliationTransactionBased{
				PaymentID:     "",
				TransactionID: tx.ID,
				RuleID:        0,
				PolicyID:      policy.ID,
				CreatedAt:     event.Date,
				Status:        models.ReconciliationStatusPending,
			})
			if err != nil {
				return err
			}

			if !inserted {
				// Already inserted, which means it was processed before. if the
				// status is pending, the worker that handled previous pending
				// row will handle this one too.
				continue
			}

			fns = append(fns, func() error {
				return l.handleCommittedTransactions(ctx, policy, rules, tx)
			})
		}

	case EventTypeRevertedTransaction:
		ev, ok := event.Payload.(revertedTransaction)
		if !ok {
			return fmt.Errorf("invalid event payload type")
		}

		inserted, err := l.store.CreateTransactionBasedReconciliation(ctx, &models.ReconciliationTransactionBased{
			PaymentID:     "",
			TransactionID: ev.RevertTransaction.ID,
			RuleID:        0,
			PolicyID:      policy.ID,
			CreatedAt:     event.Date,
			Status:        models.ReconciliationStatusPending,
		})
		if err != nil {
			return err
		}

		if inserted {
			// Was not processed before, so we need to handle it
			fns = append(fns, func() error {
				return handleRevertedTransactions(ctx, ev.RevertTransaction, ev.RevertedTransaction)
			})
		}

	case EventTypeSavedPayments:
		// TODO(polo): should we handle pending payments ?
		ev, ok := event.Payload.(paymentMessagePayload)
		if !ok {
			return fmt.Errorf("invalid event payload type")
		}

		if ev.Status != "PENDING" {
			inserted, err := l.store.CreateTransactionBasedReconciliation(ctx, &models.ReconciliationTransactionBased{
				PaymentID:     ev.ID,
				TransactionID: big.NewInt(-1),
				RuleID:        0,
				PolicyID:      policy.ID,
				CreatedAt:     event.Date,
				Status:        models.ReconciliationStatusPending,
			})
			if err != nil {
				return err
			}

			if inserted {
				// Was not processed before, so we need to handle it
				fns = append(fns, func() error {
					return handleSavedPayments(ctx, ev)
				})
			}
		}

	default:
		return nil
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (l *Listener) registerListener(logger logging.Logger, r *message.Router, s message.Subscriber, topics []string) {
	for _, topic := range topics {
		r.AddNoPublisherHandler(fmt.Sprintf("reco-listen-%s-events", topic), topic, s, func(msg *message.Message) error {
			if err := l.handleMessage(msg); err != nil {
				logger.Errorf("error handling message: %s", err)
				return err
			}
			return nil
		})
	}
}

func (l *Listener) getPolicies(ctx context.Context) ([]models.Policy, error) {
	queryBuilder := query.And(
		query.Match("enabled", true),
		query.Match("type", string(models.PolicyTypeTransactionBased)),
	)

	res, err := l.store.ListPolicies(ctx, storage.NewListPoliciesQuery(
		bunpaginate.NewPaginatedQueryOptions[storage.PoliciesFilters](storage.PoliciesFilters{}).WithQueryBuilder(queryBuilder),
	))
	if err != nil {
		return nil, err
	}

	policies := make([]models.Policy, 0)
	policies = append(policies, res.Data...)
	for res.HasMore {
		var q storage.ListPoliciesQuery
		err = bunpaginate.UnmarshalCursor(res.Next, &q)
		if err != nil {
			return nil, err
		}

		res, err = l.store.ListPolicies(ctx, q)
		if err != nil {
			return nil, err
		}

		policies = append(policies, res.Data...)
	}

	return policies, nil
}
