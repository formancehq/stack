// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type ConnectorsResponseData struct {
	ConnectorID string    `json:"connectorID"`
	Name        string    `json:"name"`
	Provider    Connector `json:"provider"`
}

// ConnectorsResponse - OK
type ConnectorsResponse struct {
	Data []ConnectorsResponseData `json:"data"`
}
