// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type ReadConnectorConfigRequest struct {
	// The name of the connector.
	Connector shared.Connector `pathParam:"style=simple,explode=false,name=connector"`
}

func (o *ReadConnectorConfigRequest) GetConnector() shared.Connector {
	if o == nil {
		return shared.Connector("")
	}
	return o.Connector
}

type ReadConnectorConfigResponse struct {
	// OK
	ConnectorConfigResponse *shared.ConnectorConfigResponse
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *ReadConnectorConfigResponse) GetConnectorConfigResponse() *shared.ConnectorConfigResponse {
	if o == nil {
		return nil
	}
	return o.ConnectorConfigResponse
}

func (o *ReadConnectorConfigResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ReadConnectorConfigResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ReadConnectorConfigResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
