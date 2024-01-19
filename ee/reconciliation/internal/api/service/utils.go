package service

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"golang.org/x/mod/semver"
)

type Version interface {
	GetVersion() string
}

func isVersionSupported(
	version Version,
	minSupportedVersion string,
) bool {
	v := "v" + version.GetVersion()
	if !semver.IsValid(v) {
		// If semver is not valid, we assume it's a commit hash, so last version
		return true
	}

	switch semver.Compare(v, minSupportedVersion) {
	case 0, 1:
		// Higher or equal, nothing to do
		return true
	default:
		return false
	}
}

func (s *Service) getAccountsAggregatedBalance(ctx context.Context, ledgerName string, ledgerAggregatedBalanceQuery map[string]interface{}, at time.Time) (map[string]*big.Int, error) {
	infoResponse, err := s.client.V2GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	if infoResponse.StatusCode != 200 {
		return nil, errors.New("failed to get ledger info")
	}

	if !isVersionSupported(infoResponse.V2ConfigInfoResponse, "v2.0.0-beta.1") {
		return nil, errors.New("ledger version not supported")
	}

	balances, err := s.client.V2GetBalancesAggregated(
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
	response, err := s.client.PaymentsgetServerInfo(ctx)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("failed to get payments info")
	}

	if !isVersionSupported(response.ServerInfo, "v1.0.0-rc.4") {
		return nil, errors.New("payments version not supported")
	}

	balances, err := s.client.GetPoolBalances(
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
