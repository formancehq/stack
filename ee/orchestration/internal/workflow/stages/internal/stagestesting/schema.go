package stagestesting

import (
	"reflect"
	"testing"

	"github.com/formancehq/orchestration/internal/schema"
	"github.com/stretchr/testify/require"
)

type SchemaTestCase struct {
	Data                    map[string]any
	ExpectedResolveError    bool
	ExpectedResolved        any
	ExpectedValidationError bool
	Name                    string
	Variables               map[string]string
}

func TestSchema(t *testing.T, stageName string, testCase SchemaTestCase) {
	t.Run(testCase.Name, func(t *testing.T) {
		t.Parallel()

		variables := testCase.Variables
		if variables == nil {
			variables = map[string]string{}
		}
		s, err := schema.Resolve(schema.Context{
			Variables: variables,
		}, testCase.Data, stageName)
		if !testCase.ExpectedResolveError {
			require.NoError(t, err, "resolving schema")
			require.Equal(t, testCase.ExpectedResolved, reflect.ValueOf(s).Elem().Interface())
		} else {
			require.Error(t, err, "resolving schema")
			return
		}

		err = schema.ValidateRequirements(s)
		if testCase.ExpectedValidationError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	})
}

func TestSchemas(t *testing.T, stageName string, testCases ...SchemaTestCase) {
	t.Parallel()
	for _, testCase := range testCases {
		TestSchema(t, stageName, testCase)
	}
}
