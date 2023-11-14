// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type CurrencyCloudConfig struct {
	APIKey string `json:"apiKey"`
	// The endpoint to use for the API. Defaults to https://devapi.currencycloud.com
	Endpoint *string `json:"endpoint,omitempty"`
	// Username of the API Key holder
	LoginID string `json:"loginID"`
	Name    string `json:"name"`
	// The frequency at which the connector will fetch transactions
	PollingPeriod *string `json:"pollingPeriod,omitempty"`
}
