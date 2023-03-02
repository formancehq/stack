package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/formancehq/ledger/pkg/api/apierrors"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type TransactionController struct{}

func NewTransactionController() TransactionController {
	return TransactionController{}
}

func (ctl *TransactionController) CountTransactions(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	var startTimeParsed, endTimeParsed time.Time
	var err error
	if r.URL.Query().Get(QueryKeyStartTime) != "" {
		startTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyStartTime))
		if err != nil {
			apierrors.ResponseError(w, r, ErrInvalidStartTime)
			return
		}
	}
	if r.URL.Query().Get(QueryKeyStartTimeDeprecated) != "" {
		startTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyStartTimeDeprecated))
		if err != nil {
			apierrors.ResponseError(w, r, ErrInvalidStartTimeDeprecated)
			return
		}
	}

	if r.URL.Query().Get(QueryKeyEndTime) != "" {
		endTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyEndTime))
		if err != nil {
			apierrors.ResponseError(w, r, ErrInvalidEndTime)
			return
		}
	}
	if r.URL.Query().Get(QueryKeyEndTimeDeprecated) != "" {
		endTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyEndTimeDeprecated))
		if err != nil {
			apierrors.ResponseError(w, r, ErrInvalidEndTimeDeprecated)
			return
		}
	}

	txQuery := ledger.NewTransactionsQuery().
		WithReferenceFilter(r.URL.Query().Get("reference")).
		WithAccountFilter(r.URL.Query().Get("account")).
		WithSourceFilter(r.URL.Query().Get("source")).
		WithDestinationFilter(r.URL.Query().Get("destination")).
		WithStartTimeFilter(startTimeParsed).
		WithEndTimeFilter(endTimeParsed)

	count, err := l.CountTransactions(r.Context(), *txQuery)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	w.Header().Set("Count", fmt.Sprint(count))
}

