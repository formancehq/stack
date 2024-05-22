package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type RuleType string

const (
	RuleTypeUnknown      RuleType = "UNKNOWN"
	RuleTypeCustom       RuleType = "CUSTOM"
	RuleTypeAccountBased RuleType = "ACCOUNT_BASED"
	RuleTypeMatchField   RuleType = "MATCH"
	RuleTypeAPICall      RuleType = "API_CALL"
)

type Rule struct {
	bun.BaseModel `bun:"reconciliationsv2.rules" json:"-"`

	ID        string    `bun:",pk,notnull" json:"id"`
	Name      string    `bun:",notnull" json:"name"`
	CreatedAt time.Time `bun:",notnull" json:"createdAt"`
	RuleType  RuleType  `bun:",notnull" json:"ruleType"`

	RuleDefinition json.RawMessage `bun:",type:jsonb,notnull" json:"ruleDefinition"`
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
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (kv *KeyValue) Validate() error {
	if kv.Key == "" {
		return errors.New("missing key")
	}

	if kv.Value == nil {
		return errors.New("missing value")
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

type RuleAPICall struct {
	Endpoint string
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
