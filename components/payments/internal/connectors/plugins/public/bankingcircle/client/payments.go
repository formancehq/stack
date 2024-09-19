package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

//nolint:tagliatelle // allow for client-side structures
type Payment struct {
	PaymentID                    string      `json:"paymentId"`
	TransactionReference         string      `json:"transactionReference"`
	ConcurrencyToken             string      `json:"concurrencyToken"`
	Classification               string      `json:"classification"`
	Status                       string      `json:"status"`
	Errors                       interface{} `json:"errors"`
	ProcessedTimestamp           time.Time   `json:"processedTimestamp"`
	LatestStatusChangedTimestamp time.Time   `json:"latestStatusChangedTimestamp"`
	LastChangedTimestamp         time.Time   `json:"lastChangedTimestamp"`
	DebtorInformation            struct {
		PaymentBulkID interface{} `json:"paymentBulkId"`
		AccountID     string      `json:"accountId"`
		Account       struct {
			Account              string `json:"account"`
			FinancialInstitution string `json:"financialInstitution"`
			Country              string `json:"country"`
		} `json:"account"`
		VibanID interface{} `json:"vibanId"`
		Viban   struct {
			Account              string `json:"account"`
			FinancialInstitution string `json:"financialInstitution"`
			Country              string `json:"country"`
		} `json:"viban"`
		InstructedDate interface{} `json:"instructedDate"`
		DebitAmount    struct {
			Currency string      `json:"currency"`
			Amount   json.Number `json:"amount"`
		} `json:"debitAmount"`
		DebitValueDate time.Time   `json:"debitValueDate"`
		FxRate         interface{} `json:"fxRate"`
		Instruction    interface{} `json:"instruction"`
	} `json:"debtorInformation"`
	Transfer struct {
		DebtorAccount interface{} `json:"debtorAccount"`
		DebtorName    interface{} `json:"debtorName"`
		DebtorAddress interface{} `json:"debtorAddress"`
		Amount        struct {
			Currency string      `json:"currency"`
			Amount   json.Number `json:"amount"`
		} `json:"amount"`
		ValueDate             interface{} `json:"valueDate"`
		ChargeBearer          interface{} `json:"chargeBearer"`
		RemittanceInformation interface{} `json:"remittanceInformation"`
		CreditorAccount       interface{} `json:"creditorAccount"`
		CreditorName          interface{} `json:"creditorName"`
		CreditorAddress       interface{} `json:"creditorAddress"`
	} `json:"transfer"`
	CreditorInformation struct {
		AccountID string `json:"accountId"`
		Account   struct {
			Account              string `json:"account"`
			FinancialInstitution string `json:"financialInstitution"`
			Country              string `json:"country"`
		} `json:"account"`
		VibanID interface{} `json:"vibanId"`
		Viban   struct {
			Account              string `json:"account"`
			FinancialInstitution string `json:"financialInstitution"`
			Country              string `json:"country"`
		} `json:"viban"`
		CreditAmount struct {
			Currency string      `json:"currency"`
			Amount   json.Number `json:"amount"`
		} `json:"creditAmount"`
		CreditValueDate time.Time   `json:"creditValueDate"`
		FxRate          interface{} `json:"fxRate"`
	} `json:"creditorInformation"`
}

func (c *Client) GetPayments(ctx context.Context, page int, pageSize int) ([]Payment, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "bankingcircle", "list_payments")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/api/v1/payments/singles", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create payments request: %w", err)
	}

	q := req.URL.Query()
	q.Add("PageSize", fmt.Sprint(pageSize))
	q.Add("PageNumber", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	type response struct {
		Result   []Payment `json:"result"`
		PageInfo struct {
			CurrentPage int `json:"currentPage"`
			PageSize    int `json:"pageSize"`
		} `json:"pageInfo"`
	}

	res := response{Result: make([]Payment, 0)}
	statusCode, err := c.httpClient.Do(req, &res, nil)
	switch err {
	case nil:
		return res.Result, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, fmt.Errorf("received status code %d for get payments", statusCode)
	}
	return nil, fmt.Errorf("failed to get payments: %w", err)
}

type StatusResponse struct {
	Status string `json:"status"`
}

func (c *Client) GetPaymentStatus(ctx context.Context, paymentID string) (*StatusResponse, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "bankingcircle", "get_payment_status")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/payments/singles/%s/status", c.endpoint, paymentID), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create payments request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	var res StatusResponse
	statusCode, err := c.httpClient.Do(req, &res, nil)
	switch err {
	case nil:
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, fmt.Errorf("received status code %d for get payment status", statusCode)
	}
	return nil, fmt.Errorf("failed to get payments status: %w", err)
}
