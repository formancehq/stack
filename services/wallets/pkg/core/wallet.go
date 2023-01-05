package core

import (
	"github.com/formancehq/go-libs/metadata"
	"github.com/google/uuid"
)

type Spec struct {
	Type string `json:"type"`
}

const (
	MetadataKeyWalletTransaction     = "wallets"
	MetadataKeyWalletSpecType        = "wallets/spec/type"
	MetadataKeyWalletID              = "wallets/id"
	MetadataKeyWalletName            = "wallets/name"
	MetadataKeyWalletCustomData      = "wallets/custom_data"
	MetadataKeyHoldWalletID          = "wallets/holds/wallet_id"
	MetadataKeyHoldAsset             = "wallets/holds/asset"
	MetadataKeyHoldID                = "wallets/holds/id"
	MetadataKeyWalletHoldDescription = "wallets/holds/description"
	MetadataKeyHoldVoidDestination   = "void_destination"
	MetadataKeyHoldDestination       = "destination"

	PrimaryWallet = "wallets.primary"
	HoldWallet    = "wallets.hold"
)

func WalletTransactionBaseMetadata() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletTransaction: true,
	}
}

func WalletTransactionBaseMetadataFilter() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletTransaction: true,
	}
}

func IsPrimary(v metadata.Owner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, PrimaryWallet)
}

func IsHold(v metadata.Owner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, HoldWallet)
}

func GetMetadata(v metadata.Owner, key string) any {
	return v.GetMetadata()[key]
}

func HasMetadata(v metadata.Owner, key, value string) bool {
	return GetMetadata(v, key) == value
}

func SpecType(v metadata.Owner) string {
	return GetMetadata(v, MetadataKeyWalletSpecType).(string)
}

type Wallet struct {
	ID       string            `json:"id"`
	Metadata metadata.Metadata `json:"metadata"`
	Name     string            `json:"name"`
}

type WalletWithBalances struct {
	Wallet
	Balances map[string]int32 `json:"balances"`
}

func (w Wallet) LedgerMetadata() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletSpecType:   PrimaryWallet,
		MetadataKeyWalletID:         w.ID,
		MetadataKeyWalletCustomData: map[string]any(w.Metadata),
		MetadataKeyWalletName:       w.Name,
	}
}

func NewWallet(name string, m metadata.Metadata) Wallet {
	if m == nil {
		m = metadata.Metadata{}
	}
	return Wallet{
		ID:       uuid.NewString(),
		Metadata: m,
		Name:     name,
	}
}

func WalletFromAccount(account metadata.Owner) Wallet {
	return Wallet{
		ID:       GetMetadata(account, MetadataKeyWalletID).(string),
		Metadata: GetMetadata(account, MetadataKeyWalletCustomData).(map[string]any),
		Name:     GetMetadata(account, MetadataKeyWalletName).(string),
	}
}

func WalletWithBalancesFromAccount(account interface {
	metadata.Owner
	GetBalances() map[string]int32
},
) WalletWithBalances {
	return WalletWithBalances{
		Wallet:   WalletFromAccount(account),
		Balances: account.GetBalances(),
	}
}
