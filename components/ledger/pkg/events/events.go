package events

const (
	EventVersion = "v2"
	EventApp     = "ledger"

	EventTypeCommittedTransactions = "COMMITTED_TRANSACTIONS"
	EventTypeSavedMetadata         = "SAVED_METADATA"
	TypeRevertedTransaction        = "REVERTED_TRANSACTION"
)
