package internal

import (
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

type Info struct {
	ID        int64     `json:"id"`
	Reference string    `json:"reference"`
	Date      time.Time `json:"date"`
}

type Posting struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Asset       string   `json:"asset"`
	Amount      *big.Int `json:"amount"`
}

type Summary struct {
	Account      string   `json:"account"`
	Asset        string   `json:"asset"`
	Movement     *big.Int `json:"movement"`
	FinalBalance *big.Int `json:"finalBalance"`
}

type ExportTransaction struct {
	Info     Info                   `json:"info"`
	Postings []*Posting             `json:"postings"`
	Summary  []*Summary             `json:"summaries,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func WithPosting(posting []*shared.Posting) []*Posting {
	var postings []*Posting
	for _, posting := range posting {
		postings = append(postings, &Posting{
			posting.Source, posting.Destination, posting.Asset, posting.Amount,
		})
	}
	return postings
}

func WithSummary(tx *ExpandedTransaction) []*Summary {
	var summaries []*Summary
	for account, postCommitVolume := range tx.PostCommitVolumes {
		for asset, volumes := range postCommitVolume {
			movement := big.NewInt(0)
			movement = movement.Sub(volumes.Balance, (tx.PreCommitVolumes[account][asset]).Balance)
			summaries = append(summaries, &Summary{
				account, asset, movement, volumes.Balance,
			})
		}
	}
	return summaries
}

func WithInfo(tx *ExpandedTransaction) Info {
	return Info{
		tx.Txid,
		*tx.Reference,
		tx.Timestamp,
	}
}

func NewExportExpandedTransaction(tx *ExpandedTransaction) *ExportTransaction {
	return &ExportTransaction{
		WithInfo(tx),
		WithPosting(tx.Postings),
		WithSummary(tx),
		tx.Metadata,
	}
}

func NewExportTransaction(tx *Transaction) *ExportTransaction {
	return &ExportTransaction{
		Info{
			tx.Txid,
			*tx.Reference,
			tx.Timestamp,
		},
		WithPosting(tx.Postings),
		make([]*Summary, 0),
		tx.Metadata,
	}
}
