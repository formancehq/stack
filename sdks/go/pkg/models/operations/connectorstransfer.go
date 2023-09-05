// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type ConnectorsTransferRequest struct {
	TransferRequest shared.TransferRequest `request:"mediaType=application/json"`
	// The name of the connector.
	Connector shared.Connector `pathParam:"style=simple,explode=false,name=connector"`
}

func (o *ConnectorsTransferRequest) GetTransferRequest() shared.TransferRequest {
	if o == nil {
		return shared.TransferRequest{}
	}
	return o.TransferRequest
}

func (o *ConnectorsTransferRequest) GetConnector() shared.Connector {
	if o == nil {
		return shared.Connector("")
	}
	return o.Connector
}

type ConnectorsTransferResponse struct {
	ContentType string
	// Error
	ErrorResponse *shared.ErrorResponse
	StatusCode    int
	RawResponse   *http.Response
	// OK
	TransferResponse *shared.TransferResponse
}

func (o *ConnectorsTransferResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ConnectorsTransferResponse) GetErrorResponse() *shared.ErrorResponse {
	if o == nil {
		return nil
	}
	return o.ErrorResponse
}

func (o *ConnectorsTransferResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ConnectorsTransferResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *ConnectorsTransferResponse) GetTransferResponse() *shared.TransferResponse {
	if o == nil {
		return nil
	}
	return o.TransferResponse
}
