/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;

public class GetTransactionsResponse {
    
    public String contentType;

    public GetTransactionsResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    /**
     * OK
     */
    
    public org.openapis.openapi.models.shared.GetTransactionsResponse getTransactionsResponse;

    public GetTransactionsResponse withGetTransactionsResponse(org.openapis.openapi.models.shared.GetTransactionsResponse getTransactionsResponse) {
        this.getTransactionsResponse = getTransactionsResponse;
        return this;
    }
    
    
    public Integer statusCode;

    public GetTransactionsResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public GetTransactionsResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    /**
     * Error
     */
    
    public org.openapis.openapi.models.shared.WalletsErrorResponse walletsErrorResponse;

    public GetTransactionsResponse withWalletsErrorResponse(org.openapis.openapi.models.shared.WalletsErrorResponse walletsErrorResponse) {
        this.walletsErrorResponse = walletsErrorResponse;
        return this;
    }
    
    public GetTransactionsResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
