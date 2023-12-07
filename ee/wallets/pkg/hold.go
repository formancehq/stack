package wallet

import (
	"math/big"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/google/uuid"
)

type DebitHold struct {
	ID          string            `json:"id"`
	WalletID    string            `json:"walletID"`
	Destination Subject           `json:"destination"`
	Asset       string            `json:"asset"`
	Metadata    metadata.Metadata `json:"metadata"`
	Description string            `json:"description"`
}

func (h DebitHold) LedgerMetadata(chart *Chart) metadata.Metadata {

	ret := metadata.Metadata{
		MetadataKeyWalletSpecType: HoldWallet,
		MetadataKeyHoldWalletID:   h.WalletID,
		MetadataKeyHoldID:         h.ID,
		MetadataKeyHoldAsset:      h.Asset,
		MetadataKeyHoldVoidDestination: metadata.MarshalValue(map[string]any{
			"type":  "account",
			"value": chart.GetMainBalanceAccount(h.WalletID),
		}),
		MetadataKeyHoldDestination: metadata.MarshalValue(map[string]any{
			"type":  "account",
			"value": h.Destination.getAccount(chart),
		}),
		MetadataKeyWalletHoldDescription: h.Description,
		MetadataKeyHoldSubject: metadata.MarshalValue(map[string]any{
			"type":       h.Destination.Type,
			"identifier": h.Destination.Identifier,
			"balance":    h.Destination.Balance,
		}),
	}

	ret = ret.Merge(EncodeCustomMetadata(h.Metadata))

	return ret
}

func NewDebitHold(walletID string, destination Subject, asset, description string, md metadata.Metadata) DebitHold {
	return DebitHold{
		ID:          uuid.NewString(),
		WalletID:    walletID,
		Destination: destination,
		Asset:       asset,
		Metadata:    md,
		Description: description,
	}
}

func DebitHoldFromLedgerAccount(account Account) DebitHold {
	subject := metadata.UnmarshalValue[Subject](account.GetMetadata()[MetadataKeyHoldSubject])

	hold := DebitHold{}
	hold.ID = account.GetMetadata()[MetadataKeyHoldID]
	hold.WalletID = account.GetMetadata()[MetadataKeyHoldWalletID]
	hold.Destination = subject
	hold.Asset = account.GetMetadata()[MetadataKeyHoldAsset]
	hold.Description = account.GetMetadata()[MetadataKeyWalletHoldDescription]
	hold.Metadata = ExtractCustomMetadata(account)

	return hold
}

type ExpandedDebitHold struct {
	DebitHold
	OriginalAmount *big.Int `json:"originalAmount"`
	Remaining      *big.Int `json:"remaining"`
}

func (h ExpandedDebitHold) IsClosed() bool {
	return h.Remaining.Uint64() == 0
}

func ExpandedDebitHoldFromLedgerAccount(account AccountWithVolumesAndBalances) ExpandedDebitHold {
	hold := ExpandedDebitHold{
		DebitHold: DebitHoldFromLedgerAccount(account.Account),
	}
	hold.OriginalAmount = account.GetVolumes()[hold.Asset]["input"]
	hold.Remaining = account.GetBalances()[hold.Asset]
	return hold
}

type ConfirmHold struct {
	HoldID    string `json:"holdID"`
	Amount    *big.Int
	Reference string
	Final     bool
}

func (c ConfirmHold) resolveAmount(hold ExpandedDebitHold) (uint64, error) {
	amount := hold.Remaining.Uint64()
	if c.Amount.Uint64() != 0 {
		if c.Amount.Uint64() > amount {
			return 0, ErrInsufficientFundError
		}
		amount = c.Amount.Uint64()
	}
	return amount, nil
}

type VoidHold struct {
	HoldID string `json:"holdID"`
}
