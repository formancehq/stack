package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverseTransaction(t *testing.T) {
	tx := &ExpandedTransaction{
		Transaction: Transaction{
			TransactionData: TransactionData{
				Postings: Postings{
					{
						Source:      "world",
						Destination: "users:001",
						Amount:      big.NewInt(100),
						Asset:       "COIN",
					},
					{
						Source:      "users:001",
						Destination: "payments:001",
						Amount:      big.NewInt(100),
						Asset:       "COIN",
					},
				},
				Reference: "foo",
			},
		},
	}

	expected := TransactionData{
		Postings: Postings{
			{
				Source:      "payments:001",
				Destination: "users:001",
				Amount:      big.NewInt(100),
				Asset:       "COIN",
			},
			{
				Source:      "users:001",
				Destination: "world",
				Amount:      big.NewInt(100),
				Asset:       "COIN",
			},
		},
		Reference: "revert_foo",
	}
	require.Equal(t, expected, tx.Reverse())
}
