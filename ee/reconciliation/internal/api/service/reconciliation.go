package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/reconciliation/internal/models"
	"golang.org/x/sync/errgroup"
)

type ReconciliationRequest struct {
	LedgerName                   string                 `json:"ledgerName"`
	LedgerAggregatedBalanceQuery map[string]interface{} `json:"ledgerAggregatedBalanceQuery"`
	PaymentPoolID                string                 `json:"paymentPoolID"`
	At                           time.Time              `json:"at"`
}

func (r *ReconciliationRequest) Validate() error {
	if r.LedgerName == "" {
		return errors.New("missing ledger name")
	}

	if r.LedgerAggregatedBalanceQuery == nil {
		return errors.New("missing ledger aggregated balance query")
	}

	if r.PaymentPoolID == "" {
		return errors.New("missing payments pool id")
	}

	if r.At.IsZero() {
		return errors.New("missing at")
	}

	if r.At.After(time.Now()) {
		return errors.New("at must be in the past")
	}

	return nil
}

type ReconciliationResponse struct {
	Status          models.ReconciliationStatus
	PaymentBalances map[string]*big.Int
	LedgerBalances  map[string]*big.Int
	Error           string
}

func (s *Service) Reconciliation(ctx context.Context, req *ReconciliationRequest) (*ReconciliationResponse, error) {
	eg, ctxGroup := errgroup.WithContext(ctx)

	var paymentBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		paymentBalance, err = s.getPaymentPoolBalance(ctxGroup, req.PaymentPoolID, req.At)
		return err
	})

	var ledgerBalance map[string]*big.Int
	eg.Go(func() error {
		var err error
		ledgerBalance, err = s.getAccountsAggregatedBalance(ctxGroup, req.LedgerName, req.LedgerAggregatedBalanceQuery, req.At)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	fmt.Println("paymentBalance")
	for asset, balance := range paymentBalance {
		fmt.Println(asset, balance)
	}

	fmt.Println("ledgerBalance")
	for asset, balance := range ledgerBalance {
		fmt.Println(asset, balance)
	}

	res := &ReconciliationResponse{
		Status:          models.ReconciliationOK,
		PaymentBalances: paymentBalance,
		LedgerBalances:  ledgerBalance,
	}

	if len(paymentBalance) != len(ledgerBalance) {
		res.Status = models.ReconciliationNotOK
		res.Error = "different number of assets"
		return res, nil
	}

	for paymentAsset, paymentBalance := range paymentBalance {
		ledgerBalance, ok := ledgerBalance[paymentAsset]
		if !ok {
			res.Status = models.ReconciliationNotOK
			res.Error = fmt.Sprintf("missing asset %s in ledger", paymentAsset)
			return res, nil
		}

		if paymentBalance.Cmp(ledgerBalance) != 0 {
			res.Status = models.ReconciliationNotOK
			res.Error = fmt.Sprintf("different balance for asset %s", paymentAsset)
			return res, nil
		}
	}

	return res, nil
}

func (s *Service) getAccountsAggregatedBalance(ctx context.Context, ledgerName string, ledgerAggregatedBalanceQuery map[string]interface{}, at time.Time) (map[string]*big.Int, error) {
	balances, err := s.client.Ledger.V2.GetBalancesAggregated(
		ctx,
		operations.GetBalancesAggregatedRequest{
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

	if balances.AggregateBalancesResponse == nil {
		return nil, errors.New("no aggregated balance")
	}

	balanceMap := make(map[string]*big.Int)
	for asset, balance := range balances.AggregateBalancesResponse.Data {
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
