/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;


public class ActivityGetWallet {
    @JsonProperty("id")
    public String id;

    public ActivityGetWallet withId(String id) {
        this.id = id;
        return this;
    }
    
    public ActivityGetWallet(@JsonProperty("id") String id) {
        this.id = id;
  }
}
