package events

const (
	TopicPayments = "payments"

	EventVersion = "v1"
	EventApp     = "payments"

	EventTypeSavedPool                = "SAVED_POOL"
	EventTypeDeletePool               = "DELETED_POOL"
	EventTypeSavedPayments            = "SAVED_PAYMENT"
	EventTypeSavedAccounts            = "SAVED_ACCOUNT"
	EventTypeSavedBalances            = "SAVED_BALANCE"
	EventTypeSavedBankAccount         = "SAVED_BANK_ACCOUNT"
	EventTypeSavedTransferInitiation  = "SAVED_TRANSFER_INITIATION"
	EventTypeDeleteTransferInitiation = "DELETED_TRANSFER_INITIATION"
	EventTypeConnectorReset           = "CONNECTOR_RESET"
)
