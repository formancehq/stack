package triggers

import (
	"math"
	"net/http"
	"net/http/httptest"
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

			e := NewExpressionEvaluator(http.DefaultClient)

			ok, err := e.evalFilter(tc.object, tc.filter)
			require.NoError(t, err)
			require.Equal(t, tc.shouldBeOk, ok)
		})
	}
}

func TestEvalVariables(t *testing.T) {

	type testCase struct {
		name           string
		rawObject      any
		variables      map[string]string
		expectedResult map[string]string
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data": {"role": "admin"}}`))
	}))
	t.Cleanup(srv.Close)

	for _, testCase := range []testCase{
		{
			name: "nominal",
			rawObject: map[string]any{
				"metadata": map[string]any{
					"psp": "stripe",
				},
			},
			variables: map[string]string{
				"psp": "event.metadata.psp",
			},
			expectedResult: map[string]string{
				"psp": "stripe",
			},
		},
		{
			name: "using links",
			rawObject: map[string]any{
				"links": []map[string]any{
					{
						"name": "source_account",
						"uri":  srv.URL,
					},
				},
			},
			variables: map[string]string{
				"role": `link(event, "source_account").role`,
			},
			expectedResult: map[string]string{
				"role": "admin",
			},
		},
		{
			name: "using float",
			rawObject: map[string]any{
				"amount": float64(1200000),
			},
			variables: map[string]string{
				"amount": "event.amount",
			},
			expectedResult: map[string]string{
				"amount": "1200000",
			},
		},
		{
			name: "using bool",
			rawObject: map[string]any{
				"test": true,
			},
			variables: map[string]string{
				"test": "event.test",
			},
			expectedResult: map[string]string{
				"test": "true",
			},
		},
		{
			name: "using integer",
			rawObject: map[string]any{
				"amount": int64(1200000),
			},
			variables: map[string]string{
				"amount": "event.amount",
			},
			expectedResult: map[string]string{
				"amount": "1200000",
			},
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			e := NewExpressionEvaluator(http.DefaultClient)
			evaluated, err := e.evalVariables(testCase.rawObject, testCase.variables)
			require.NoError(t, err)
			require.Equal(t, testCase.expectedResult, evaluated)
		})
	}
}
