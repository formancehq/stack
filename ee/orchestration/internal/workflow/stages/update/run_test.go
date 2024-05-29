package send

import (
	"testing"

	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/stretchr/testify/mock"

	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
)

func TestUpdateSchemaValidation(t *testing.T) {
	stagestesting.TestSchemas(t, "update", []stagestesting.SchemaTestCase{
		{
			Name: "nominal",
			Data: map[string]any{
				"account": map[string]interface{}{
					"id":     "merchant:${merchantID}",
					"ledger": "demo",
					"metadata": map[string]string{
						"formanceAccountID": "${paymentAccountID}",
					},
				},
			},
			ExpectedResolved: Update{
				Account: &AccountUpdate{
					ID:     "merchant:1234",
					Ledger: "demo",
					Metadata: map[string]string{
						"formanceAccountID": "ABCD",
					},
				},
			},
			Variables: map[string]string{
				"merchantID":       "1234",
				"paymentAccountID": "ABCD",
			},
		},
	}...)
}

func TestUpdate(t *testing.T) {
	stagestesting.RunWorkflows(t,
		stagestesting.WorkflowTestCase[Update]{
			Stage: Update{
				Account: &AccountUpdate{
					ID:     "abcd",
					Ledger: "default",
					Metadata: map[string]string{
						"foo": "bar",
					},
				},
			},
			MockedActivities: []stagestesting.MockedActivity{
				{
					Activity: activities.AddAccountMetadataActivity,
					Args: []any{mock.Anything, activities.AddAccountMetadataRequest{
						Account: "abcd",
						Ledger:  "default",
						Metadata: map[string]string{
							"foo": "bar",
						},
					}},
					Returns: []any{nil},
				},
			},
			Name: "nominal",
		},
	)
}
