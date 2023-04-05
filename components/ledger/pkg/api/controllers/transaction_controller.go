package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/formancehq/ledger/pkg/api/apierrors"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/ledger/pkg/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/errorsutil"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func CountTransactions(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	var startTimeParsed, endTimeParsed core.Time
	var err error
	if r.URL.Query().Get(QueryKeyStartTime) != "" {
		startTimeParsed, err = core.ParseTime(r.URL.Query().Get(QueryKeyStartTime))
		if err != nil {
			apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation, ErrInvalidStartTime))
			return
		}
	}

	if r.URL.Query().Get(QueryKeyEndTime) != "" {
		endTimeParsed, err = core.ParseTime(r.URL.Query().Get(QueryKeyEndTime))
		if err != nil {
			apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation, ErrInvalidEndTime))
			return
		}
	}

	txQuery := storage.NewTransactionsQuery().
		WithReferenceFilter(r.URL.Query().Get("reference")).
		WithAccountFilter(r.URL.Query().Get("account")).
		WithSourceFilter(r.URL.Query().Get("source")).
		WithDestinationFilter(r.URL.Query().Get("destination")).
		WithStartTimeFilter(startTimeParsed).
		WithEndTimeFilter(endTimeParsed).
		WithMetadataFilter(sharedapi.GetQueryMap(r.URL.Query(), "metadata"))

	count, err := l.CountTransactions(r.Context(), txQuery)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	w.Header().Set("Count", fmt.Sprint(count))
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	txQuery := storage.NewTransactionsQuery()

	if r.URL.Query().Get(QueryKeyCursor) != "" {
		if r.URL.Query().Get("after") != "" ||
			r.URL.Query().Get("reference") != "" ||
			r.URL.Query().Get("account") != "" ||
			r.URL.Query().Get("source") != "" ||
			r.URL.Query().Get("destination") != "" ||
			r.URL.Query().Get(QueryKeyStartTime) != "" ||
			r.URL.Query().Get(QueryKeyEndTime) != "" ||
			r.URL.Query().Get(QueryKeyPageSize) != "" {
			apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation,
				errors.Errorf("no other query params can be set with '%s'", QueryKeyCursor)))
			return
		}

		err := storage.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &txQuery)
		if err != nil {
			apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation,
				errors.Errorf("invalid '%s' query param", QueryKeyCursor)))
			return
		}
	} else {
		var (
			err             error
			afterTxIDParsed uint64
		)
		if r.URL.Query().Get("after") != "" {
			afterTxIDParsed, err = strconv.ParseUint(r.URL.Query().Get("after"), 10, 64)
			if err != nil {
				apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation,
					errors.New("invalid 'after' query param")))
				return
			}
		}

		var startTimeParsed, endTimeParsed core.Time
		if r.URL.Query().Get(QueryKeyStartTime) != "" {
			startTimeParsed, err = core.ParseTime(r.URL.Query().Get(QueryKeyStartTime))
			if err != nil {
				apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation, ErrInvalidStartTime))
				return
			}
		}

		if r.URL.Query().Get(QueryKeyEndTime) != "" {
			endTimeParsed, err = core.ParseTime(r.URL.Query().Get(QueryKeyEndTime))
			if err != nil {
				apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation, ErrInvalidEndTime))
				return
			}
		}

		pageSize, err := getPageSize(r)
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

	cursor, err := l.GetTransactions(r.Context(), txQuery)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.RenderCursor(w, *cursor)
}

type PostTransactionRequest struct {
	Postings  core.Postings     `json:"postings"`
	Script    core.Script       `json:"script"`
	Timestamp core.Time         `json:"timestamp"`
	Reference string            `json:"reference"`
	Metadata  metadata.Metadata `json:"metadata" swaggertype:"object"`
}

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	value := r.URL.Query().Get("preview")
	preview := strings.ToUpper(value) == "YES" || strings.ToUpper(value) == "TRUE" || value == "1"

	payload := PostTransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		apierrors.ResponseError(w, r,
			errorsutil.NewError(ledger.ErrValidation,
				errors.New("invalid transaction format")))
		return
	}

	if len(payload.Postings) > 0 && payload.Script.Plain != "" ||
		len(payload.Postings) == 0 && payload.Script.Plain == "" {
		apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation,
			errors.New("invalid payload: should contain either postings or script")))
		return
	} else if len(payload.Postings) > 0 {
		if i, err := payload.Postings.Validate(); err != nil {
			apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation, errors.Wrap(err,
				fmt.Sprintf("invalid posting %d", i))))
			return
		}
		txData := core.TransactionData{
			Postings:  payload.Postings,
			Timestamp: payload.Timestamp,
			Reference: payload.Reference,
			Metadata:  payload.Metadata,
		}

		res, err := l.CreateTransaction(r.Context(), preview, core.TxToScriptData(txData))
		if err != nil {
			apierrors.ResponseError(w, r, err)
			return
		}

		sharedapi.Ok(w, res)
		return
	}

	script := core.RunScript{
		Script:    payload.Script,
		Timestamp: payload.Timestamp,
		Reference: payload.Reference,
		Metadata:  payload.Metadata,
	}

	res, err := l.CreateTransaction(r.Context(), preview, script)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, res)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	tx, err := l.GetTransaction(r.Context(), chi.URLParam(r, "txid"))
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, tx)
}

func RevertTransaction(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	tx, err := l.RevertTransaction(r.Context(), chi.URLParam(r, "txid"))
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, tx)
}

func PostTransactionMetadata(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	var m metadata.Metadata
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		apierrors.ResponseError(w, r, errorsutil.NewError(ledger.ErrValidation,
			errors.New("invalid metadata format")))
		return
	}

	if err := l.SaveMeta(r.Context(), core.MetaTargetTypeTransaction, chi.URLParam(r, "txid"), m); err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.NoContent(w)
}
