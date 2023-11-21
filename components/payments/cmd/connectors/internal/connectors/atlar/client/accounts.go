package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	accountsEndpointPath = "/v1/accounts"
)

type AccountsListResponse struct {
	ListResponse
	Items []*Account `json:"items"`
}

type Account struct {
	ID                     string                        `json:"id"`
	OrganizationID         string                        `json:"organizationId"`
	ThirdPartyID           string                        `json:"thirdPartyId"`
	Bank                   AccountBank                   `json:"bank"`
	Affiliation            AccountAffiliation            `json:"affiliation"`
	Identifiers            []AccountIdentifier           `json:"identifiers"`
	ConnectivityReferences AccountConnectivityReferences `json:"connectivityReferences"`
	Balance                AccountBalance                `json:"balance"`
	Market                 string                        `json:"market"`
	Bic                    string                        `json:"bic"`
	Currency               string                        `json:"currency"`
	Owner                  AccountOwner                  `json:"owner"`
	Name                   string                        `json:"name"`
	Created                time.Time                     `json:"created"`
	Updated                time.Time                     `json:"updated"`
	Alias                  string                        `json:"alias"`
	Capabilities           AccountCapabilities           `json:"capabilities"`
	Fictive                bool                          `json:"fictive"`
	Version                int                           `json:"version"`
}

type AccountBank struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bic  string `json:"bic"`
}

type AccountBalance struct {
	ID             string               `json:"id"`
	OrganizationID string               `json:"organizationId"`
	AccountID      string               `json:"accountId"`
	Amount         AccountBalanceAmount `json:"amount"`
	Type           string               `json:"type"`
	ReportedType   string               `json:"reportedType"`
	Timestamp      time.Time            `json:"timestamp"`
	LocalDate      string               `json:"localDate"`
	Version        int                  `json:"version"`
}

type AccountBalanceAmount struct {
	Currency    string `json:"currency"`
	Value       int    `json:"value"`
	StringValue string `json:"stringValue"`
}

type AccountOwner struct {
	Name string `json:"name"`
}

type AccountAffiliation struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ThirdPartyID string `json:"thirdPartyId"`
}

const (
	AccountIdentifierTypeIban          = "IBAN"
	AccountIdentifierTypeBban          = "BBAN"
	AccountIdentifierTypeAccountNumber = "ACCOUNT_NUMBER"
	AccountIdentifierTypeSeBg          = "SE_BG"
	AccountIdentifierTypeSePg          = "SE_PG"
)

type AccountIdentifier struct {
	Market     string `json:"market"`
	Type       string `json:"type"`
	Number     string `json:"number"`
	HolderName string `json:"holderName"`
}

type AccountConnectivityReferences struct {
	BanksCustomerID string `json:"banksCustomerId"`
	BanksAccountID  string `json:"banksAccountId"`
}

type AccountCapabilities struct {
	TransferCapabilities    []interface{} `json:"transferCapabilities"`
	DirectDebitCapabilities []interface{} `json:"directDebitCapabilities"`
	SignatureCapabilities   []interface{} `json:"signatureCapabilities"`
}

// https://docs.atlar.com/reference/get_v1-accounts
// TODO: consume token for pagination
func (d *DefaultClient) Accounts(ctx context.Context,
	options ...ClientOption,
) ([]*Account, string, error) {
	baseUrl, err := url.Parse(d.baseUrl)
	if err != nil {
		return nil, "", errors.Wrap(err, "parsing base URL")
	}
	path, err := url.Parse(accountsEndpointPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "parsing endpoint path")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl.ResolveReference(path).String(), nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.Apply(req)
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(d.accessKey, d.secret)

	var httpResponse *http.Response

	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, "", errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type listResponse struct {
		AccountsListResponse
		Items []json.RawMessage `json:"items"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, "", errors.Wrap(err, "decoding response")
	}

	accounts := make([]*Account, 0)

	if len(rsp.Items) > 0 {
		for _, item := range rsp.Items {
			account := &Account{}

			err = json.Unmarshal(item, &account)
			if err != nil {
				return nil, "", err
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, rsp.NextToken, nil
}
