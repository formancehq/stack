package client

type Notification struct {
	RequestID string `json:"requestId"`

	Amount string `json:"amount"`

	PaymentID     string `json:"paymentId"`
	TransactionID string `json:"transactionId"`

	BillingAccount       string `json:"billingAcct"`
	ExtraBillingAccounts any    `json:"extraBillingAccts"`

	CustomerReferenceNumber string `json:"customerReferenceNumber"`
	TerminalID              string `json:"terminalId"`
	BillNumber              string `json:"billNumber"`
}
