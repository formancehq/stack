package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type RuleType string

const (
	RuleTypeUnknown      RuleType = "UNKNOWN"
	RuleTypeFilters      RuleType = "FILTERS"
	RuleTypeOr           RuleType = "OR"
	RuleTypeAccountBased RuleType = "ACCOUNT_BASED"
	RuleTypeMatchField   RuleType = "MATCH"
	RuleTypeAPICall      RuleType = "API_CALL"
)

type Rule struct {
	bun.BaseModel `bun:"reconciliationsv2.rules" json:"-"`

	ID             uuid.UUID       `bun:",pk,notnull" json:"id"`
	Name           string          `bun:",notnull" json:"name"`
	CreatedAt      time.Time       `bun:",notnull" json:"createdAt"`
	Type           RuleType        `bun:",notnull" json:"ruleType"`
	Discard        bool            `bun:",notnull" json:"discard"`
	RuleDefinition json.RawMessage `bun:",type:jsonb,notnull" json:"ruleDefinition"`
}

type RuleFilters struct {
	ChildrenRules []uuid.UUID `json:"childrenRules"`
}

func ValidateRuleFilters(rule json.RawMessage) (RuleFilters, error) {
	var ruleAnd RuleFilters
	if err := json.Unmarshal(rule, &ruleAnd); err != nil {
		return RuleFilters{}, errors.New("cannot unmarshal rule into ruleAnd")
	}

	if len(ruleAnd.ChildrenRules) == 0 {
		return RuleFilters{}, errors.New("missing childrenRules in ruleAnd")
	}

	return ruleAnd, nil
}

type RuleOr struct {
	ChildrenRules []uuid.UUID `json:"childrenRules"`
}

func ValidateRuleOr(rule json.RawMessage) (RuleOr, error) {
	var ruleOr RuleOr
	if err := json.Unmarshal(rule, &ruleOr); err != nil {
		return RuleOr{}, errors.New("cannot unmarshal rule into ruleOr")
	}

	if len(ruleOr.ChildrenRules) == 0 {
		return RuleOr{}, errors.New("missing childrenRules in ruleOr")
	}

	return ruleOr, nil
}

type RuleAccountBased struct {
	Ledger struct {
		Name  string                 `json:"name"`
		Query map[string]interface{} `json:"query"`
	} `json:"ledger"`
	Payments struct {
		PoolID uuid.UUID `json:"poolID"`
	} `json:"payments"`
}

type KeyValue struct {
	Key        string `json:"key"`
	IsMetadata bool   `json:"isMetadata"`
}

func (kv *KeyValue) Validate() error {
	if kv.Key == "" {
		return errors.New("missing key")
	}

	return nil
}

type RuleMatch struct {
	Ledger struct {
		Name  string   `json:"name"`
		Match KeyValue `json:"match"`
	} `json:"ledger"`
	Payment struct {
		ConnectorID string   `json:"connectorID"`
		Match       KeyValue `json:"match"`
	}
}

func ValidateRuleMatch(rule json.RawMessage) (RuleMatch, error) {
	var ruleMatch RuleMatch
	if err := json.Unmarshal(rule, &ruleMatch); err != nil {
		return RuleMatch{}, errors.New("cannot unmarshal rule into ruleMatch")
	}

	if ruleMatch.Ledger.Name == "" {
		return RuleMatch{}, errors.New("missing name in ledger rule match")
	}

	if err := ruleMatch.Ledger.Match.Validate(); err != nil {
		return RuleMatch{}, errors.Wrap(err, "invalid key value in ledger rule match")
	}

	if ruleMatch.Payment.ConnectorID == "" {
		return RuleMatch{}, errors.New("missing name in payment rule match")
	}

	if err := ruleMatch.Payment.Match.Validate(); err != nil {
		return RuleMatch{}, errors.New("invalid key value in payment rule match")
	}

	return ruleMatch, nil
}

func (r *RuleMatch) BuildLedgerFilters(value interface{}) []Filter {
	filters := []Filter{
		{
			MatchPhrase: map[string]interface{}{
				"indexed.ledger": r.Ledger.Name,
			},
		},
	}

	if r.Ledger.Match.IsMetadata {
		// TODO(polo): Metadata are not indexed yet, we need to index a special
		// type of metadata
	} else {
		filters = append(filters, Filter{
			MatchPhrase: map[string]interface{}{
				fmt.Sprintf("indexed.%s", r.Ledger.Match.Key): value,
			},
		})
	}

	return filters
}

func (r *RuleMatch) BuildPaymentFilters(value interface{}) []Filter {
	filters := []Filter{
		{
			MatchPhrase: map[string]interface{}{
				"indexed.connectorID": r.Payment.ConnectorID,
			},
		},
	}

	if r.Payment.Match.IsMetadata {
		// TODO(polo): Metadata are not indexed yet, we need to index a special
		// type of metadata
	} else {
		filters = append(filters, Filter{
			MatchPhrase: map[string]interface{}{
				fmt.Sprintf("indexed.%s", r.Payment.Match.Key): value,
			},
		})
	}

	return filters
}

type RuleAPICall struct {
	Endpoint string

	// Transaction
	// Payment
	// Reponse -> Transaction ou payment
}

func ValidateRuleAPICall(rule json.RawMessage) (RuleAPICall, error) {
	var ruleAPICall RuleAPICall
	if err := json.Unmarshal(rule, &ruleAPICall); err != nil {
		return RuleAPICall{}, errors.New("cannot unmarshal rule into ruleAPICall")
	}

	if ruleAPICall.Endpoint == "" {
		return RuleAPICall{}, errors.New("missing endpoint in ruleAPICall")
	}

	return ruleAPICall, nil
}
