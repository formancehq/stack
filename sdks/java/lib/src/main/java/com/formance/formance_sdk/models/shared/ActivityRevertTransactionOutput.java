/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

public class ActivityRevertTransactionOutput {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("data")
    public OrchestrationTransaction data;

    public ActivityRevertTransactionOutput withData(OrchestrationTransaction data) {
        this.data = data;
        return this;
    }
    
    public ActivityRevertTransactionOutput(){}
}