func (ctl *TransactionController) GetTransactions(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	txQuery := ledger.NewTransactionsQuery()

	if r.URL.Query().Get(QueryKeyCursor) != "" {
		if r.URL.Query().Get("after") != "" ||
			r.URL.Query().Get("reference") != "" ||
			r.URL.Query().Get("account") != "" ||
			r.URL.Query().Get("source") != "" ||
			r.URL.Query().Get("destination") != "" ||
			r.URL.Query().Get(QueryKeyStartTime) != "" ||
			r.URL.Query().Get(QueryKeyStartTimeDeprecated) != "" ||
			r.URL.Query().Get(QueryKeyEndTime) != "" ||
			r.URL.Query().Get(QueryKeyEndTimeDeprecated) != "" ||
			r.URL.Query().Get(QueryKeyPageSize) != "" ||
			r.URL.Query().Get(QueryKeyPageSizeDeprecated) != "" {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("no other query params can be set with '%s'", QueryKeyCursor)))
			return
		}

		res, err := base64.RawURLEncoding.DecodeString(r.URL.Query().Get(QueryKeyCursor))
		if err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("invalid '%s' query param", QueryKeyCursor)))
			return
		}

		token := sqlstorage.TxsPaginationToken{}
		if err = json.Unmarshal(res, &token); err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("invalid '%s' query param", QueryKeyCursor)))
			return
		}

		txQuery = txQuery.
			WithAfterTxID(token.AfterTxID).
			WithReferenceFilter(token.ReferenceFilter).
			WithAccountFilter(token.AccountFilter).
			WithSourceFilter(token.SourceFilter).
			WithDestinationFilter(token.DestinationFilter).
			WithStartTimeFilter(token.StartTime).
			WithEndTimeFilter(token.EndTime).
			WithMetadataFilter(token.MetadataFilter).
			WithPageSize(token.PageSize)

	} else if r.URL.Query().Get(QueryKeyCursorDeprecated) != "" {
		if r.URL.Query().Get("after") != "" ||
			r.URL.Query().Get("reference") != "" ||
			r.URL.Query().Get("account") != "" ||
			r.URL.Query().Get("source") != "" ||
			r.URL.Query().Get("destination") != "" ||
			r.URL.Query().Get(QueryKeyStartTime) != "" ||
			r.URL.Query().Get(QueryKeyStartTimeDeprecated) != "" ||
			r.URL.Query().Get(QueryKeyEndTime) != "" ||
			r.URL.Query().Get(QueryKeyEndTimeDeprecated) != "" ||
			r.URL.Query().Get(QueryKeyPageSize) != "" ||
			r.URL.Query().Get(QueryKeyPageSizeDeprecated) != "" {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("no other query params can be set with '%s'", QueryKeyCursorDeprecated)))
			return
		}

		res, err := base64.RawURLEncoding.DecodeString(r.URL.Query().Get(QueryKeyCursorDeprecated))
		if err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("invalid '%s' query param", QueryKeyCursorDeprecated)))
			return
		}

		token := sqlstorage.TxsPaginationToken{}
		if err = json.Unmarshal(res, &token); err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(
				fmt.Sprintf("invalid '%s' query param", QueryKeyCursorDeprecated)))
			return
		}

		txQuery = txQuery.
			WithAfterTxID(token.AfterTxID).
			WithReferenceFilter(token.ReferenceFilter).
			WithAccountFilter(token.AccountFilter).
			WithSourceFilter(token.SourceFilter).
			WithDestinationFilter(token.DestinationFilter).
			WithStartTimeFilter(token.StartTime).
			WithEndTimeFilter(token.EndTime).
			WithMetadataFilter(token.MetadataFilter).
			WithPageSize(token.PageSize)

	} else {
		var err error
		var afterTxIDParsed uint64
		if r.URL.Query().Get("after") != "" {
			afterTxIDParsed, err = strconv.ParseUint(r.URL.Query().Get("after"), 10, 64)
			if err != nil {
				apierrors.ResponseError(w, r, ledger.NewValidationError(
					"invalid 'after' query param"))
				return
			}
		}

		var startTimeParsed, endTimeParsed time.Time
		if r.URL.Query().Get(QueryKeyStartTime) != "" {
			startTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyStartTime))
			if err != nil {
				apierrors.ResponseError(w, r, ErrInvalidStartTime)
				return
			}
		}
		if r.URL.Query().Get(QueryKeyStartTimeDeprecated) != "" {
			startTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyStartTimeDeprecated))
			if err != nil {
				apierrors.ResponseError(w, r, ErrInvalidStartTimeDeprecated)
				return
			}
		}

		if r.URL.Query().Get(QueryKeyEndTime) != "" {
			endTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyEndTime))
			if err != nil {
				apierrors.ResponseError(w, r, ErrInvalidEndTime)
				return
			}
		}
		if r.URL.Query().Get(QueryKeyEndTimeDeprecated) != "" {
			endTimeParsed, err = time.Parse(time.RFC3339, r.URL.Query().Get(QueryKeyEndTimeDeprecated))
			if err != nil {
				apierrors.ResponseError(w, r, ErrInvalidEndTimeDeprecated)
				return
			}
		}

		pageSize, err := getPageSize(w, r)
		if err != nil {
			apierrors.ResponseError(w, r, err)
			return
		}

		txQuery = txQuery.
			WithAfterTxID(afterTxIDParsed).
			WithReferenceFilter(r.URL.Query().Get("reference")).
			WithAccountFilter(r.URL.Query().Get("account")).
			WithSourceFilter(r.URL.Query().Get("source")).
			WithDestinationFilter(r.URL.Query().Get("destination")).
			WithStartTimeFilter(startTimeParsed).
			WithEndTimeFilter(endTimeParsed).
			WithMetadataFilter(sharedapi.GetQueryMap(r.URL.Query(), "metadata")).
			WithPageSize(pageSize)
	}

	cursor, err := l.GetTransactions(r.Context(), *txQuery)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.RenderCursor(w, cursor)
}

type PostTransaction struct {
	Postings  core.Postings `json:"postings"`
	Script    core.Script   `json:"script"`
	Timestamp time.Time     `json:"timestamp"`
	Reference string        `json:"reference"`
	Metadata  core.Metadata `json:"metadata" swaggertype:"object"`
}

