package models

import "errors"

type TransferInitiationStatus int

const (
	TransferInitiationStatusWaitingForValidation TransferInitiationStatus = iota
	TransferInitiationStatusProcessing
	TransferInitiationStatusProcessed
	TransferInitiationStatusFailed
	TransferInitiationStatusRejected
	TransferInitiationStatusValidated
	TransferInitiationStatusAskRetried
	TransferInitiationStatusAskReversed
	TransferInitiationStatusReverseProcessing
	TransferInitiationStatusReverseFailed
	TransferInitiationStatusPartiallyReversed
	TransferInitiationStatusReversed
)

func (s TransferInitiationStatus) String() string {
	return [...]string{
		"WAITING_FOR_VALIDATION",
		"PROCESSING",
		"PROCESSED",
		"FAILED",
		"REJECTED",
		"VALIDATED",
		"ASK_RETRIED",
		"ASK_REVERSED",
		"REVERSE_PROCESSING",
		"REVERSE_FAILED",
		"PARTIALLY_REVERSED",
		"REVERSED",
	}[s]
}

func TransferInitiationStatusFromString(s string) (TransferInitiationStatus, error) {
	switch s {
	case "WAITING_FOR_VALIDATION":
		return TransferInitiationStatusWaitingForValidation, nil
	case "PROCESSING":
		return TransferInitiationStatusProcessing, nil
	case "PROCESSED":
		return TransferInitiationStatusProcessed, nil
	case "FAILED":
		return TransferInitiationStatusFailed, nil
	case "REJECTED":
		return TransferInitiationStatusRejected, nil
	case "VALIDATED":
		return TransferInitiationStatusValidated, nil
	case "ASK_RETRIED":
		return TransferInitiationStatusAskRetried, nil
	case "ASK_REVERSED":
		return TransferInitiationStatusAskReversed, nil
	case "REVERSE_PROCESSING":
		return TransferInitiationStatusReverseProcessing, nil
	case "REVERSE_FAILED":
		return TransferInitiationStatusReverseFailed, nil
	case "PARTIALLY_REVERSED":
		return TransferInitiationStatusPartiallyReversed, nil
	case "REVERSED":
		return TransferInitiationStatusReversed, nil
	default:
		return TransferInitiationStatusWaitingForValidation, errors.New("invalid status")
	}
}

type TransferReversalStatus int

const (
	TransferReversalStatusProcessing TransferReversalStatus = iota
	TransferReversalStatusProcessed
	TransferReversalStatusFailed
)

func (s TransferReversalStatus) String() string {
	return [...]string{
		"CREATED",
		"PROCESSING",
		"PROCESSED",
		"FAILED",
	}[s]
}

func TransferReversalStatusFromString(s string) (TransferReversalStatus, error) {
	switch s {
	case "PROCESSING":
		return TransferReversalStatusProcessing, nil
	case "PROCESSED":
		return TransferReversalStatusProcessed, nil
	case "FAILED":
		return TransferReversalStatusFailed, nil
	default:
		return TransferReversalStatusProcessing, errors.New("invalid status")
	}
}

func (s TransferReversalStatus) ToTransferInitiationStatus(isFullyReversed bool) TransferInitiationStatus {
	switch s {
	case TransferReversalStatusProcessing:
		return TransferInitiationStatusReverseProcessing
	case TransferReversalStatusProcessed:
		if isFullyReversed {
			return TransferInitiationStatusReversed
		}
		return TransferInitiationStatusPartiallyReversed
	case TransferReversalStatusFailed:
		return TransferInitiationStatusReverseFailed
	default:
		return TransferInitiationStatusProcessed
	}
}
