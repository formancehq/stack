package service

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type ReconciliationRequest struct {
	ReconciledAtLedger   time.Time `json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time `json:"reconciledAtPayment"`
}

func (r *ReconciliationRequest) Validate() error {
	if r.ReconciledAtLedger.IsZero() {
		return errors.New("missing reconciledAtLedger")
	}

	if r.ReconciledAtLedger.After(time.Now()) {
		return errors.New("reconciledAtLedger must be in the past")
	}

	if r.ReconciledAtPayments.IsZero() {
		return errors.New("missing reconciledAtPayments")
	}

	if r.ReconciledAtPayments.After(time.Now()) {
		return errors.New("ReconciledAtPayments must be in the past")
	}

	return nil
}

func (s *Service) Reconciliation(ctx context.Context, policyID string, req *ReconciliationRequest) (*models.Reconciliation, error) {
	id, err := uuid.Parse(policyID)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	eg, ctxGroup := errgroup.WithContext(ctx)
	policy, err := s.store.GetPolicy(ctx, id)
	if err != nil {
		return nil, newStorageError(err, "failed to get policy")
	}

	var paymentBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		paymentBalance, err = s.getPaymentPoolBalance(ctxGroup, policy.PaymentsPoolID.String(), req.ReconciledAtPayments)
		return err
	})

	var ledgerBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		ledgerBalance, err = s.getAccountsAggregatedBalance(ctxGroup, policy.LedgerName, policy.LedgerQuery, req.ReconciledAtLedger)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	res := &models.Reconciliation{
		ID:                   uuid.New(),
		PolicyID:             policy.ID,
		CreatedAt:            time.Now().UTC(),
		ReconciledAtLedger:   req.ReconciledAtLedger,
		ReconciledAtPayments: req.ReconciledAtPayments,
		Status:               models.ReconciliationOK,
		PaymentsBalances:     paymentBalance,
		LedgerBalances:       ledgerBalance,
	}

	var reconciliationError bool
	if len(paymentBalance) != len(ledgerBalance) {
		res.Status = models.ReconciliationNotOK
		res.Error = "different number of assets"
		reconciliationError = true
		return res, nil
	}

	if !reconciliationError {
		for paymentAsset, paymentBalance := range paymentBalance {
			ledgerBalance, ok := ledgerBalance[paymentAsset]
			if !ok {
				res.Status = models.ReconciliationNotOK
				res.Error = fmt.Sprintf("missing asset %s in ledger", paymentAsset)
				break
			}

			if paymentBalance.Cmp(ledgerBalance) != 0 {
				res.Status = models.ReconciliationNotOK
				res.Error = fmt.Sprintf("different balance for asset %s", paymentAsset)
				break
			}
		}
	}

	if err := s.store.CreateReconciation(ctx, res); err != nil {
		return nil, newStorageError(err, "failed to create reconciliation")
	}

	return res, nil
}

func (s *Service) getAccountsAggregatedBalance(ctx context.Context, ledgerName string, ledgerAggregatedBalanceQuery map[string]interface{}, at time.Time) (map[string]*big.Int, error) {
	balances, err := s.client.Ledger.V2GetBalancesAggregated(
		ctx,
		operations.V2GetBalancesAggregatedRequest{
			RequestBody: ledgerAggregatedBalanceQuery,
			Ledger:      ledgerName,
			Pit:         &at,
		},
	)
	if err != nil {
		return nil, err
	}

	if balances.StatusCode != 200 {
		return nil, errors.New("failed to get aggregated balances")
	}

	if balances.V2AggregateBalancesResponse == nil {
		return nil, errors.New("no aggregated balance")
	}

	balanceMap := make(map[string]*big.Int)
	for asset, balance := range balances.V2AggregateBalancesResponse.Data {
		balanceMap[asset] = balance
	}

	return balanceMap, nil
}

func (s *Service) getPaymentPoolBalance(ctx context.Context, paymentPoolID string, at time.Time) (map[string]*big.Int, error) {
	balances, err := s.client.Payments.GetPoolBalances(
		ctx,
		operations.GetPoolBalancesRequest{
			At:     at,
			PoolID: paymentPoolID,
		},
	)
	if err != nil {
		return nil, err
	}

	if balances.StatusCode != 200 {
		fmt.Println(err)
		body, _ := io.ReadAll(balances.RawResponse.Body)
		fmt.Println(string(body))
		return nil, errors.New("failed to get pool balances")
	}

	if balances.PoolBalancesResponse == nil {
		return nil, errors.New("no pool balance")
	}

	balanceMap := make(map[string]*big.Int)
	for _, balance := range balances.PoolBalancesResponse.Data.Balances {
		balanceMap[balance.GetAsset()] = balance.GetAmount()
	}

	return balanceMap, nil
}

func (s *Service) GetReconciliation(ctx context.Context, id string) (*models.Reconciliation, error) {
	rID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	reco, err := s.store.GetReconciliation(ctx, rID)
	return reco, newStorageError(err, "getting reconciliation")
}

func (s *Service) ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*api.Cursor[models.Reconciliation], error) {
	reconciliations, err := s.store.ListReconciliations(ctx, q)
	return reconciliations, newStorageError(err, "listing reconciliations")
}
