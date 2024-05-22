package service

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Service) GetReconciliation(ctx context.Context, id string) (*models.Reconciliation, error) {
	reconciliationID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	reconciliation, err := s.store.GetReconciliation(ctx, reconciliationID)
	if err != nil {
		return nil, newStorageError(err, "getting reconciliation")
	}

	return reconciliation, nil
}

func (s *Service) ListReconciliations(ctx context.Context, q storage.ListReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	reconciliations, err := s.store.ListReconciliations(ctx, q)
	return reconciliations, newStorageError(err, "listing reconciliations")
}

type CreateAccountBasedReconciliationRequest struct {
	ReconciledAtLedger   time.Time `json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time `json:"reconciledAtPayments"`
}

func (r *CreateAccountBasedReconciliationRequest) Validate() error {
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

type CreateReconciliationRequest struct {
	Name     string    `json:"name"`
	PolicyID uuid.UUID `json:"policyID"`

	CreateAccountBasedReconciliationRequest *CreateAccountBasedReconciliationRequest `json:"account_based"`
}

func (r *CreateReconciliationRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	if r.PolicyID == uuid.Nil {
		return fmt.Errorf("policyID is required")
	}

	if r.CreateAccountBasedReconciliationRequest != nil {
		if err := r.CreateAccountBasedReconciliationRequest.Validate(); err != nil {
			return errors.Wrap(err, "invalid account based reconciliation request")
		}
	}

	return nil
}

func (s *Service) CreateReconciliation(ctx context.Context, req *CreateReconciliationRequest) (*models.Reconciliation, error) {
	policy, err := s.store.GetPolicy(ctx, req.PolicyID)
	if err != nil {
		return nil, newStorageError(err, "getting policy")
	}

	reconciliation := &models.Reconciliation{
		ID:         uuid.New(),
		Name:       req.Name,
		PolicyID:   req.PolicyID,
		PolicyType: policy.Type,
		CreatedAt:  time.Now().UTC(),
	}

	switch policy.Type {
	case models.PolicyTypeAccountBased:
		if req.CreateAccountBasedReconciliationRequest == nil {
			return nil, fmt.Errorf("missing account based reconciliation request")
		}

		if err := s.handleAccountBasedReconciliation(ctx, req.CreateAccountBasedReconciliationRequest, policy, reconciliation); err != nil {
			return nil, errors.Wrap(err, "handling account based reconciliation")
		}
		return reconciliation, nil

	case models.PolicyTypeTransactionBased:
		if err := s.handleTransactionBasedReconciliation(ctx, reconciliation); err != nil {
			return nil, errors.Wrap(err, "handling transaction based reconciliation")
		}
		return reconciliation, nil

	default:
		return nil, fmt.Errorf("unsupported policy type: %s", policy.Type)
	}
}

func (s *Service) handleTransactionBasedReconciliation(ctx context.Context, reconciliation *models.Reconciliation) error {
	return nil
}
