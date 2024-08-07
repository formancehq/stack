package utils

type AccountIdentifier struct {
    OrganisationName string
    Ledger           string
    AccountName      string
}

func NewAccountIdentifier(org, ledger, account string) AccountIdentifier {
    return AccountIdentifier{OrganisationName: org, Ledger: ledger, AccountName: account}
}