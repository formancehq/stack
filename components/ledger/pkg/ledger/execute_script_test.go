package ledger_test

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/ledger/pkg/cache"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/ledger/pkg/machine"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name             string
	setup            func(t *testing.T, l *ledger.Ledger)
	script           string
	reference        string
	expectedError    error
	expectedTx       core.ExpandedTransaction
	expectedLogs     []core.Log
	expectedAccounts map[string]core.AccountWithVolumes
}

var testCases = []testCase{
	{
		name: "nominal",
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		expectedTx: core.ExpandedTransaction{
			Transaction: core.NewTransaction().WithPostings(
				core.NewPosting("world", "mint", "GEM", core.NewMonetaryInt(100)),
			),
			PreCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"GEM": core.NewEmptyVolumes(),
				},
				"mint": {
					"GEM": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"GEM": core.NewEmptyVolumes().WithOutput(core.NewMonetaryInt(100)),
				},
				"mint": {
					"GEM": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(100)),
				},
			},
		},
		expectedLogs: []core.Log{
			core.NewTransactionLog(
				core.NewTransaction().WithPostings(
					core.NewPosting("world", "mint", "GEM", core.NewMonetaryInt(100))),
				map[string]core.Metadata{},
			),
		},
		expectedAccounts: map[string]core.AccountWithVolumes{
			"mint": {
				Account: core.NewAccount("mint"),
				Volumes: core.AssetsVolumes{
					"GEM": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(100)),
				},
				Balances: map[string]*core.MonetaryInt{
					"GEM": core.NewMonetaryInt(100),
				},
			},
		},
	},
	{
		name:          "no script",
		script:        ``,
		expectedError: machine.NewScriptError(machine.ScriptErrorNoScript, ""),
	},
	{
		name:          "invalid script",
		script:        `XXX`,
		expectedError: machine.NewScriptError(machine.ScriptErrorCompilationFailed, ""),
	},
	{
		name: "set reference conflict",
		setup: func(t *testing.T, l *ledger.Ledger) {
			require.NoError(t, l.GetLedgerStore().InsertTransactions(context.Background(), core.ExpandedTransaction{
				Transaction: core.NewTransaction().
					WithPostings(core.NewPosting("world", "mint", "GEM", core.NewMonetaryInt(100))).
					WithReference("tx_ref"),
			}))
		},
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		reference:     "tx_ref",
		expectedError: ledger.NewConflictError(),
	},
	{
		name: "set reference",
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		reference: "tx_ref",
		expectedTx: core.ExpandedTransaction{
			Transaction: core.NewTransaction().
				WithPostings(
					core.NewPosting("world", "mint", "GEM", core.NewMonetaryInt(100)),
				).
				WithReference("tx_ref"),
			PreCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"GEM": core.NewEmptyVolumes(),
				},
				"mint": {
					"GEM": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"GEM": core.NewEmptyVolumes().WithOutput(core.NewMonetaryInt(100)),
				},
				"mint": {
					"GEM": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(100)),
				},
			},
		},
		expectedLogs: []core.Log{
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(
						core.NewPosting("world", "mint", "GEM", core.NewMonetaryInt(100)),
					).
					WithReference("tx_ref"),
				map[string]core.Metadata{},
			),
		},
		expectedAccounts: map[string]core.AccountWithVolumes{
			"mint": {
				Account: core.NewAccount("mint"),
				Volumes: core.AssetsVolumes{
					"GEM": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(100)),
				},
				Balances: map[string]*core.MonetaryInt{
					"GEM": core.NewMonetaryInt(100),
				},
			},
		},
	},
}

func TestExecuteScript(t *testing.T) {
	t.Parallel()
	now := core.Now()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runOnLedger(t, func(l *ledger.Ledger) {

				if tc.setup != nil {
					tc.setup(t, l)
				}

				ret, _, err := l.CreateTransactionAndWait(context.Background(), false, core.ScriptData{
					Script: core.Script{
						Plain: tc.script,
					},
					Timestamp: now,
					Reference: tc.reference,
				})
				if tc.expectedError != nil {
					require.True(t, errors.Is(err, tc.expectedError))
				} else {
					require.NoError(t, err)
					tc.expectedTx.Timestamp = now
					require.Equal(t, tc.expectedTx, *ret)

					logs, err := l.GetLedgerStore().ReadLogsStartingFromID(context.Background(), 0)
					require.NoError(t, err)
					require.Len(t, logs, len(tc.expectedLogs))
					for ind := range tc.expectedLogs {
						var previous *core.Log
						if ind > 0 {
							previous = &tc.expectedLogs[ind-1]
						}
						expectedLog := tc.expectedLogs[ind]
						switch v := expectedLog.Data.(type) {
						case core.NewTransactionLogPayload:
							v.Transaction.Timestamp = now
							expectedLog.Data = v
						}
						expectedLog.Date = now
						require.Equal(t, expectedLog.ComputeHash(previous), logs[ind])
					}

					lastTXInfo, err := l.GetDBCache().GetLastTransaction(context.Background())
					require.NoError(t, err)
					require.NotNil(t, lastTXInfo)
					require.Equal(t, cache.TxInfo{
						Date: tc.expectedTx.Timestamp,
					}, *lastTXInfo)

					<-time.After(2 * time.Second)
					for address, account := range tc.expectedAccounts {
						accountFromCache, err := l.GetDBCache().GetAccountWithVolumes(context.Background(), address)
						require.NoError(t, err)
						require.NotNil(t, accountFromCache)
						require.Equal(t, account, *accountFromCache)
					}
				}
			})
		})
	}
}
