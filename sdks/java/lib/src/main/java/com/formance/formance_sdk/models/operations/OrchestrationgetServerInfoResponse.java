/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;


public class OrchestrationgetServerInfoResponse {
    
    public String contentType;

    public OrchestrationgetServerInfoResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    /**
     * General error
     */
    
    public com.formance.formance_sdk.models.shared.Error error;

    public OrchestrationgetServerInfoResponse withError(com.formance.formance_sdk.models.shared.Error error) {
        this.error = error;
        return this;
    }
    
    /**
     * Server information
     */
    
    public com.formance.formance_sdk.models.shared.ServerInfo serverInfo;

    public OrchestrationgetServerInfoResponse withServerInfo(com.formance.formance_sdk.models.shared.ServerInfo serverInfo) {
        this.serverInfo = serverInfo;
        return this;
    }
    
    
    public Integer statusCode;

    public OrchestrationgetServerInfoResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public OrchestrationgetServerInfoResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    public OrchestrationgetServerInfoResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
