package client

import (
	"context"
	"encoding/json"
	"time"
)

const (
	transactionsEndpoint = "https://api.atlar.com/v1/transactions"
)

type Transaction struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Amount         struct {
		Currency    string `json:"currency"`
		Value       int    `json:"value"`
		StringValue string `json:"stringValue"`
	} `json:"amount"`
	Date      string `json:"date"`
	ValueDate string `json:"valueDate"`
	Account   struct {
		ID          string             `json:"id"`
		Name        string             `json:"name"`
		Bank        AccountBank        `json:"bank"`
		Affiliation AccountAffiliation `json:"affiliation"`
	} `json:"account"`
	Counterparty struct {
		ID             string      `json:"id"`
		Name           string      `json:"name"`
		PartyType      string      `json:"partyType"`
		Identifiers    interface{} `json:"identifiers"`
		ContactDetails struct {
			Address struct {
				StreetName      string      `json:"streetName"`
				StreetNumber    string      `json:"streetNumber"`
				PostalCode      string      `json:"postalCode"`
				City            string      `json:"city"`
				RawAddressLines interface{} `json:"rawAddressLines"`
				Country         string      `json:"country"`
			} `json:"address"`
		} `json:"contactDetails"`
	} `json:"counterparty"`
	CounterpartyDetails struct {
		Name           string `json:"name"`
		PartyType      string `json:"partyType"`
		ContactDetails struct {
			NationalID string `json:"nationalId"`
			Address    struct {
				StreetName      string        `json:"streetName"`
				StreetNumber    string        `json:"streetNumber"`
				PostalCode      string        `json:"postalCode"`
				City            string        `json:"city"`
				RawAddressLines []interface{} `json:"rawAddressLines"`
				Country         string        `json:"country"`
			} `json:"address"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		} `json:"contactDetails"`
		NationalIdentifier interface{} `json:"nationalIdentifier"`
		ExternalAccount    struct {
			Market        string      `json:"market"`
			Identifier    interface{} `json:"identifier"`
			RawIdentifier string      `json:"rawIdentifier"`
			Bank          struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Bic  string `json:"bic"`
			} `json:"bank"`
		} `json:"externalAccount"`
		MandateReference string `json:"mandateReference"`
	} `json:"counterpartyDetails"`
	RemittanceInformation struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"remittanceInformation"`
	Reconciliation struct {
		Status              string `json:"status"`
		BookedTransactionID string `json:"bookedTransactionId"`
	} `json:"reconciliation"`
	Characteristics struct {
		BankTransactionCode struct {
			Domain      string `json:"domain"`
			Family      string `json:"family"`
			Description string `json:"description"`
			Proprietary struct {
			} `json:"proprietary"`
		} `json:"bankTransactionCode"`
		Returned     bool `json:"returned"`
		ReturnReason struct {
			OriginalBankTransactionCode struct {
				Proprietary struct {
				} `json:"proprietary"`
			} `json:"originalBankTransactionCode"`
		} `json:"returnReason"`
		InstructedAmount struct {
			Currency    string `json:"currency"`
			Value       int    `json:"value"`
			StringValue string `json:"stringValue"`
		} `json:"instructedAmount"`
		CurrencyExchange interface{} `json:"currencyExchange"`
	} `json:"characteristics"`
	Description string    `json:"description"`
	Version     int       `json:"version"`
	Created     time.Time `json:"created"`
}

type TransactionsListResponse struct {
	ListResponse
	Items []*Transaction `json:"items"`
}

func (d *DefaultClient) Transactions(ctx context.Context,
	options ...ClientOption,
) ([]*Transaction, string, error) {
	res := RawListResponse{}

	d.SendGetListRequest(ctx, &res, options...)

	accounts := make([]*Transaction, 0)

	if len(res.Items) > 0 {
		for _, item := range res.Items {
			account := &Account{}

			err := json.Unmarshal(item, &account)
			if err != nil {
				return nil, "", err
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, res.NextToken, nil
}
