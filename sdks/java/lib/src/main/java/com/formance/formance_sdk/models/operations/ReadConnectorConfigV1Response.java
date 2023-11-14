/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;

public class ReadConnectorConfigV1Response {
    /**
     * OK
     */
    
    public com.formance.formance_sdk.models.shared.ConnectorConfigResponse connectorConfigResponse;

    public ReadConnectorConfigV1Response withConnectorConfigResponse(com.formance.formance_sdk.models.shared.ConnectorConfigResponse connectorConfigResponse) {
        this.connectorConfigResponse = connectorConfigResponse;
        return this;
    }
    
    
    public String contentType;

    public ReadConnectorConfigV1Response withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    
    public Integer statusCode;

    public ReadConnectorConfigV1Response withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public ReadConnectorConfigV1Response withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    public ReadConnectorConfigV1Response(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
