package triggers

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilters(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		object     map[string]any
		filter     string
		shouldBeOk bool
	}

	testCases := []testCase{
		{
			name: "nominal",
			object: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": 3,
					},
				},
			},
			filter:     "event.a.b.c == 3",
			shouldBeOk: true,
		},
		{
			name: "comparison with $gt and float",
			object: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": math.Pi,
					},
				},
			},
			filter:     "event.a.b.c > 3",
			shouldBeOk: true,
		},
		{
			name: "comparison with $lt and float",
			object: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": 3.14,
					},
				},
			},
			filter: "event.a.b.c < 3.14",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ok, err := evalFilter(tc.object, tc.filter)
			require.NoError(t, err)
			require.Equal(t, tc.shouldBeOk, ok)
		})
	}
}

func TestEvalVariables(t *testing.T) {
	evaluated, err := evalVariables(map[string]any{
		"metadata": map[string]any{
			"psp": "stripe",
		},
	}, map[string]string{
		"psp": "event.metadata.psp",
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"psp": "stripe",
	}, evaluated)
}
