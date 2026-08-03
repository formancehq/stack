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
