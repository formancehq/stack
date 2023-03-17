package bus

import (
	"time"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/events/ledger"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	EventApp = "ledger"
)

func newEventCommittedTransactions(ledgerName string, txs ...core.ExpandedTransaction) (*events.Event, error) {
	transactions, err := buildTransactionsProto(txs...)
	if err != nil {
		return nil, err
	}

	preCommitVolumes, err := buildAccountsAssetsVolumesProto(core.AggregatePreCommitVolumes(txs...))
	if err != nil {
		return nil, err
	}
	postCommitVolumes, err := buildAccountsAssetsVolumesProto(core.AggregatePostCommitVolumes(txs...))
	if err != nil {
		return nil, err
	}

	return &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       EventApp,
		Event: &events.Event_TransactionsCommitted{
			TransactionsCommitted: &ledger.TransactionsCommitted{
				Ledger:            ledgerName,
				Transactions:      transactions,
				PreCommitVolumes:  preCommitVolumes,
				PostCommitVolumes: postCommitVolumes,
			},
		},
	}, nil
}

func newEventSavedMetadata(ledgerName, targetID string, targetType core.TargetType, metadata core.Metadata) (*events.Event, error) {
	tt := ledger.TargetType_TARGET_TYPE_UNKNOWN
	switch targetType {
	case core.MetaTargetTypeAccount:
		tt = ledger.TargetType_TARGET_TYPE_ACCOUNT
	case core.MetaTargetTypeTransaction:
		tt = ledger.TargetType_TARGET_TYPE_TRANSACTION
	}

	md, err := structpb.NewValue(metadata)
	if err != nil {
		return nil, err
	}

	return &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       EventApp,
		Event: &events.Event_MetadataSaved{
			MetadataSaved: &ledger.MetadataSaved{
				Ledger:     ledgerName,
				TargetType: tt,
				TargetId:   targetID,
				Metadata:   md,
			},
		},
	}, nil
}

func newEventRevertedTransaction(
	ledgerName string,
	revertedTransaction,
	revertTransaction *core.ExpandedTransaction,
) (*events.Event, error) {
	revertedTrans, err := buildTransactionProto(revertedTransaction)
	if err != nil {
		return nil, err
	}

	revertTrans, err := buildTransactionProto(revertTransaction)
	if err != nil {
		return nil, err
	}

	return &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       EventApp,
		Event: &events.Event_TransactionReverted{
			TransactionReverted: &ledger.TransactionReverted{
				Ledger:              ledgerName,
				RevertedTransaction: revertedTrans,
				RevertTransaction:   revertTrans,
			},
		},
	}, nil
}

func buildTransactionsProto(txs ...core.ExpandedTransaction) ([]*ledger.ExpandedTransaction, error) {
	transactions := make([]*ledger.ExpandedTransaction, len(txs))
	var err error
	for i, tx := range txs {
		transactions[i], err = buildTransactionProto(&tx)
		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func buildTransactionProto(tx *core.ExpandedTransaction) (*ledger.ExpandedTransaction, error) {
	postings := make([]*ledger.Posting, len(tx.Postings))
	for j, posting := range tx.Postings {
		amount, err := posting.Amount.MarshalText()
		if err != nil {
			return nil, err
		}
		postings[j] = &ledger.Posting{
			Source:      posting.Source,
			Destination: posting.Destination,
			Amount:      amount,
			Asset:       posting.Asset,
		}
	}

	metadata, err := structpb.NewValue(tx.Metadata)
	if err != nil {
		return nil, err
	}

	transaction := &ledger.ExpandedTransaction{
		Transaction: &ledger.Transaction{
			Id:        tx.ID,
			Postings:  postings,
			Reference: tx.Reference,
			Metadata:  metadata,
			Timestamp: timestamppb.New(tx.Timestamp),
		},
		PreCommitVolumes:  &ledger.AccountsAssetsVolumes{},
		PostCommitVolumes: &ledger.AccountsAssetsVolumes{},
	}

	return transaction, nil
}

func buildAccountsAssetsVolumesProto(volumes core.AccountsAssetsVolumes) (*ledger.AccountsAssetsVolumes, error) {
	ac := &ledger.AccountsAssetsVolumes{
		Accounts: make(map[string]*ledger.AssetsVolumes),
	}
	for account, assets := range volumes {
		ac.Accounts[account] = &ledger.AssetsVolumes{
			Assets: make(map[string]*ledger.VolumeWithBalance),
		}
		for asset, volume := range assets {
			input, err := volume.Input.MarshalText()
			if err != nil {
				return nil, err
			}
			output, err := volume.Output.MarshalText()
			if err != nil {
				return nil, err
			}
			balance, err := volume.Input.Sub(volume.Output).MarshalText()
			if err != nil {
				return nil, err
			}

			ac.Accounts[account].Assets[asset] = &ledger.VolumeWithBalance{
				Input:   input,
				Output:  output,
				Balance: balance,
			}
		}
	}

	return ac, nil
}
