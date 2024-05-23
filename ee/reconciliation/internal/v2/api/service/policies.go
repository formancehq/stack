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

type CreatePolicyRequest struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Rules []uint32 `json:"rules"`
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

	if len(r.Rules) == 0 {
		return errors.New("missing rules list")
	}

	return nil
}

func (s *Service) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*models.Policy, error) {
	policyType := models.PolicyType(req.Type)

	switch policyType {
	case models.PolicyTypeAccountBased:
		if len(req.Rules) != 1 {
			return nil, errors.New("account-based policy must have exactly one rule")
		}

		rule, err := s.store.GetRule(ctx, req.Rules[0])
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
		Enabled:   true,
		Rules:     req.Rules,
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
