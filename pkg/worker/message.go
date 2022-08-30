package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	ledger "github.com/numary/ledger/pkg/bus"
	payments "github.com/numary/payments/pkg"
	paymentIngestion "github.com/numary/payments/pkg/bridge/ingestion"
)

type EventMessage struct {
	Date    time.Time       `json:"date"`
	App     string          `json:"app"`
	Version string          `json:"version"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var ErrUnknownEventType = errors.New("unknown event type")

const (
	PrefixLedger   = "ledger"
	PrefixPayments = "payments"
)

func FilterMessage(msgValue []byte) (string, error) {
	var ev EventMessage
	if err := json.Unmarshal(msgValue, &ev); err != nil {
		return "", fmt.Errorf("json.Unmarshal event message: %w", err)
	}
	eventType := strcase.ToSnake(ev.Type)

	switch ev.Type {
	case ledger.EventTypeCommittedTransactions:
		committedTxs := new(ledger.CommittedTransactions)
		if err := json.Unmarshal(ev.Payload, committedTxs); err != nil {
			return "", fmt.Errorf("json.Unmarshal event message payload: %w", err)
		}
		eventType = strings.Join([]string{PrefixLedger, eventType}, ".")
		fmt.Printf("\nEVENT FETCHED: %s\n%+v\n", eventType, committedTxs)
	case ledger.EventTypeSavedMetadata:
		metadata := new(ledger.SavedMetadata)
		if err := json.Unmarshal(ev.Payload, metadata); err != nil {
			return "", fmt.Errorf("json.Unmarshal event message payload: %w", err)
		}
		eventType = strings.Join([]string{PrefixLedger, eventType}, ".")
		fmt.Printf("\nEVENT FETCHED: %s\n%+v\n", eventType, metadata)
	case ledger.EventTypeUpdatedMapping:
		mapping := new(ledger.UpdatedMapping)
		if err := json.Unmarshal(ev.Payload, mapping); err != nil {
			return "", fmt.Errorf("json.Unmarshal event message payload: %w", err)
		}
		eventType = strings.Join([]string{PrefixLedger, eventType}, ".")
		fmt.Printf("\nEVENT FETCHED: %s\n%+v\n", eventType, mapping)
	case ledger.EventTypeRevertedTransaction:
		revertedTx := new(ledger.RevertedTransaction)
		if err := json.Unmarshal(ev.Payload, revertedTx); err != nil {
			return "", fmt.Errorf("json.Unmarshal event message payload: %w", err)
		}
		eventType = strings.Join([]string{PrefixLedger, eventType}, ".")
		fmt.Printf("\nEVENT FETCHED: %s\n%+v\n", eventType, revertedTx)
	case paymentIngestion.EventTypeSavedPayment:
		savedPayment := new(payments.SavedPayment)
		if err := json.Unmarshal(ev.Payload, savedPayment); err != nil {
			return "", fmt.Errorf("json.Unmarshal event message payload: %w", err)
		}
		eventType = strings.Join([]string{PrefixPayments, eventType}, ".")
		fmt.Printf("\nEVENT FETCHED: %s\n%+v\n", eventType, savedPayment)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownEventType, ev.Type)
	}

	return ev.Type, nil
}
