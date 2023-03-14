package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/ledger/pkg/aggregator"
	"github.com/formancehq/ledger/pkg/cache"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/machine"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/pkg/errors"
)

func (l *Ledger) runScript(ctx context.Context, script core.ScriptData) (*core.ExpandedTransaction, map[string]core.Metadata, error) {

	if script.Timestamp.IsZero() {
		script.Timestamp = core.Now()
	} else {
		script.Timestamp = script.Timestamp.UTC().Round(core.DatePrecision)
	}

	lastTx, err := l.dbCache.GetLastTransaction(ctx)
	if err != nil {
		return nil, nil, err
	}

	past := false
	if lastTx != nil && script.Timestamp.Before(lastTx.Date) {
		past = true
	}
	if past && !l.allowPastTimestamps {
		return nil, nil, NewValidationError(fmt.Sprintf(
			"cannot pass a timestamp prior to the last transaction: %s (passed) is %s before %s (last)",
			script.Timestamp.Format(time.RFC3339Nano),
			lastTx.Date.Sub(script.Timestamp),
			lastTx.Date.Format(time.RFC3339Nano)))
	}

	result, err := machine.Run(ctx, l.dbCache, script)
	if err != nil {
		return nil, nil, err
	}

	if len(result.Postings) == 0 {
		return nil, nil, NewValidationError("transaction has no postings")
	}

	var nextTxId uint64
	if lastTx != nil {
		nextTxId = lastTx.ID + 1
	}

	//TODO(gfyrag): Since the CQRS is in place, this code is really not safe as a transaction could be in logs but not available
	//on global store
	if script.Reference != "" {
		txs, err := l.GetTransactions(ctx, *storage.NewTransactionsQuery().WithReferenceFilter(script.Reference))
		if err != nil {
			return nil, nil, errors.Wrap(err, "get transactions with reference")
		}
		if len(txs.Data) > 0 {
			return nil, nil, NewConflictError()
		}
	}
	vAggr := aggregator.Volumes(l.dbCache)
	txVolumeAggr := vAggr.NextTx()
	postings := make([]core.Posting, len(result.Postings))
	for j, posting := range result.Postings {
		if err := txVolumeAggr.Transfer(ctx,
			posting.Source, posting.Destination, posting.Asset, posting.Amount); err != nil {
			return nil, nil, errors.Wrap(err, "transferring volumes")
		}
		postings[j] = core.Posting{
			Source:      posting.Source,
			Destination: posting.Destination,
			Amount:      posting.Amount,
			Asset:       posting.Asset,
		}
	}

	return &core.ExpandedTransaction{
		Transaction: core.Transaction{
			TransactionData: core.TransactionData{
				Postings:  postings,
				Reference: script.Reference,
				Metadata:  result.Metadata,
				Timestamp: script.Timestamp,
			},
			ID: nextTxId,
		},
		PreCommitVolumes:  txVolumeAggr.PreCommitVolumes,
		PostCommitVolumes: txVolumeAggr.PostCommitVolumes,
	}, result.AccountMetadata, nil
}

func (l *Ledger) CreateTransactionAndWait(ctx context.Context, preview bool, script core.ScriptData) (*core.ExpandedTransaction, *LogHandler, error) {
	ret, logs, err := l.CreateTransaction(ctx, preview, script)
	if err != nil {
		return nil, nil, err
	}
	if err := logs.Wait(ctx); err != nil {
		return nil, nil, err
	}
	return ret, logs, nil
}

func (l *Ledger) CreateTransaction(ctx context.Context, preview bool, script core.ScriptData) (*core.ExpandedTransaction, *LogHandler, error) {

	// TODO: Add LockWithContext with multi lock level accounts and reference check
	unlock, err := l.locker.Lock(ctx, l.store.Name())
	if err != nil {
		panic(err)
	}
	defer unlock(context.Background()) // Use a background context instead of the request one as it could have been cancelled

	expandedTx, accountMetadata, err := l.runScript(ctx, script)
	if err != nil {
		return nil, nil, err
	}
	if preview {
		return expandedTx, &LogHandler{}, nil
	}

	logHandler, err := writeLog(ctx, l.store.AppendLogs, core.NewTransactionLog(expandedTx.Transaction, accountMetadata))
	if err != nil {
		return nil, nil, err
	}

	l.dbCache.Update(ctx, &cache.TxInfo{
		Date: expandedTx.Timestamp,
		ID:   expandedTx.ID,
	}, expandedTx.PostCommitVolumes)

	return expandedTx, logHandler, nil
}
