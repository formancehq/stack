package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/stack/libs/go-libs/api/apitesting"
	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/go-chi/chi/v5"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestTestTrigger(t *testing.T) {

	type testCase struct {
		name     string
		data     triggers.TriggerData
		expected triggers.TestTriggerResult
		event    map[string]any
	}

	for _, testCase := range []testCase{
		{
			name:     "nominal",
			data:     triggers.TriggerData{},
			expected: triggers.TestTriggerResult{},
		},
		{
			name: "with filter and variable",
			event: map[string]any{
				"user": map[string]any{
					"role": "admin",
				},
			},
			data: triggers.TriggerData{
				Filter: pointer.For("event.user.role == \"admin\""),
				Vars: map[string]string{
					"role": "event.user.role",
				},
			},
			expected: triggers.TestTriggerResult{
				Filter: &triggers.FilterEvaluationResult{
					Match: true,
				},
				Variables: map[string]triggers.VariableEvaluationResult{
					"role": {
						Value: "admin",
					},
				},
			},
		},
		{
			name: "error on filter",
			data: triggers.TriggerData{
				Filter: pointer.For("a + b"),
			},
			expected: triggers.TestTriggerResult{
				Filter: &triggers.FilterEvaluationResult{
					Match: false,
					Error: "invalid operation: <nil> + <nil> (1:3)\n | a + b\n | ..^",
				},
			},
		},
		{
			name: "error on variables",
			data: triggers.TriggerData{
				Vars: map[string]string{
					"test": "a + b",
				},
			},
			expected: triggers.TestTriggerResult{
				Variables: map[string]triggers.VariableEvaluationResult{
					"test": {
						Error: "invalid operation: <nil> + <nil> (1:3)\n | a + b\n | ..^",
					},
				},
			},
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {

				w, err := m.Create(context.Background(), workflow.Config{})
				require.NoError(t, err)

				testCase.data.WorkflowID = w.ID
				testCase.data.Event = "XXX"
				trigger, err := m.CreateTrigger(context.Background(), testCase.data)
				require.NoError(t, err)

				payload := []byte("{}")
				if testCase.event != nil {
					payload, err = json.Marshal(testCase.event)
					require.NoError(t, err)
				}

				req := httptest.NewRequest(http.MethodPost, "/triggers/"+trigger.ID+"/test",
					bytes.NewBuffer(payload))
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				require.Equal(t, http.StatusOK, rec.Result().StatusCode)

				result := triggers.TestTriggerResult{}
				apitesting.ReadResponse(t, rec, &result)

				require.Equal(t, testCase.expected, result)
			})
		})
	}
}
