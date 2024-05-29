package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateRuleRequest struct {
	Name string `json:"name"`

	RuleType       models.RuleType `json:"ruleType"`
	Discard        bool            `json:"discard"`
	RuleDefinition json.RawMessage `json:"rule"`
}

func (r *CreateRuleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.RuleType == "" {
		return errors.New("missing ruleType")
	}

	if r.RuleDefinition == nil {
		return errors.New("missing rules")
	}

	switch r.RuleType {
	case models.RuleTypeFilters:
		if _, err := models.ValidateRuleFilters(r.RuleDefinition); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	case models.RuleTypeOr:
		if _, err := models.ValidateRuleOr(r.RuleDefinition); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	case models.RuleTypeMatchField:
		if _, err := models.ValidateRuleMatch(r.RuleDefinition); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	case models.RuleTypeAPICall:
		if _, err := models.ValidateRuleAPICall(r.RuleDefinition); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	default:
		return errors.New("unsupported rule type")
	}

	return nil
}

func (s *Service) CreateRule(ctx context.Context, req *CreateRuleRequest) (*models.Rule, error) {
	rule := &models.Rule{
		ID:             uuid.New(),
		Name:           req.Name,
		CreatedAt:      time.Now(),
		Type:           req.RuleType,
		Discard:        req.Discard,
		RuleDefinition: req.RuleDefinition,
	}

	err := s.store.CreateRule(ctx, rule)
	if err != nil {
		return nil, newStorageError(err, "creating rule")
	}

	return rule, nil
}

func (s *Service) DeleteRule(ctx context.Context, id string) error {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrap(ErrValidation, "parsing rule ID")
	}

	return newStorageError(s.store.DeleteRule(ctx, ruleID), "deleting rule")
}

func (s *Service) GetRule(ctx context.Context, id string) (*models.Rule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, "parsing rule ID")
	}

	rule, err := s.store.GetRule(ctx, ruleID)
	if err != nil {
		return nil, newStorageError(err, "getting rule")
	}

	return rule, nil
}

func (s *Service) ListRules(ctx context.Context, q storage.ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error) {
	rules, err := s.store.ListRules(ctx, q)
	return rules, newStorageError(err, "listing rules")
}
