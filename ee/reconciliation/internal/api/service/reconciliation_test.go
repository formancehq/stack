package service

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/reconciliation/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestReconciliation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name             string
		ledgerVersion    string
		paymentsVersion  string
		ledgerBalances   map[string]*big.Int
		paymentsBalances map[string]*big.Int
		expectedReco     *models.Reconciliation
		expectedError    bool
	}

	testCases := []testCase{
		{
			name:            "nominal with drift = 0",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(100),
				"EUR": big.NewInt(200),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"EUR": big.NewInt(-200),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(200),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(-200),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(0),
				},
				Error: "",
			},
		},
		{
			name:            "nominal with drift >= 0",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(200),
				"EUR": big.NewInt(300),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"EUR": big.NewInt(-200),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(200),
					"EUR": big.NewInt(300),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(-200),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(100),
				},
				Error: "",
			},
		},
		{
			name:            "nominal with drift < 0",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(100),
				"EUR": big.NewInt(200),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"EUR": big.NewInt(-400),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationNotOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(200),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(-400),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(200),
				},
				Error: "balance drift for asset EUR",
			},
		},
		{
			name:            "different length, no drift",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"EUR": big.NewInt(200),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"EUR": big.NewInt(-200),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationNotOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(200),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(-200),
				},
				DriftBalances: map[string]*big.Int{
					"EUR": big.NewInt(0),
					"USD": big.NewInt(100),
				},
				Error: "balance drift for asset USD",
			},
		},
		{
			name:            "same length, different asset and no drift",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(100),
				"EUR": big.NewInt(200),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"DKK": big.NewInt(-200),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationNotOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(200),
					"DKK": big.NewInt(0),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"DKK": big.NewInt(-200),
					"EUR": big.NewInt(0),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(200),
					"DKK": big.NewInt(200),
				},
				Error: "balance drift for asset DKK",
			},
		},
		{
			name:            "missing payments balance with ledger balance at 0",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(100),
				"EUR": big.NewInt(0),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(0),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(0),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(0),
				},
				Error: "",
			},
		},
		{
			name:            "missing ledger balance with payments balance at 0",
			ledgerVersion:   "v2.0.0-beta.1",
			paymentsVersion: "v1.0.0-rc.4",
			ledgerBalances: map[string]*big.Int{
				"USD": big.NewInt(100),
			},
			paymentsBalances: map[string]*big.Int{
				"USD": big.NewInt(-100),
				"EUR": big.NewInt(0),
			},
			expectedReco: &models.Reconciliation{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Time{},
				Status:               models.ReconciliationOK,
				LedgerBalances: map[string]*big.Int{
					"USD": big.NewInt(100),
					"EUR": big.NewInt(0),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD": big.NewInt(-100),
					"EUR": big.NewInt(0),
				},
				DriftBalances: map[string]*big.Int{
					"USD": big.NewInt(0),
					"EUR": big.NewInt(0),
				},
				Error: "",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(newMockStore(), newMockSDKFormanceClient(
				tc.ledgerVersion,
				tc.ledgerBalances,
				tc.paymentsVersion,
				tc.paymentsBalances,
			))

			reco, err := s.Reconciliation(context.Background(), uuid.New().String(), &ReconciliationRequest{})
			if tc.expectedError {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedReco.Status, reco.Status)
			compareBalancesMap(t, tc.expectedReco.LedgerBalances, reco.LedgerBalances)
			compareBalancesMap(t, tc.expectedReco.PaymentsBalances, reco.PaymentsBalances)
			compareBalancesMap(t, tc.expectedReco.DriftBalances, reco.DriftBalances)
			require.Equal(t, tc.expectedReco.Error, reco.Error)
		})
	}
}

func compareBalancesMap(t *testing.T, expected, actual map[string]*big.Int) {
	require.Equal(t, len(expected), len(actual))
	for k, v := range expected {
		require.Equal(t, v.Cmp(actual[k]), 0)
	}
}
