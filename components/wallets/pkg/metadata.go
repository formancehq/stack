package wallet

import (
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

const (
	MetadataKeyWalletTransaction     = "wallets/transaction"
	MetadataKeyWalletSpecType        = "wallets/spec/type"
	MetadataKeyWalletID              = "wallets/id"
	MetadataKeyWalletName            = "wallets/name"
	MetadataKeyHoldWalletID          = "wallets/holds/wallet_id"
	MetadataKeyHoldAsset             = "wallets/holds/asset"
	MetadataKeyHoldSubject           = "wallets/holds/subject"
	MetadataKeyHoldID                = "wallets/holds/id"
	MetadataKeyWalletHoldDescription = "wallets/holds/description"
	MetadataKeyHoldVoidDestination   = "void_destination"
	MetadataKeyHoldDestination       = "destination"
	MetadataKeyBalanceName           = "wallets/balances/name"
	MetadataKeyBalanceExpiresAt      = "wallets/balances/expiresAt"
	MetadataKeyBalancePriority       = "wallets/balances/priority"
	MetadataKeyWalletBalance         = "wallets/balances"
	MetadataKeyCreatedAt             = "wallets/createdAt"

	PrimaryWallet = "wallets.primary"
	HoldWallet    = "wallets.hold"

	TrueValue = "true"

	MetadataKeyWalletCustomData = "wallets/custom_data"
)

func TransactionMetadata(customMetadata metadata.Metadata) map[string]any {
	if customMetadata == nil {
		customMetadata = metadata.Metadata{}
	}
	return map[string]any{
		MetadataKeyWalletTransaction: "true",
		MetadataKeyWalletCustomData:  customMetadata,
	}
}

func TransactionBaseMetadataFilter() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletTransaction: "true",
	}
}

type MetadataOwner interface {
	GetMetadata() map[string]any
}

func IsPrimary(v MetadataOwner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, PrimaryWallet)
}

func IsHold(v MetadataOwner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, HoldWallet)
}

func GetMetadata(v MetadataOwner, key string) any {
	return v.GetMetadata()[key]
}

func HasMetadata(v MetadataOwner, key, value string) bool {
	return GetMetadata(v, key) == value
}

func LedgerMetadataToWalletMetadata(m map[string]any) metadata.Metadata {
	ret := metadata.Metadata{}
	for k, v := range m {
		ret[k] = v.(string)
	}
	return ret
}

func GetCustomMetadata(owner MetadataOwner) metadata.Metadata {
	ret := owner.GetMetadata()[MetadataKeyWalletCustomData]
	if ret == nil {
		return metadata.Metadata{}
	}

	switch v := ret.(type) {
	case string:
		return metadata.Metadata{}
	case map[string]any:
		return LedgerMetadataToWalletMetadata(v)
	default:
		return metadata.Metadata{}
	}
}
