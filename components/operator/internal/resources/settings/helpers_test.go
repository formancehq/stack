package settings

import (
	"fmt"
	"testing"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/stretchr/testify/require"
)

func TestFindMatchingSettings(t *testing.T) {
	t.Parallel()
	type settings struct {
		key        string
		value      string
		isWildcard bool
	}
	type testCase struct {
		settings       []settings
		key            string
		expectedResult string
	}
	testCases := []testCase{
		{
			settings: []settings{
				{"postgres.ledger.dsn", "postgresql://localhost:5433", false},
				{"postgres.*.dsn", "postgresql://localhost:5432", false},
			},
			key:            "postgres.payments.dsn",
			expectedResult: "postgresql://localhost:5432",
		},
		{
			settings: []settings{
				{"postgres.*.dsn", "postgresql://localhost:5432", false},
				{"postgres.ledger.dsn", "postgresql://localhost:5433", false},
			},
			key:            "postgres.ledger.dsn",
			expectedResult: "postgresql://localhost:5433",
		},
		{
			settings: []settings{
				{"resource-requirements.*.containers.*.limits", "vvv", false},
				{"resource-requirements.ledger.containers.*.limits", "xxx", false},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "xxx",
		},
		{
			settings: []settings{
				{"resource-requirements.*.containers.*.limits", "vvv", false},
				{"resource-requirements.*.containers.ledger.limits", "xxx", false},
			},
			key:            "resource-requirements.payments.containers.payments.limits",
			expectedResult: "vvv",
		},
		{
			settings: []settings{
				{"resource-requirements.*.containers.ledger.limits", "xxx", false},
				{"resource-requirements.*.containers.*.limits", "vvv", false},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "xxx",
		},
		{
			settings: []settings{
				{"resource-requirements.*.containers.*.limits", "memory=512Mi", true},
				{"resource-requirements.*.containers.*.limits", "memory=1024Mi", false},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "memory=1024Mi",
		},
		{
			settings: []settings{
				{"resource-requirements.ledger.containers.ledger.limits", "memory=512Mi", true},
				{"resource-requirements.*.containers.*.limits", "memory=1024Mi", false},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "memory=1024Mi",
		},
		{
			settings: []settings{
				{"resource-requirements.ledger.containers.ledger.limits", "memory=512Mi", true},
				{"resource-requirements.*.containers.*.limits", "memory=1024Mi", false},
			},
			key:            "resource-requirements.payments.containers.payments.limits",
			expectedResult: "memory=1024Mi",
		},
		{
			settings: []settings{
				{"resource-requirements.*.containers.payments.limits", "memory=512Mi", true},
				{"resource-requirements.*.containers.*.limits", "memory=1024Mi", false},
			},
			key:            "resource-requirements.payments.containers.payments.limits",
			expectedResult: "memory=1024Mi",
		},
	}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			value, err := findMatchingSettings(Map(tc.settings, func(from settings) v1beta1.Settings {
				ret := v1beta1.Settings{
					Spec: v1beta1.SettingsSpec{
						Key:   from.key,
						Value: from.value,
					},
				}
				if from.isWildcard {
					ret.Spec.Stacks = []string{"*"}
				}
				return ret
			}), tc.key)
			require.NoError(t, err)
			require.NotNil(t, value)
			require.Equal(t, tc.expectedResult, *value)
		})
	}

}
