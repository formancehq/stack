package service

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreatePoolRequest struct {
	Name       string   `json:"name"`
	AccountIDs []string `json:"accountIDs"`
}

func (c *CreatePoolRequest) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if len(c.AccountIDs) == 0 {
		return errors.New("accountIDs is required")
	}

	return nil
}

func (s *Service) CreatePool(
	ctx context.Context,
	req *CreatePoolRequest,
) (*models.Pool, error) {
	pool := &models.Pool{
		Name:      req.Name,
		CreatedAt: time.Now().UTC(),
	}

	err := s.store.CreatePool(ctx, pool)
	if err != nil {
		return nil, newStorageError(err, "creating pool")
	}

	poolAccounts := make([]*models.PoolAccounts, len(req.AccountIDs))
	for i, accountID := range req.AccountIDs {
		aID, err := models.AccountIDFromString(accountID)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}

		poolAccounts[i] = &models.PoolAccounts{
			PoolID:    pool.ID,
			AccountID: *aID,
		}
	}

	err = s.store.AddAccountsToPool(ctx, poolAccounts)
	if err != nil {
		return nil, newStorageError(err, "adding accounts to pool")
	}
	pool.PoolAccounts = poolAccounts

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventSavedPool(pool)))
	if err != nil {
		return nil, errors.Wrap(err, "publishing message")
	}

	return pool, nil
}

type AddAccountToPoolRequest struct {
	AccountID string `json:"accountID"`
}

func (c *AddAccountToPoolRequest) Validate() error {
	if c.AccountID == "" {
		return errors.New("accountID is required")
	}

	return nil
}

func (s *Service) AddAccountToPool(
	ctx context.Context,
	poolID string,
	req *AddAccountToPoolRequest,
) error {
	id, err := uuid.Parse(poolID)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	aID, err := models.AccountIDFromString(req.AccountID)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	if err := s.store.AddAccountToPool(ctx, &models.PoolAccounts{
		PoolID:    id,
		AccountID: *aID,
	}); err != nil {
		return newStorageError(err, "adding account to pool")
	}

	pool, err := s.store.GetPool(ctx, id)
	if err != nil {
		return newStorageError(err, "getting pool")
	}

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventSavedPool(pool)))
	if err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return nil
}

func (s *Service) RemoveAccountFromPool(
	ctx context.Context,
	poolID string,
	accountID string,
) error {
	id, err := uuid.Parse(poolID)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	aID, err := models.AccountIDFromString(accountID)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	if err := s.store.RemoveAccountFromPool(ctx, &models.PoolAccounts{
		PoolID:    id,
		AccountID: *aID,
	}); err != nil {
		return newStorageError(err, "removing account from pool")
	}

	pool, err := s.store.GetPool(ctx, id)
	if err != nil {
		return newStorageError(err, "getting pool")
	}

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventSavedPool(pool)))
	if err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return nil
}

func (s *Service) ListPools(ctx context.Context, q storage.ListPoolsQuery) (*api.Cursor[models.Pool], error) {
	cursor, err := s.store.ListPools(ctx, q)
	return cursor, newStorageError(err, "listing pools")
}

func (s *Service) GetPool(
	ctx context.Context,
	poolID string,
) (*models.Pool, error) {
	id, err := uuid.Parse(poolID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	pool, err := s.store.GetPool(ctx, id)
	return pool, newStorageError(err, "getting pool")
}

type GetPoolBalanceResponse struct {
	Balances []*Balance
}

type Balance struct {
	Amount *big.Int
	Asset  string
}

func (s *Service) GetPoolBalance(
	ctx context.Context,
	poolID string,
	atTime string,
) (*GetPoolBalanceResponse, error) {
	id, err := uuid.Parse(poolID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	at, err := time.Parse(time.RFC3339, atTime)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	pool, err := s.store.GetPool(ctx, id)
	if err != nil {
		return nil, newStorageError(err, "getting pool")
	}

	res := make(map[string]*big.Int)
	for _, poolAccount := range pool.PoolAccounts {
		balances, err := s.store.GetBalancesAt(ctx, poolAccount.AccountID, at)
		if err != nil {
			return nil, newStorageError(err, "getting balances")
		}

		for _, balance := range balances {
			amount, ok := res[balance.Asset.String()]
			if !ok {
				amount = big.NewInt(0)
			}

			amount.Add(amount, balance.Balance)
			res[balance.Asset.String()] = amount
		}
	}

	balances := make([]*Balance, 0, len(res))
	for asset, amount := range res {
		balances = append(balances, &Balance{
			Asset:  asset,
			Amount: amount,
		})
	}

	return &GetPoolBalanceResponse{
		Balances: balances,
	}, nil
}

func (s *Service) DeletePool(
	ctx context.Context,
	poolID string,
) error {
	id, err := uuid.Parse(poolID)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	if err := s.store.DeletePool(ctx, id); err != nil {
		return newStorageError(err, "deleting pool")
	}

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventDeletePool(id)))
	if err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return nil
}
