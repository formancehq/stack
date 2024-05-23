package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/pkg/errors"
)

type CreateRuleRequest struct {
	Name string `json:"name"`

	RuleType models.RuleType `json:"ruleType"`
	Discard  bool            `json:"discard"`
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
	rule := &models.Rule{
		ID:             getRuleID(req.Rule),
		Name:           req.Name,
		CreatedAt:      time.Now(),
		Type:           req.RuleType,
		Discard:        req.Discard,
		RuleDefinition: req.Rule,
	}

	err := s.store.CreateRule(ctx, rule)
	if err != nil {
		return nil, newStorageError(err, "creating rule")
	}

	return rule, nil
}

func (s *Service) DeleteRule(ctx context.Context, id string) error {
	ruleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.Wrap(ErrValidation, "parsing rule ID")
	}

	return newStorageError(s.store.DeleteRule(ctx, uint32(ruleID)), "deleting rule")
}

func (s *Service) GetRule(ctx context.Context, id string) (*models.Rule, error) {
	ruleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, "parsing rule ID")
	}

	rule, err := s.store.GetRule(ctx, uint32(ruleID))
	if err != nil {
		return nil, newStorageError(err, "getting rule")
	}

	return rule, nil
}

func (s *Service) ListRules(ctx context.Context, q storage.ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error) {
	rules, err := s.store.ListRules(ctx, q)
	return rules, newStorageError(err, "listing rules")
}
