package query

import (
	"context"
	"time"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

type workerConfig struct {
	Interval time.Duration
}

type Worker struct {
	workerConfig
	stopChan chan chan struct{}
	driver   *sqlstorage.Driver
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case stopChan := <-w.stopChan:
			close(stopChan)
			return nil
		case <-time.After(w.Interval):
			if err := w.run(ctx); err != nil {
				if err == context.Canceled {
					return err
				}
				logging.FromContext(ctx).Error(err)
			}
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

func (w *Worker) run(ctx context.Context) error {

	ctx = storage.NewCQRSContext(ctx)

	ledgers, err := w.driver.GetSystemStore().ListLedgers(ctx)
	if err != nil {
		return err
	}

	for _, ledger := range ledgers {
		if err := w.processLedger(ctx, ledger); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) processLedger(ctx context.Context, ledger string) error {

	store, _, err := w.driver.GetLedgerStore(ctx, ledger, false)
	if err != nil && err != storage.ErrLedgerStoreNotFound {
		return err
	}
	if err == storage.ErrLedgerStoreNotFound {
		return nil
	}

	lastReadLogID, err := store.GetNextLogID(ctx)
	if err != nil {
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

	return nil
}

func (w *Worker) processLogs(ctx context.Context, store storage.LedgerStore, logs ...core.Log) error {
	volumeAggregator := newVolumeAggregator(store)
	for _, log := range logs {
		var err error
		switch log.Type {
		case core.NewTransactionLogType:
			tx := log.Data.(core.Transaction)
			txVolumeAggregator := volumeAggregator.NextTx()
			for _, posting := range tx.Postings {
				if err := txVolumeAggregator.Transfer(ctx, posting.Source, posting.Destination, posting.Asset, posting.Amount); err != nil {
					return errors.Wrap(err, "aggregating volumes")
				}
			}

			//TODO(gfyrag): when all the mess will be cleaned, the InsertTransactions method should be rewrite to ignore potential conflict
			// This way, we don't need any sql transactions
			if err := store.InsertTransactions(ctx, core.ExpandedTransaction{
				Transaction:       tx,
				PreCommitVolumes:  txVolumeAggregator.PreCommitVolumes,
				PostCommitVolumes: txVolumeAggregator.PostCommitVolumes,
			}); err != nil {
				return errors.Wrap(err, "inserting transactions")
			}

			for account := range txVolumeAggregator.PostCommitVolumes {
				if err := store.EnsureAccountExists(ctx, account); err != nil {
					return errors.Wrap(err, "ensuring account exists")
				}
			}

			if err := store.UpdateVolumes(ctx, txVolumeAggregator.PostCommitVolumes); err != nil {
				return errors.Wrap(err, "updating volumes")
			}

		case core.SetMetadataLogType:
			switch setMetadata := log.Data.(core.SetMetadata); setMetadata.TargetType {
			case core.MetaTargetTypeAccount:
				if err := store.UpdateAccountMetadata(ctx, setMetadata.TargetID.(string), setMetadata.Metadata); err != nil {
					return errors.Wrap(err, "updating account metadata")
				}
			case core.MetaTargetTypeTransaction:
				if err := store.UpdateTransactionMetadata(ctx, setMetadata.TargetID.(uint64), setMetadata.Metadata); err != nil {
					return errors.Wrap(err, "updating transactions metadata")
				}
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func NewWorker(config workerConfig, driver *sqlstorage.Driver) *Worker {
	return &Worker{
		stopChan:     make(chan chan struct{}),
		workerConfig: config,
		driver:       driver,
	}
}
