package wallet

import (
	"strings"

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

	MetadataKeyWalletCustomDataPrefix = "wallets/custom_data_"
)

func TransactionMetadata(customMetadata metadata.Metadata) map[string]string {
	if customMetadata == nil {
		customMetadata = metadata.Metadata{}
	}
	return metadata.Metadata{
		MetadataKeyWalletTransaction: "true",
	}.Merge(EncodeCustomMetadata(customMetadata))
}

func TransactionBaseMetadataFilter() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletTransaction: "true",
	}
}

type MetadataOwner interface {
	GetMetadata() map[string]string
}

func IsPrimary(v MetadataOwner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, PrimaryWallet)
}

func IsHold(v MetadataOwner) bool {
	return HasMetadata(v, MetadataKeyWalletSpecType, HoldWallet)
}

func GetMetadata(v MetadataOwner, key string) string {
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

func EncodeCustomMetadata(m metadata.Metadata) metadata.Metadata {
	ret := metadata.Metadata{}
	for key, value := range m {
		ret[MetadataKeyWalletCustomDataPrefix+key] = value
	}
	return ret
}

func ExtractCustomMetadata(account interface {
	GetMetadata() map[string]string
}) metadata.Metadata {
	ret := metadata.Metadata{}
	for key, value := range account.GetMetadata() {
		if strings.HasPrefix(key, MetadataKeyWalletCustomDataPrefix) {
			ret[strings.TrimPrefix(key, MetadataKeyWalletCustomDataPrefix)] = value
		}
	}
	return ret
}
