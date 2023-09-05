/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Error - General error
 */

public class Error {
    @JsonProperty("errorCode")
    public ErrorErrorCode errorCode;

    public Error withErrorCode(ErrorErrorCode errorCode) {
        this.errorCode = errorCode;
        return this;
    }
    
    @JsonProperty("errorMessage")
    public String errorMessage;

    public Error withErrorMessage(String errorMessage) {
        this.errorMessage = errorMessage;
        return this;
    }
    
    public Error(@JsonProperty("errorCode") ErrorErrorCode errorCode, @JsonProperty("errorMessage") String errorMessage) {
        this.errorCode = errorCode;
        this.errorMessage = errorMessage;
  }
}
