package worker

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/google/uuid"
)

// TODO(polo): Refine storage to actually represents a tree

type RuleTreeFor string

const (
	RuleTreeForLedger   RuleTreeFor = "ledger"
	RuleTreeForPayments RuleTreeFor = "payments"
)

type RuleTree struct {
	Rule        *models.Rule
	RuleTreeFor RuleTreeFor
	Children    []*RuleTree

	ruleMatch   *models.RuleMatch
	ruleAPICall *models.RuleAPICall
}

func (l *Listener) getRulesTree(ctx context.Context, ruleID uuid.UUID) (*RuleTree, error) {
	r, ok := l.ruleTreeCache[ruleID]
	if ok {
		return r, nil
	}

	rule, err := l.store.GetRule(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	ruleTree := &RuleTree{
		Rule: rule,
	}

	switch rule.Type {
	case models.RuleTypeFilters:
		var ruleFilters models.RuleFilters
		if err := json.Unmarshal(rule.RuleDefinition, &ruleFilters); err != nil {
			return nil, err
		}

		ruleTree.Children = make([]*RuleTree, len(ruleFilters.ChildrenRules))
		for i, childRuleID := range ruleFilters.ChildrenRules {
			childRuleTree, err := l.getRulesTree(ctx, childRuleID)
			if err != nil {
				return nil, err
			}

			ruleTree.Children[i] = childRuleTree
		}

	case models.RuleTypeOr:
		var ruleOr models.RuleFilters
		if err := json.Unmarshal(rule.RuleDefinition, &ruleOr); err != nil {
			return nil, err
		}

		ruleTree.Children = make([]*RuleTree, len(ruleOr.ChildrenRules))
		for i, childRuleID := range ruleOr.ChildrenRules {
			childRuleTree, err := l.getRulesTree(ctx, childRuleID)
			if err != nil {
				return nil, err
			}

			ruleTree.Children[i] = childRuleTree
		}

	case models.RuleTypeAPICall:
		var ruleAPICall models.RuleAPICall
		if err := json.Unmarshal(rule.RuleDefinition, &ruleAPICall); err != nil {
			return nil, err
		}

		ruleTree.ruleAPICall = &ruleAPICall

	case models.RuleTypeMatchField:
		var ruleMatch models.RuleMatch
		if err := json.Unmarshal(rule.RuleDefinition, &ruleMatch); err != nil {
			return nil, err
		}

		ruleTree.ruleMatch = &ruleMatch

	case models.RuleTypeAccountBased:
		return nil, errors.New("unsupported rule type: account based rules")
	case models.RuleTypeUnknown:
		return nil, errors.New("unknown rule type")
	}

	l.ruleTreeCache[ruleID] = ruleTree

	return ruleTree, nil
}

type processRulesTreeResults struct {
	ForLedger struct {
		ledgerTransactionID *big.Int
		paymentIDs          []string
	}
	ForPayment struct {
		paymentID            string
		ledgerTransactionIDs []*big.Int
	}
}

func (l *Listener) processRulesTreeRoot(
	ctx context.Context,
	r *RuleTree,
	event interface{},
	additionalFilters ...models.Filter,
) (*big.Int, uuid.UUID, error) {
	switch r.Rule.Type {
	case models.RuleTypeFilters:
		filters, err := l.processRulesTreeFilters(ctx, r, event, additionalFilters...)
		if err != nil {
			return nil, uuid.Nil, err
		}

		query, err := models.BuildESQuery(filters)
		if err != nil {
			return nil, uuid.Nil, err
		}

		response, err := l.client.RawSearch(ctx, query)
		if err != nil {
			return nil, uuid.Nil, err
		}

	}
	return nil, uuid.Nil, nil
}

func (l *Listener) processRulesTreeFilters(
	ctx context.Context,
	r *RuleTree,
	event interface{},
	additionalFilters ...models.Filter,
) ([]models.Filter, error) {
	filters := make([]models.Filter, 0, len(r.Children)+len(additionalFilters))
	filters = append(filters, additionalFilters...)

	for _, child := range r.Children {
		if child.Rule.Type != models.RuleTypeMatchField {
			return nil, errors.New("unsupported rule type for filters")
		}

		if child.ruleMatch == nil {
			return nil, errors.New("missing rule match")
		}

		switch r.RuleTreeFor {
		case RuleTreeForLedger:
			value, err := extractValue(event, child.ruleMatch.Ledger.Match.Key)
			if err != nil {
				return nil, err
			}

			filters = append(filters, child.ruleMatch.BuildLedgerFilters(value)...)
		case RuleTreeForPayments:
			value, err := extractValue(event, child.ruleMatch.Payment.Match.Key)
			if err != nil {
				return nil, err
			}

			filters = append(filters, child.ruleMatch.BuildPaymentFilters(value)...)
		}
	}

	return filters, nil
}
