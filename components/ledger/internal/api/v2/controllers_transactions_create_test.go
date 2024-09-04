package v2

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/pointer"
	ledger "github.com/formancehq/ledger/internal"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransactionCreate(t *testing.T) {
	type testCase struct {
		name                 string
		expectedDryRun       bool
		expectedRunScript    ledger.RunScript
		returnError          error
		payload              any
		expectedStatusCode   int
		expectedErrorCode    string
		queryParams          url.Values
		expectControllerCall bool
	}

	testCases := []testCase{
		{
			name: "using plain numscript",
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `XXX`,
					},
				},
			},
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `XXX`,
					Vars:  map[string]string{},
				},
			},
			expectControllerCall: true,
		},
		{
			name: "using plain numscript with variables",
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `vars {
						monetary $val
					}

					send $val (
						source = @world
						destination = @bank
					)`,
					},
					Vars: map[string]any{
						"val": "USD/2 100",
					},
				},
			},
			expectControllerCall: true,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `vars {
						monetary $val
					}

					send $val (
						source = @world
						destination = @bank
					)`,
					Vars: map[string]string{
						"val": "USD/2 100",
					},
				},
			},
		},
		{
			name:                 "using plain numscript with variables (legacy format)",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `vars {
						monetary $val
					}

					send $val (
						source = @world
						destination = @bank
					)`,
					},
					Vars: map[string]any{
						"val": map[string]any{
							"asset":  "USD/2",
							"amount": 100,
						},
					},
				},
			},
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `vars {
						monetary $val
					}

					send $val (
						source = @world
						destination = @bank
					)`,
					Vars: map[string]string{
						"val": "USD/2 100",
					},
				},
			},
		},
		{
			name:                 "using plain numscript and dry run",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `send (
						source = @world
						destination = @bank
					)`,
					},
				},
			},
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `send (
						source = @world
						destination = @bank
					)`,
					Vars: map[string]string{},
				},
			},
			expectedDryRun: true,
			queryParams: url.Values{
				"dryRun": []string{"true"},
			},
		},
		{
			name:                 "using JSON postings",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Postings: []ledger.Posting{
					ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
				},
			},
			expectedRunScript: ledger.TxToScriptData(ledger.NewTransactionData().WithPostings(
				ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
			), false),
		},
		{
			name:                 "using JSON postings and dry run",
			expectControllerCall: true,
			queryParams: url.Values{
				"dryRun": []string{"true"},
			},
			payload: ledger.TransactionRequest{
				Postings: []ledger.Posting{
					ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
				},
			},
			expectedDryRun: true,
			expectedRunScript: ledger.TxToScriptData(ledger.NewTransactionData().WithPostings(
				ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
			), false),
		},
		{
			name: "no postings or script",
			payload: ledger.TransactionRequest{
				Metadata: map[string]string{},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrNoPostings,
			returnError:        errors.New("you need to pass either a posting array or a numscript script"),
		},
		{
			name: "postings and script",
			payload: ledger.TransactionRequest{
				Postings: ledger.Postings{
					{
						Source:      "world",
						Destination: "alice",
						Amount:      big.NewInt(100),
						Asset:       "COIN",
					},
				},
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `
						send [COIN 100] (
						  source = @world
						  destination = @bob
						)`,
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "using invalid body",
			payload:            "not a valid payload",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:                 "with insufficient funds",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `XXX`,
					},
				},
			},
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `XXX`,
					Vars:  map[string]string{},
				},
			},
			returnError:        &ledgercontroller.ErrInsufficientFunds{},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInsufficientFund,
		},
		{
			name: "using JSON postings and negative amount",
			payload: ledger.TransactionRequest{
				Postings: []ledger.Posting{
					ledger.NewPosting("world", "bank", "USD", big.NewInt(-100)),
				},
			},
			expectControllerCall: true,
			expectedStatusCode:   http.StatusBadRequest,
			expectedErrorCode:    ErrCompilationFailed,
			expectedRunScript: ledger.TxToScriptData(ledger.NewTransactionData().WithPostings(
				ledger.NewPosting("world", "bank", "USD", big.NewInt(-100)),
			), false),
			returnError: &ledgercontroller.ErrInvalidVars{},
		},
		{
			expectControllerCall: true,
			name:                 "numscript and negative amount",
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `send [COIN -100] (
						source = @world
						destination = @bob
					)`,
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrCompilationFailed,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `send [COIN -100] (
						source = @world
						destination = @bob
					)`,
					Vars: map[string]string{},
				},
			},
			returnError: &ledgercontroller.ErrInvalidVars{},
		},
		{
			name:                 "numscript and compilation failed",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `send [COIN XXX] (
						source = @world
						destination = @bob
					)`,
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrCompilationFailed,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `send [COIN XXX] (
						source = @world
						destination = @bob
					)`,
					Vars: map[string]string{},
				},
			},
			returnError: ledgercontroller.ErrCompilationFailed{},
		},
		{
			name:                 "numscript and no postings",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `vars {}`,
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrNoPostings,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `vars {}`,
					Vars:  map[string]string{},
				},
			},
			returnError: ledgercontroller.ErrNoPostings,
		},
		{
			name:                 "numscript and metadata override",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `send [COIN 100] (
						source = @world
						destination = @bob
					)
					set_tx_meta("foo", "bar")`,
					},
				},
				Reference: "xxx",
				Metadata: map[string]string{
					"foo": "baz",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMetadataOverride,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `send [COIN 100] (
						source = @world
						destination = @bob
					)
					set_tx_meta("foo", "bar")`,
					Vars: map[string]string{},
				},
				Reference: "xxx",
				Metadata: map[string]string{
					"foo": "baz",
				},
			},
			returnError: &ledgercontroller.ErrMetadataOverride{},
		},
		{
			name:                 "unexpected error",
			expectControllerCall: true,
			payload: ledger.TransactionRequest{
				Script: ledger.ScriptV1{
					Script: ledger.Script{
						Plain: `send [COIN 100] (
						source = @world
						destination = @bob
					)`,
					},
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  api.ErrorInternal,
			expectedRunScript: ledger.RunScript{
				Script: ledger.Script{
					Plain: `send [COIN 100] (
						source = @world
						destination = @bob
					)`,
					Vars: map[string]string{},
				},
			},
			returnError: errors.New("unexpected error"),
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			expectedTx := ledger.NewTransaction().WithPostings(
				ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
			)

			systemController, ledgerController := newTestingSystemController(t, true)
			if testCase.expectControllerCall {
				expect := ledgerController.EXPECT().
					CreateTransaction(gomock.Any(), ledgercontroller.Parameters{
						DryRun: tc.expectedDryRun,
					}, testCase.expectedRunScript)

				if tc.returnError == nil {
					expect.Return(pointer.For(expectedTx), nil)
				} else {
					expect.Return(nil, tc.returnError)
				}
			}

			router := NewRouter(systemController, auth.NewNoAuth(), "develop", testing.Verbose())

			req := httptest.NewRequest(http.MethodPost, "/xxx/transactions", api.Buffer(t, testCase.payload))
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				tx, ok := api.DecodeSingleResponse[ledger.Transaction](t, rec.Body)
				require.True(t, ok)
				require.Equal(t, expectedTx, tx)
			} else {
				err := api.ErrorResponse{}
				api.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