func (ctl *TransactionController) PostTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	value := r.URL.Query().Get("preview")
	preview := strings.ToUpper(value) == "YES" || strings.ToUpper(value) == "TRUE" || value == "1"

	payload := PostTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		apierrors.ResponseError(w, r,
			ledger.NewValidationError("invalid transaction format"))
		return
	}

	if len(payload.Postings) > 0 && payload.Script.Plain != "" ||
		len(payload.Postings) == 0 && payload.Script.Plain == "" {
		apierrors.ResponseError(w, r, ledger.NewValidationError(
			"invalid payload: should contain either postings or script"))
		return
	} else if len(payload.Postings) > 0 {
		if i, err := payload.Postings.Validate(); err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(errors.Wrap(err,
				fmt.Sprintf("invalid posting %d", i)).Error()))
			return
		}
		txData := core.TransactionData{
			Postings:  payload.Postings,
			Timestamp: payload.Timestamp,
			Reference: payload.Reference,
			Metadata:  payload.Metadata,
		}
		res, err := l.ExecuteTxsData(r.Context(), preview, txData)
		if err != nil {
			apierrors.ResponseError(w, r, err)
			return
		}

		sharedapi.Ok(w, res)
		return
	}

	script := core.ScriptData{
		Script:    payload.Script,
		Timestamp: payload.Timestamp,
		Reference: payload.Reference,
		Metadata:  payload.Metadata,
	}
	res, err := l.ExecuteScript(r.Context(), preview, script)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, []core.ExpandedTransaction{res})
}

func (ctl *TransactionController) GetTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	txId, err := strconv.ParseUint(chi.URLParam(r, "txid"), 10, 64)
	if err != nil {
		apierrors.ResponseError(w, r, ledger.NewValidationError("invalid transaction ID"))
		return
	}

	tx, err := l.GetTransaction(r.Context(), txId)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, tx)
}

func (ctl *TransactionController) RevertTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	txId, err := strconv.ParseUint(chi.URLParam(r, "txid"), 10, 64)
	if err != nil {
		apierrors.ResponseError(w, r, ledger.NewValidationError("invalid transaction ID"))
		return
	}

	tx, err := l.RevertTransaction(r.Context(), txId)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, tx)
}

func (ctl *TransactionController) PostTransactionMetadata(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	var m core.Metadata
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		apierrors.ResponseError(w, r, ledger.NewValidationError("invalid metadata format"))
		return
	}

	txId, err := strconv.ParseUint(chi.URLParam(r, "txid"), 10, 64)
	if err != nil {
		apierrors.ResponseError(w, r, ledger.NewValidationError("invalid transaction ID"))
		return
	}

	_, err = l.GetTransaction(r.Context(), txId)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	if err := l.SaveMeta(r.Context(),
		core.MetaTargetTypeTransaction, txId, m); err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.NoContent(w)
}

func (ctl *TransactionController) PostTransactionsBatch(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	var txs core.Transactions
	if err := json.NewDecoder(r.Body).Decode(&txs); err != nil {
		apierrors.ResponseError(w, r, ledger.NewValidationError("invalid transactions format"))
		return
	}

	if len(txs.Transactions) == 0 {
		apierrors.ResponseError(w, r, ledger.NewValidationError("no transaction to insert"))
		return
	}

	for i, tx := range txs.Transactions {
		if len(tx.Postings) == 0 {
			apierrors.ResponseError(w, r, ledger.NewValidationError(errors.New(fmt.Sprintf(
				"invalid transaction %d: no postings", i)).Error()))
			return
		}
		if j, err := tx.Postings.Validate(); err != nil {
			apierrors.ResponseError(w, r, ledger.NewValidationError(errors.Wrap(err,
				fmt.Sprintf("invalid transaction %d: posting %d", i, j)).Error()))
			return
		}
	}

	res, err := l.ExecuteTxsData(r.Context(), false, txs.Transactions...)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, res)
}
