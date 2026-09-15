package events

import (
	"fmt"
	"testing"
)

func TestReconciliationAlertEventSchemas(t *testing.T) {
	eventTypes := map[string]string{
		"OPENED_ALERT":       "fail",
		"UPDATED_ALERT":      "fail",
		"ACKNOWLEDGED_ALERT": "ack",
		"RESOLVED_ALERT":     "resolve",
		"ACCEPTED_ALERT":     "accept",
		"REOPENED_ALERT":     "fail",
		"SNOOZED_ALERT":      "snooze",
		"UNSNOOZED_ALERT":    "unsnooze",
	}

	for eventName, alertEventType := range eventTypes {
		t.Run(eventName, func(t *testing.T) {
			event := fmt.Sprintf(`{
				"app": "reconciliation",
				"version": "v1",
				"date": "2026-08-03T12:00:00Z",
				"type": %q,
				"payload": {
					"alert": {
						"id": "11111111-1111-1111-1111-111111111111",
						"ruleID": "22222222-2222-2222-2222-222222222222",
						"fingerprint": "asset:USD/2",
						"periodID": "2026-08",
						"status": "OPEN",
						"severity": "high",
						"firstSeenAt": "2026-08-03T12:00:00Z",
						"lastSeenAt": "2026-08-03T12:00:00Z",
						"occurrenceCount": 1,
						"lastEvaluationID": "33333333-3333-3333-3333-333333333333",
						"createdAt": "2026-08-03T12:00:00Z",
						"updatedAt": "2026-08-03T12:00:00Z"
					},
					"event": {
						"id": "44444444-4444-4444-4444-444444444444",
						"alertID": "11111111-1111-1111-1111-111111111111",
						"type": %q,
						"newStatus": "OPEN",
						"notify": true,
						"at": "2026-08-03T12:00:00Z",
						"createdAt": "2026-08-03T12:00:00Z"
					}
				}
			}`, eventName, alertEventType)

			if err := Check([]byte(event), "reconciliation", eventName); err != nil {
				t.Fatalf("validate event: %v", err)
			}
		})
	}
}

func TestLedgerV3EventSchemas(t *testing.T) {
	t.Parallel()

	events := map[string]string{
		"CREATED_LEDGER": `{
			"type":"CREATED_LEDGER",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":1,
			"log":{
				"sequence":1,
				"payload":{"createLedger":{"name":"orders","createdAt":"2026-08-28T22:29:41.133323Z"}},
				"responseSignature":{}
			}
		}`,
		"COMMITTED_TRANSACTION": `{
			"type":"COMMITTED_TRANSACTION",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":2,
			"log":{
				"sequence":2,
				"payload":{"apply":{"ledgerName":"orders","log":{
					"type":"NEW_TRANSACTION",
					"data":{"createdTransaction":{
						"transaction":{
							"postings":[{"source":"world","destination":"users:123","amount":1000,"asset":"USD/2","color":""}],
							"metadata":{"kind":"sale"},
							"timestamp":"2026-08-28T22:29:41.133323Z",
							"reference":"order-42",
							"id":42,
							"reverted":false
						},
						"accountMetadata":{"users:123":{"tier":"gold"}},
						"chapterId":7
					}},
					"date":"2026-08-28T22:29:41.133323Z",
					"id":1
				}}},
				"responseSignature":{}
			}
		}`,
		"REVERTED_TRANSACTION": `{
			"type":"REVERTED_TRANSACTION",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":3,
			"log":{
				"sequence":3,
				"payload":{"apply":{"ledgerName":"orders","log":{
					"type":"REVERTED_TRANSACTION",
					"data":{"revertedTransaction":{
						"revertedTransactionId":42,
						"revertTransaction":{
							"postings":[{"source":"users:123","destination":"world","amount":1000,"asset":"USD/2","color":""}],
							"metadata":{},
							"timestamp":"2026-08-28T22:29:41.133323Z",
							"id":43,
							"reverted":false
						}
					}},
					"date":"2026-08-28T22:29:41.133323Z",
					"id":2
				}}},
				"responseSignature":{}
			}
		}`,
		"SAVED_METADATA": `{
			"type":"SAVED_METADATA",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":4,
			"log":{
				"sequence":4,
				"payload":{"apply":{"ledgerName":"orders","log":{
					"type":"SET_METADATA",
					"data":{"savedMetadata":{"targetType":"ACCOUNT","accountId":"users:123","metadata":{"tier":"platinum"}}},
					"date":"2026-08-28T22:29:41.133323Z",
					"id":3
				}}},
				"responseSignature":{}
			}
		}`,
		"DELETED_METADATA": `{
			"type":"DELETED_METADATA",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":5,
			"log":{
				"sequence":5,
				"payload":{"apply":{"ledgerName":"orders","log":{
					"type":"DELETE_METADATA",
					"data":{"deletedMetadata":{"targetType":"TRANSACTION","transactionId":42,"key":"kind"}},
					"date":"2026-08-28T22:29:41.133323Z",
					"id":4
				}}},
				"responseSignature":{}
			}
		}`,
		"DELETED_LEDGER": `{
			"type":"DELETED_LEDGER",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":6,
			"log":{
				"sequence":6,
				"payload":{"deleteLedger":{"name":"orders","deletedAt":"2026-08-28T22:29:41.133323Z"}},
				"responseSignature":{}
			}
		}`,
		"SKIPPED_ORDER": `{
			"type":"SKIPPED_ORDER",
			"ledger":"orders",
			"date":"2026-08-28T22:29:41.133323Z",
			"logSequence":7,
			"log":{
				"sequence":7,
				"payload":{"apply":{"ledgerName":"orders","log":{
					"type":"ORDER_SKIPPED",
					"data":{"reason":"TRANSACTION_REFERENCE_CONFLICT","context":{"reference":"order-42","existingTransactionId":"42"}},
					"date":"2026-08-28T22:29:41.133323Z",
					"id":5
				}}},
				"responseSignature":{}
			}
		}`,
	}

	for eventName, event := range events {
		eventName, event := eventName, event
		t.Run(eventName, func(t *testing.T) {
			t.Parallel()

			if err := CheckForVersion([]byte(event), "ledger", "v3.0.0", eventName); err != nil {
				t.Fatalf("validate event: %v", err)
			}
		})
	}

	if err := CheckForVersion([]byte(events["COMMITTED_TRANSACTION"]), "ledger", "v3.0.0", "SAVED_METADATA"); err == nil {
		t.Fatal("expected a committed transaction to be rejected by the saved metadata schema")
	}
}

func TestComputeSchemaUsesLatestVersionContainingEvent(t *testing.T) {
	t.Parallel()

	legacyEvent := []byte(`{
		"app":"ledger",
		"version":"v2",
		"date":"2026-08-28T22:29:41Z",
		"type":"COMMITTED_TRANSACTIONS",
		"payload":{
			"ledger":"orders",
			"transactions":[{
				"postings":[],
				"metadata":{},
				"id":42,
				"timestamp":"2026-08-28T22:29:41Z",
				"reverted":false
			}]
		}
	}`)

	if err := Check(legacyEvent, "ledger", "COMMITTED_TRANSACTIONS"); err != nil {
		t.Fatalf("validate legacy event from newest version containing it: %v", err)
	}
}
