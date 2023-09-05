/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;


public class ActivateConfigResponse {
    /**
     * Config successfully activated.
     */
    
    public com.formance.formance_sdk.models.shared.ConfigResponse configResponse;

    public ActivateConfigResponse withConfigResponse(com.formance.formance_sdk.models.shared.ConfigResponse configResponse) {
        this.configResponse = configResponse;
        return this;
    }
    
    
    public String contentType;

    public ActivateConfigResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    
    public Integer statusCode;

    public ActivateConfigResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public ActivateConfigResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    /**
     * Error
     */
    
    public com.formance.formance_sdk.models.shared.WebhooksErrorResponse webhooksErrorResponse;

    public ActivateConfigResponse withWebhooksErrorResponse(com.formance.formance_sdk.models.shared.WebhooksErrorResponse webhooksErrorResponse) {
        this.webhooksErrorResponse = webhooksErrorResponse;
        return this;
    }
    
    public ActivateConfigResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
