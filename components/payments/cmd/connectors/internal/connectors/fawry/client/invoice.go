package client

type InvoiceInfo struct {
	Number            string `json:"number"`
	BusinessRefNumber string `json:"businessRefNumber"`
	DueDate           string `json:"dueDate"`
	ExpiryDate        uint64 `json:"expiryDate"`
}
