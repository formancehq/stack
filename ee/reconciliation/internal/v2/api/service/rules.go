package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/pkg/errors"
)

type CreateRuleRequest struct {
	Name string `json:"name"`

	RuleType models.RuleType `json:"ruleType"`
	Rule     json.RawMessage `json:"rule"`
}

func (r *CreateRuleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.RuleType == "" {
		return errors.New("missing ruleType")
	}

	if r.Rule == nil {
		return errors.New("missing rule")
	}

	switch r.RuleType {
	case models.RuleTypeMatchField:
		if _, err := models.ValidateRuleMatch(r.Rule); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	case models.RuleTypeAPICall:
		if _, err := models.ValidateRuleAPICall(r.Rule); err != nil {
			return errors.Wrap(err, "invalid rule definition")
		}
	default:
		return errors.New("unsupported rule type")
	}

	return nil
}

func (s *Service) CreateRule(ctx context.Context, req *CreateRuleRequest) (*models.Rule, error) {
	// TODO(polo): hash rule json to get id
	id := ""

	rule := &models.Rule{
		ID:             id,
		Name:           req.Name,
		CreatedAt:      time.Now(),
		Type:           req.RuleType,
		RuleDefinition: req.Rule,
	}

	err := s.store.CreateRule(ctx, rule)
	if err != nil {
		return nil, newStorageError(err, "creating rule")
	}

	return rule, nil
}

func (s *Service) DeleteRule(ctx context.Context, id string) error {
	return newStorageError(s.store.DeleteRule(ctx, id), "deleting rule")
}

func (s *Service) GetRule(ctx context.Context, id string) (*models.Rule, error) {
	rule, err := s.store.GetRule(ctx, id)
	if err != nil {
		return nil, newStorageError(err, "getting rule")
	}

	return rule, nil
}

func (s *Service) ListRules(ctx context.Context, q storage.ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error) {
	rules, err := s.store.ListRules(ctx, q)
	return rules, newStorageError(err, "listing rules")
}
