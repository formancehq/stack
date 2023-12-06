package models

type ReconciliationStatus string

var (
	ReconciliationNotOK ReconciliationStatus = "NOT_OK"
	ReconciliationOK    ReconciliationStatus = "OK"
)

func (r ReconciliationStatus) String() string {
	return string(r)
}
