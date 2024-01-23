package settings

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindMatchingSettings(t *testing.T) {
	t.Parallel()
	type settings struct {
		key   string
		value string
	}
	type testCase struct {
		name           string
		settings       []settings
		key            string
		expectedResult string
	}
	testCases := []testCase{
		{
			name: "priority using wildcard",
			settings: []settings{
				{"postgres.ledger.dsn", "postgresql://localhost:5433"},
				{"postgres.*.dsn", "postgresql://localhost:5432"},
			},
			key:            "postgres.payments.dsn",
			expectedResult: "postgresql://localhost:5432",
		},
		{
			name: "priority using specific dsn",
			settings: []settings{
				{"postgres.*.dsn", "postgresql://localhost:5432"},
				{"postgres.ledger.dsn", "postgresql://localhost:5433"},
			},
			key:            "postgres.ledger.dsn",
			expectedResult: "postgresql://localhost:5433",
		},
		{
			name: "priority using specific dsn",
			settings: []settings{
				{"resource-requirements.*.containers.*.limits", "vvv"},
				{"resource-requirements.ledger.containers.*.limits", "xxx"},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "xxx",
		},
		{
			name: "priority using specific dsn and multiple wildcard",
			settings: []settings{
				{"resource-requirements.*.containers.*.limits", "vvv"},
				{"resource-requirements.*.containers.ledger.limits", "xxx"},
			},
			key:            "resource-requirements.payments.containers.payments.limits",
			expectedResult: "vvv",
		},
		{
			name: "priority using specific dsn and multiple wildcard",
			settings: []settings{
				{"resource-requirements.*.containers.ledger.limits", "xxx"},
				{"resource-requirements.*.containers.*.limits", "vvv"},
			},
			key:            "resource-requirements.ledger.containers.ledger.limits",
			expectedResult: "xxx",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			value, err := findMatchingSettings(Map(tc.settings, func(from settings) v1beta1.Settings {
				return v1beta1.Settings{
					Spec: v1beta1.SettingsSpec{
						Key:   from.key,
						Value: from.value,
					},
				}
			}), tc.key)
			require.NoError(t, err)
			require.NotNil(t, value)
			require.Equal(t, tc.expectedResult, *value)
		})
	}

}
