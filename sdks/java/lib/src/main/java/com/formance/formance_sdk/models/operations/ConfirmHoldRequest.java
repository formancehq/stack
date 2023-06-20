/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;

public class ConfirmHoldRequest {
    @SpeakeasyMetadata("request:mediaType=application/json")
    public com.formance.formance_sdk.models.shared.ConfirmHoldRequest confirmHoldRequest;

    public ConfirmHoldRequest withConfirmHoldRequest(com.formance.formance_sdk.models.shared.ConfirmHoldRequest confirmHoldRequest) {
        this.confirmHoldRequest = confirmHoldRequest;
        return this;
    }
    
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=hold_id")
    public String holdId;

    public ConfirmHoldRequest withHoldId(String holdId) {
        this.holdId = holdId;
        return this;
    }
    
    public ConfirmHoldRequest(@JsonProperty("hold_id") String holdId) {
        this.holdId = holdId;
  }
}
