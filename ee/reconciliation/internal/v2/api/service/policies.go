package service

import (
	"context"
	"time"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AdditionalConfig struct {
	LedgerOverdraftAccountReference string `json:"ledgerOverdraftAccountReference"`
}

type CreatePolicyRequest struct {
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	AdditionalConfig AdditionalConfig `json:"additionalConfig"`
	Rule             uuid.UUID        `json:"rule"`
}

func (r *CreatePolicyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.Type == "" {
		return errors.New("missing type")
	}

	_, err := models.PolicyTypeFromString(r.Type)
	if err != nil {
		return errors.Wrap(err, "invalid type")
	}

	return nil
}

func (s *Service) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*models.Policy, error) {
	policyType := models.PolicyType(req.Type)

	switch policyType {
	case models.PolicyTypeAccountBased:
		rule, err := s.store.GetRule(ctx, req.Rule)
		if err != nil {
			return nil, newStorageError(err, "getting rule")
		}

		if rule.Type != models.RuleTypeAccountBased {
			return nil, errors.New("rule must be account-based")
		}
	}

	now := time.Now()
	policy := &models.Policy{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
		Type:      policyType,
		AdditionalConfig: models.AdditionalConfig(
			req.AdditionalConfig.LedgerOverdraftAccountReference,
		),
		Enabled: true,
		Rule:    req.Rule,
	}

	err := s.store.CreatePolicy(ctx, policy)
	if err != nil {
		return nil, newStorageError(err, "creating policy")
	}

	return policy, nil
}

type UpdatePolicyRulesRequest struct {
	Rules []string `json:"rules"`
}

func (r *UpdatePolicyRulesRequest) Validate() error {
	if len(r.Rules) == 0 {
		return errors.New("missing rules")
	}

	return nil
}

func (s *Service) UpdatePolicyRules(ctx context.Context, id string, req *UpdatePolicyRulesRequest) error {
	policyID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	return newStorageError(s.store.UpdatePolicyRules(ctx, policyID, req.Rules), "failed to update policy rules")
}

func (s *Service) EnablePolicy(ctx context.Context, id string) error {
	policyID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	return newStorageError(s.store.UpdatePolicyStatus(ctx, policyID, true), "failed to enable policy")
}

func (s *Service) DisablePolicy(ctx context.Context, id string) error {
	policyID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	return newStorageError(s.store.UpdatePolicyStatus(ctx, policyID, false), "failed to disable policy")
}

func (s *Service) DeletePolicy(ctx context.Context, id string) error {
	policyID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	return newStorageError(s.store.DeletePolicy(ctx, policyID), "failed to delete policy")
}

func (s *Service) GetPolicy(ctx context.Context, id string) (*models.Policy, error) {
	policyID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	policy, err := s.store.GetPolicy(ctx, policyID)
	if err != nil {
		return nil, newStorageError(err, "getting policy")
	}

	return policy, nil
}

func (s *Service) ListPolicies(ctx context.Context, q storage.ListPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error) {
	policies, err := s.store.ListPolicies(ctx, q)
	return policies, newStorageError(err, "listing policies")
}
