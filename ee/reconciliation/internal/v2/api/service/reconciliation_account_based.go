package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetAccountBasedReconciliationDetails(ctx context.Context, id string) (*models.ReconciliationAccountBased, error) {
	reconciliationID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	accountBasedReconciliation, err := s.store.GetAccountBasedReconciliation(ctx, reconciliationID)
	if err != nil {
		return nil, newStorageError(err, "getting account based reconciliation")
	}

	return accountBasedReconciliation, nil

}

func (s *Service) handleAccountBasedReconciliation(
	ctx context.Context,
	req *CreateAccountBasedReconciliationRequest,
	policy *models.Policy,
	reconciliation *models.Reconciliation,
) error {
	eg, ctxGroup := errgroup.WithContext(ctx)

	accountBasedRule, err := s.getAccountBasedRule(ctx, policy.Rules[0])
	if err != nil {
		return errors.Wrap(err, "getting account-based rule")
	}

	var paymentBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		paymentBalance, err = s.getPaymentPoolBalance(ctxGroup, accountBasedRule.Payments.PoolID.String(), req.ReconciledAtPayments)
		return err
	})

	var ledgerBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		ledgerBalance, err = s.getAccountsAggregatedBalance(ctxGroup, accountBasedRule.Ledger.Name, accountBasedRule.Ledger.Query, req.ReconciledAtLedger)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	res := models.ReconciliationAccountBased{
		ReconciliationID:     reconciliation.ID,
		ReconciledAtLedger:   req.ReconciledAtLedger,
		ReconciledAtPayments: req.ReconciledAtPayments,
		ReconciliationStatus: models.ReconciliationStatusSucceeded,
		LedgerBalances:       ledgerBalance,
		PaymentsBalances:     paymentBalance,
		DriftBalances:        make(map[string]*big.Int),
	}

	if len(paymentBalance) != len(ledgerBalance) {
		reconciliation.Status = models.ReconciliationStatusFailed
		res.Error = "different number of assets"
	} else {
		for asset, ledgerBalance := range ledgerBalance {
			err := s.computeDrift(&res, asset, ledgerBalance, paymentBalance[asset])
			if err != nil {
				reconciliation.Status = models.ReconciliationStatusFailed
				if res.Error == "" {
					res.Error = err.Error()
				} else {
					res.Error = res.Error + "; " + err.Error()
				}
			}
		}

		for asset, paymentBalance := range paymentBalance {
			if _, ok := res.DriftBalances[asset]; ok {
				// Already computed
				continue
			}

			err := s.computeDrift(&res, asset, ledgerBalance[asset], paymentBalance)
			if err != nil {
				reconciliation.Status = models.ReconciliationStatusFailed
				res.Error = res.Error + "; " + err.Error()
			}
		}
	}

	if err := s.store.CreateAccountBasedReconciliation(ctx, &res); err != nil {
		return newStorageError(err, "failed to create account based reconciliation")
	}

	return nil
}

func (s *Service) computeDrift(
	res *models.ReconciliationAccountBased,
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

func (s *Service) getAccountBasedRule(ctx context.Context, ruleID uint32) (*models.RuleAccountBased, error) {
	rule, err := s.store.GetRule(ctx, ruleID)
	if err != nil {
		return nil, newStorageError(err, "getting rule")
	}

	var accountBasedRule models.RuleAccountBased
	if err := json.Unmarshal(rule.RuleDefinition, &accountBasedRule); err != nil {
		return nil, errors.Wrap(err, "unmarshalling rule")
	}

	return &accountBasedRule, nil
}
