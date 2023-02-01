package delay

import (
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/schema"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestDelaySchema(t *testing.T) {
	type testCase struct {
		data          map[string]any
		expectedError bool
	}
	testCases := []testCase{
		{
			data: map[string]any{
				"until": time.Now().Format(time.RFC3339),
			},
			expectedError: false,
		},
		{
			data: map[string]any{
				"duration": "10s",
			},
			expectedError: false,
		},
		{
			data:          map[string]any{},
			expectedError: true,
		},
		{
			data: map[string]any{
				"until":    time.Now().Format(time.RFC3339),
				"duration": "10s",
			},
			expectedError: true,
		},
	}
	for _, testCase := range testCases {
		s, err := schema.Resolve(schema.Context{
			Variables: map[string]string{},
		}, testCase.data, "delay")
		require.NoError(t, err, "resolving schema")
		err = schema.ValidateRequirements(s)
		if testCase.expectedError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

type mockedActivity struct {
	activity any
	args     []any
	returns  []any
}

type testCase struct {
	delay          Delay
	mockedActivity []mockedActivity
	name           string
}

func ptr[T any](v T) *T {
	return &v
}

var testCases = []testCase{
	{
		delay: Delay{
			Until: ptr(time.Now().Add(time.Second)),
		},
		mockedActivity: nil,
		name:           "delay-until",
	},
	{
		delay: Delay{
			Duration: (*schema.Duration)(ptr(time.Second)),
		},
		mockedActivity: nil,
		name:           "delay-duration",
	},
}

func TestDelay(t *testing.T) {
	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}

			env := testSuite.NewTestWorkflowEnvironment()
			for _, ma := range tc.mockedActivity {
				env.OnActivity(ma.activity, ma.args...).Return(ma.returns...)
			}

			env.ExecuteWorkflow(RunDelay, tc.delay)
			require.True(t, env.IsWorkflowCompleted())
			require.NoError(t, env.GetWorkflowError())
		})

	}
}
