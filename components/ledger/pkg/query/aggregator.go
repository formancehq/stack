package query

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/pkg/errors"
)

type TxVolumeAggregator struct {
	agg        *VolumeAggregator
	previousTx *TxVolumeAggregator

	PreCommitVolumes  core.AccountsAssetsVolumes
	PostCommitVolumes core.AccountsAssetsVolumes
}

func (tva *TxVolumeAggregator) FindInPreviousTxs(addr, asset string) *core.Volumes {
	current := tva.previousTx
	for current != nil {
		if v, ok := current.PostCommitVolumes[addr][asset]; ok {
			return &v
		}
		current = current.previousTx
	}
	return nil
}

func (tva *TxVolumeAggregator) Transfer(
	ctx context.Context,
	from, to, asset string,
	amount *core.MonetaryInt,
) error {
	for _, addr := range []string{from, to} {
		if !tva.PreCommitVolumes.HasAccountAndAsset(addr, asset) {
			previousVolumes := tva.FindInPreviousTxs(addr, asset)
			if previousVolumes != nil {
				tva.PreCommitVolumes.SetVolumes(addr, asset, *previousVolumes)
			} else {
				acc, err := tva.agg.store.GetAccountWithVolumes(ctx, addr)
				if err != nil {
					return errors.Wrap(err, "getting account while transferring")
				}
				tva.PreCommitVolumes.SetVolumes(addr, asset, acc.Volumes[asset])
			}
		}
		if !tva.PostCommitVolumes.HasAccountAndAsset(addr, asset) {
			tva.PostCommitVolumes.SetVolumes(addr, asset, tva.PreCommitVolumes.GetVolumes(addr, asset))
		}
	}
	tva.PostCommitVolumes.AddOutput(from, asset, amount)
	tva.PostCommitVolumes.AddInput(to, asset, amount)

	return nil
}

type VolumeAggregator struct {
	txs   []*TxVolumeAggregator
	store storage.LedgerStore
}

func (agg *VolumeAggregator) NextTx() *TxVolumeAggregator {
	var previousTx *TxVolumeAggregator
	if len(agg.txs) > 0 {
		previousTx = agg.txs[len(agg.txs)-1]
	}
	tva := &TxVolumeAggregator{
		agg:        agg,
		previousTx: previousTx,
	}
	agg.txs = append(agg.txs, tva)
	return tva
}

func newVolumeAggregator(store storage.LedgerStore) *VolumeAggregator {
	return &VolumeAggregator{
		store: store,
	}
}
