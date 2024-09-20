package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"

	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type ReconciliationRequest struct {
	ReconciledAtLedger   time.Time `json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time `json:"reconciledAtPayments"`
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

	var paymentsBalances map[string]*big.Int
	eg.Go(func() error {
		var err error
		paymentsBalances, err = s.getPaymentPoolBalance(ctxGroup, policy.PaymentsPoolID.String(), req.ReconciledAtPayments)
		return err
	})

	var ledgerBalances map[string]*big.Int
	eg.Go(func() error {
		var err error
		ledgerBalances, err = s.getAccountsAggregatedBalance(ctxGroup, policy.LedgerName, policy.LedgerQuery, req.ReconciledAtLedger)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ledgerBalances, paymentsBalances = harmonizeBalances(ledgerBalances, paymentsBalances)

	res := &models.Reconciliation{
		ID:                   uuid.New(),
		PolicyID:             policy.ID,
		CreatedAt:            time.Now().UTC(),
		ReconciledAtLedger:   req.ReconciledAtLedger,
		ReconciledAtPayments: req.ReconciledAtPayments,
		Status:               models.ReconciliationOK,
		PaymentsBalances:     paymentsBalances,
		LedgerBalances:       ledgerBalances,
		DriftBalances:        make(map[string]*big.Int),
	}

	var reconciliationError bool
	if len(paymentsBalances) != len(ledgerBalances) {
		res.Status = models.ReconciliationNotOK
		res.Error = "different number of assets"
		reconciliationError = true
		return res, nil
	}

	if !reconciliationError {
		for asset, ledgerBalance := range ledgerBalances {
			err := s.computeDrift(res, asset, ledgerBalance, paymentsBalances[asset])
			if err != nil {
				res.Status = models.ReconciliationNotOK
				if res.Error == "" {
					res.Error = err.Error()
				} else {
					res.Error = res.Error + "; " + err.Error()
				}
			}
		}

		for asset, paymentBalance := range paymentsBalances {
			if _, ok := res.DriftBalances[asset]; ok {
				// Already computed
				continue
			}

			err := s.computeDrift(res, asset, ledgerBalances[asset], paymentBalance)
			if err != nil {
				res.Status = models.ReconciliationNotOK
				res.Error = res.Error + "; " + err.Error()
			}
		}
	}

	if err := s.store.CreateReconciation(ctx, res); err != nil {
		return nil, newStorageError(err, "failed to create reconciliation")
	}

	return res, nil
}

func (s *Service) computeDrift(
	res *models.Reconciliation,
	asset string,
	ledgerBalance *big.Int,
	paymentBalance *big.Int,
) error {
	switch {
	case ledgerBalance == nil && paymentBalance == nil:
		// Not possible
		return nil
	case ledgerBalance == nil && paymentBalance != nil:
		var balance big.Int
		balance.Set(paymentBalance).Abs(&balance)
		res.DriftBalances[asset] = &balance
		return fmt.Errorf("missing asset %s in ledgerBalances", asset)
	case ledgerBalance != nil && paymentBalance == nil:
		var balance big.Int
		balance.Set(ledgerBalance).Abs(&balance)
		res.DriftBalances[asset] = &balance
		res.DriftBalances[asset] = ledgerBalance
		return fmt.Errorf("missing asset %s in paymentBalances", asset)
	case ledgerBalance != nil && paymentBalance != nil:
		var drift big.Int
		drift.Set(paymentBalance).Add(&drift, ledgerBalance)

		var err error
		switch drift.Cmp(big.NewInt(0)) {
		case 0, 1:
		default:
			err = fmt.Errorf("balance drift for asset %s", asset)
		}

		res.DriftBalances[asset] = drift.Abs(&drift)
		return err
	}

	return nil
}

// Missing asset should be considered as asset with balance 0
func harmonizeBalances(ledgerBalances, paymentsBalances map[string]*big.Int) (map[string]*big.Int, map[string]*big.Int) {
	allAssets := make(map[string]struct{})
	for asset := range ledgerBalances {
		allAssets[asset] = struct{}{}
	}
	for asset := range paymentsBalances {
		allAssets[asset] = struct{}{}
	}

	for asset := range allAssets {
		if _, ok := ledgerBalances[asset]; !ok {
			ledgerBalances[asset] = big.NewInt(0)
		}
		if _, ok := paymentsBalances[asset]; !ok {
			paymentsBalances[asset] = big.NewInt(0)
		}
	}

	return ledgerBalances, paymentsBalances
}

func (s *Service) GetReconciliation(ctx context.Context, id string) (*models.Reconciliation, error) {
	rID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	reco, err := s.store.GetReconciliation(ctx, rID)
	return reco, newStorageError(err, "getting reconciliation")
}

func (s *Service) ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	reconciliations, err := s.store.ListReconciliations(ctx, q)
	return reconciliations, newStorageError(err, "listing reconciliations")
}
