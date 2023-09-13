/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * AccountResponse - OK
 */
public class AccountResponse {
    @JsonProperty("data")
    public AccountWithVolumesAndBalances data;

    public AccountResponse withData(AccountWithVolumesAndBalances data) {
        this.data = data;
        return this;
    }
    
    public AccountResponse(@JsonProperty("data") AccountWithVolumesAndBalances data) {
        this.data = data;
  }
}
