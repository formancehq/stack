package query

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger/aggregator"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

var (
	DefaultWorkerConfig = workerConfig{
		ChanSize: 100,
	}
)

type workerConfig struct {
	ChanSize int
}

type logHolder struct {
	*core.LogHolder
	store storage.LedgerStore
}

type Worker struct {
	workerConfig
	ctx                context.Context
	logChan            chan logHolder
	stopChan           chan chan struct{}
	driver             storage.Driver
	monitor            Monitor
	lastProcessedLogID *uint64
}

func (w *Worker) Run(ctx context.Context) error {
	logging.FromContext(ctx).Debugf("Start CQRS worker")

	w.ctx = ctx

	for {
		select {
		case <-w.ctx.Done():
			// Stop the worker if the context is done
			return w.ctx.Err()
		default:
			if err := w.run(); err != nil {
				// TODO(polo): add metrics
				if err == context.Canceled {
					// Stop the worker if the context is canceled
					return err
				}

				// Restart the worker if there is an error
			} else {
				// No error was returned, it means the worker was stopped
				// using the stopChan, let's stop this loop too
				return nil
			}
		}
	}
}

func (w *Worker) run() error {
	if err := w.initLedgers(w.ctx); err != nil {
		if err == context.Canceled {
			logging.FromContext(w.ctx).Debugf("CQRS worker canceled")
		} else {
			logging.FromContext(w.ctx).Errorf("CQRS worker error: %s", err)
		}

		return err
	}

	for {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case stopChan := <-w.stopChan:
			logging.FromContext(w.ctx).Debugf("CQRS worker stopped")
			close(stopChan)
			return nil
		case wl := <-w.logChan:
			if w.lastProcessedLogID != nil && wl.Log.ID <= *w.lastProcessedLogID {
				close(wl.Ingested)
				continue
			}
			if err := w.processLogs(w.ctx, wl.store, *wl.Log); err != nil {
				if err == context.Canceled {
					logging.FromContext(w.ctx).Debugf("CQRS worker canceled")
				} else {
					logging.FromContext(w.ctx).Errorf("CQRS worker error: %s", err)
				}
				close(wl.Ingested)

				// Return the error to restart the worker
				return err
			}

			if err := wl.store.UpdateNextLogID(w.ctx, wl.Log.ID+1); err != nil {
				logging.FromContext(w.ctx).Errorf("CQRS worker error: %s", err)
				close(wl.Ingested)
				// TODO(polo/gfyrag): add indempotency tests
				// Return the error to restart the worker
				return err
			}
			close(wl.Ingested)
		}
	}
}

func (w *Worker) Stop(ctx context.Context) error {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.stopChan <- ch:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	}

	return nil
}

