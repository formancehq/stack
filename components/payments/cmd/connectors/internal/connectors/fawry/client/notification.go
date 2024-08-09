package client

type Notification struct {
	RequestID         string      `json:"requestId"`
	FawryRefNumber    string      `json:"fawryRefNumber"`
	MerchantRefNumber string      `json:"merchantRefNumber"`
	PaymentAmount     string      `json:"paymentAmount"`
	OrderAmount       string      `json:"orderAmount"`
	FawryFees         string      `json:"fawryFees"`
	OrderStatus       string      `json:"orderStatus"`
	PaymentMethod     string      `json:"paymentMethod"`
	MessageSignature  string      `json:"messageSignature"`
	InvoiceInfo       InvoiceInfo `json:"invoiceInfo"`
}
