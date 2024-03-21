package delay

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/orchestration/internal/schema"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
)

func TestDelaySchema(t *testing.T) {
	now := time.Now().Round(time.Second).UTC()
	stagestesting.TestSchemas(t, "delay", []stagestesting.SchemaTestCase{
		{
			Name: "valid case using until property",
			Data: map[string]any{
				"until": now.Format(time.RFC3339Nano),
			},
			ExpectedResolved: Delay{
				Until: &now,
			},
			ExpectedValidationError: false,
		},
		{
			Name: "valid case using duration",
			Data: map[string]any{
				"duration": "10s",
			},
			ExpectedValidationError: false,
			ExpectedResolved: Delay{
				Duration: (*schema.Duration)(ptr(time.Second * 10)),
			},
		},
		{
			Name:                    "invalid case, missing until or duration",
			Data:                    map[string]any{},
			ExpectedResolved:        Delay{},
			ExpectedValidationError: true,
		},
		{
			Name: "invalid case, both until and duration specified",
			Data: map[string]any{
				"until":    now.Format(time.RFC3339Nano),
				"duration": "10s",
			},
			ExpectedResolved: Delay{
				Duration: (*schema.Duration)(ptr(time.Second * 10)),
				Until:    &now,
			},
			ExpectedValidationError: true,
		},
	}...)
}

func ptr[T any](v T) *T {
	return &v
}

var testCases = []stagestesting.WorkflowTestCase[Delay]{
	{
		Stage: Delay{
			Until: ptr(time.Now().Add(time.Second)),
		},
		Name: "delay-until",
	},
	{
		Stage: Delay{
			Duration: (*schema.Duration)(ptr(time.Second)),
		},
		Name: "delay-duration",
	},
}

func TestDelay(t *testing.T) {
	stagestesting.RunWorkflows(t, testCases...)
}