func (w *Worker) initLedgers(ctx context.Context) error {
	ledgers, err := w.driver.GetSystemStore().ListLedgers(ctx)
	if err != nil {
		return err
	}

	for _, ledger := range ledgers {
		if err := w.initLedger(ctx, ledger); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) initLedger(ctx context.Context, ledger string) error {
	store, _, err := w.driver.GetLedgerStore(ctx, ledger, false)
	if err != nil && err != storage.ErrLedgerStoreNotFound {
		return err
	}
	if err == storage.ErrLedgerStoreNotFound {
		return nil
	}

	if !store.IsInitialized() {
		return nil
	}

	lastReadLogID, err := store.GetNextLogID(ctx)
	if err != nil && !storage.IsNotFound(err) {
		return errors.Wrap(err, "reading last log")
	}

	logs, err := store.ReadLogsStartingFromID(ctx, lastReadLogID)
	if err != nil {
		return errors.Wrap(err, "reading logs since last ID")
	}

	if len(logs) == 0 {
		return nil
	}

	if err := w.processLogs(ctx, store, logs...); err != nil {
		return errors.Wrap(err, "processing logs")
	}

	if err := store.UpdateNextLogID(ctx, logs[len(logs)-1].ID+1); err != nil {
		return errors.Wrap(err, "updating last read log")
	}
	lastProcessedLogID := logs[len(logs)-1].ID
	w.lastProcessedLogID = &lastProcessedLogID

	return nil
}

func (w *Worker) processLogs(ctx context.Context, store storage.LedgerStore, logs ...core.Log) error {

	accountsToUpdate, ensureAccountsExist, transactionsToInsert,
		transactionsToUpdate, volumesToUpdate, err := w.buildData(ctx, store, logs...)
	if err != nil {
		return errors.Wrap(err, "building data")
	}

	return store.RunInTransaction(ctx, func(ctx context.Context, tx storage.LedgerStore) error {
		if len(accountsToUpdate) > 0 {
			if err := tx.UpdateAccountsMetadata(ctx, accountsToUpdate); err != nil {
				return errors.Wrap(err, "updating accounts metadata")
			}
		}

		if len(transactionsToInsert) > 0 {
			if err := tx.InsertTransactions(ctx, transactionsToInsert...); err != nil {
				return errors.Wrap(err, "inserting transactions")
			}
		}

		if len(transactionsToUpdate) > 0 {
			if err := tx.UpdateTransactionsMetadata(ctx, transactionsToUpdate...); err != nil {
				return errors.Wrap(err, "updating transactions")
			}
		}

		if len(ensureAccountsExist) > 0 {
			if err := tx.EnsureAccountsExist(ctx, ensureAccountsExist); err != nil {
				return errors.Wrap(err, "ensuring accounts exist")
			}
		}

		if len(volumesToUpdate) > 0 {
			return tx.UpdateVolumes(ctx, volumesToUpdate...)
		}

		return nil
	})
}

func (w *Worker) buildData(
	ctx context.Context,
	store storage.LedgerStore,
	logs ...core.Log,
) ([]core.Account, []string, []core.ExpandedTransaction, []core.TransactionWithMetadata, []core.AccountsAssetsVolumes, error) {
	var accountsToUpdate []core.Account
	var ensureAccountsExist []string

	var transactionsToInsert []core.ExpandedTransaction
	var transactionsToUpdate []core.TransactionWithMetadata

	var volumesToUpdate []core.AccountsAssetsVolumes

	volumeAggregator := aggregator.Volumes(store)

	for _, log := range logs {
		switch log.Type {
		case core.NewTransactionLogType:
			payload := log.Data.(core.NewTransactionLogPayload)
			txVolumeAggregator, err := volumeAggregator.NextTxWithPostings(ctx, payload.Transaction.Postings...)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}

			if payload.AccountMetadata != nil {
				for account, metadata := range payload.AccountMetadata {
					accountsToUpdate = append(accountsToUpdate, core.Account{
						Address:  account,
						Metadata: metadata,
					})
				}
			}

			expandedTx := core.ExpandedTransaction{
				Transaction:       payload.Transaction,
				PreCommitVolumes:  txVolumeAggregator.PreCommitVolumes,
				PostCommitVolumes: txVolumeAggregator.PostCommitVolumes,
			}

			transactionsToInsert = append(transactionsToInsert, expandedTx)

			for account := range txVolumeAggregator.PostCommitVolumes {
				ensureAccountsExist = append(ensureAccountsExist, account)
			}

			volumesToUpdate = append(volumesToUpdate, txVolumeAggregator.PostCommitVolumes)

			// if w.monitor != nil {
			// 	w.monitor.CommittedTransactions(ctx, store.Name(), expandedTx)
			// 	for account, metadata := range payload.AccountMetadata {
			// 		w.monitor.SavedMetadata(ctx, store.Name(), core.MetaTargetTypeAccount, account, metadata)
			// 	}
			// }

		case core.SetMetadataLogType:
			setMetadata := log.Data.(core.SetMetadataLogPayload)
			switch setMetadata.TargetType {
			case core.MetaTargetTypeAccount:
				accountsToUpdate = append(accountsToUpdate, core.Account{
					Address:  setMetadata.TargetID.(string),
					Metadata: setMetadata.Metadata,
				})
			case core.MetaTargetTypeTransaction:
				transactionsToUpdate = append(transactionsToUpdate, core.TransactionWithMetadata{
					ID:       setMetadata.TargetID.(uint64),
					Metadata: setMetadata.Metadata,
				})
			}
			// if w.monitor != nil {
			// 	w.monitor.SavedMetadata(ctx, store.Name(), store.Name(), fmt.Sprint(setMetadata.TargetID), setMetadata.Metadata)
			// }
		case core.RevertedTransactionLogType:
			payload := log.Data.(core.RevertedTransactionLogPayload)
			transactionsToUpdate = append(transactionsToUpdate, core.TransactionWithMetadata{
				ID:       payload.RevertedTransactionID,
				Metadata: core.RevertedMetadata(payload.RevertTransaction.ID),
			})
			txVolumeAggregator, err := volumeAggregator.NextTxWithPostings(ctx, payload.RevertTransaction.Postings...)
			if err != nil {
				return nil, nil, nil, nil, nil, errors.Wrap(err, "aggregating volumes")
			}

			expandedTx := core.ExpandedTransaction{
				Transaction:       payload.RevertTransaction,
				PreCommitVolumes:  txVolumeAggregator.PreCommitVolumes,
				PostCommitVolumes: txVolumeAggregator.PostCommitVolumes,
			}
			transactionsToInsert = append(transactionsToInsert, expandedTx)

			// if w.monitor != nil {
			// 	revertedTx, err := store.GetTransaction(ctx, payload.RevertedTransactionID)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	w.monitor.RevertedTransaction(ctx, store.Name(), revertedTx, &expandedTx)
			// }
		}
	}

	return accountsToUpdate, ensureAccountsExist, transactionsToInsert,
		transactionsToUpdate, volumesToUpdate, nil
}

func (w *Worker) QueueLog(ctx context.Context, log *core.LogHolder, store storage.LedgerStore) {
	select {
	case <-w.ctx.Done():
	case w.logChan <- logHolder{
		LogHolder: log,
		store:     store,
	}:
	}
}

func NewWorker(config workerConfig, driver storage.Driver, monitor Monitor) *Worker {
	return &Worker{
		logChan:      make(chan logHolder, config.ChanSize),
		stopChan:     make(chan chan struct{}),
		workerConfig: config,
		driver:       driver,
		monitor:      monitor,
	}
}
