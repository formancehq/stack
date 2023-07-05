package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//nolint:tagliatelle // allow for client-side structures
type Payment struct {
	PaymentID            string      `json:"paymentId"`
	TransactionReference string      `json:"transactionReference"`
	ConcurrencyToken     string      `json:"concurrencyToken"`
	Classification       string      `json:"classification"`
	Status               string      `json:"status"`
	Errors               interface{} `json:"errors"`
	LastChangedTimestamp time.Time   `json:"lastChangedTimestamp"`
	DebtorInformation    struct {
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
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
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
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
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
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
		} `json:"creditAmount"`
		CreditValueDate time.Time   `json:"creditValueDate"`
		FxRate          interface{} `json:"fxRate"`
	} `json:"creditorInformation"`
}

func (c *Client) GetPayments(ctx context.Context, page int) ([]*Payment, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/api/v1/payments/singles", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("PageSize", "100")
	q.Add("PageNumber", fmt.Sprint(page))

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response body: %w", err)
	}

	type response struct {
		Result   []*Payment `json:"result"`
		PageInfo struct {
			CurrentPage int `json:"currentPage"`
			PageSize    int `json:"pageSize"`
		} `json:"pageInfo"`
	}

	var res response

	if err = json.Unmarshal(responseBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	return res.Result, nil
}
