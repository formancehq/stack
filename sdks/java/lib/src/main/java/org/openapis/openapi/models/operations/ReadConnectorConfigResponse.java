/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;

public class ReadConnectorConfigResponse {
    /**
     * OK
     */
    
    public org.openapis.openapi.models.shared.ConnectorConfigResponse connectorConfigResponse;

    public ReadConnectorConfigResponse withConnectorConfigResponse(org.openapis.openapi.models.shared.ConnectorConfigResponse connectorConfigResponse) {
        this.connectorConfigResponse = connectorConfigResponse;
        return this;
    }
    
    
    public String contentType;

    public ReadConnectorConfigResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    
    public Integer statusCode;

    public ReadConnectorConfigResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public ReadConnectorConfigResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    public ReadConnectorConfigResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
