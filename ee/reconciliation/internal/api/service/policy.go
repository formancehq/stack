package service

import (
	"context"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"

	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreatePolicyRequest struct {
	Name           string                 `json:"name"`
	LedgerName     string                 `json:"ledgerName"`
	LedgerQuery    map[string]interface{} `json:"ledgerQuery"`
	PaymentsPoolID string                 `json:"paymentsPoolID"`
}

func (r *CreatePolicyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.LedgerName == "" {
		return errors.New("missing ledgerName")
	}

	if r.PaymentsPoolID == "" {
		return errors.New("missing paymentsPoolId")
	}

	return nil
}

func (s *Service) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*models.Policy, error) {
	paymentPoolID, err := uuid.Parse(req.PaymentsPoolID)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	policy := &models.Policy{
		ID:             uuid.New(),
		Name:           req.Name,
		CreatedAt:      time.Now().UTC(),
		LedgerName:     req.LedgerName,
		LedgerQuery:    req.LedgerQuery,
		PaymentsPoolID: paymentPoolID,
	}

	err = s.store.CreatePolicy(ctx, policy)
	if err != nil {
		return nil, newStorageError(err, "creating policy")
	}

	return policy, nil
}

func (s *Service) DeletePolicy(ctx context.Context, id string) error {
	pID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	return newStorageError(s.store.DeletePolicy(ctx, pID), "deleting policy")
}

func (s *Service) GetPolicy(ctx context.Context, id string) (*models.Policy, error) {
	pID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	policy, err := s.store.GetPolicy(ctx, pID)
	if err != nil {
		return nil, newStorageError(err, "getting policy")
	}

	return policy, nil
}

func (s *Service) ListPolicies(ctx context.Context, q storage.GetPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error) {
	policies, err := s.store.ListPolicies(ctx, q)
	return policies, newStorageError(err, "listing policies")
}
