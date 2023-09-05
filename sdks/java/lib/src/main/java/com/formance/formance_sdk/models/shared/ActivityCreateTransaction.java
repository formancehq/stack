/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;


public class ActivityCreateTransaction {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("data")
    public PostTransaction data;

    public ActivityCreateTransaction withData(PostTransaction data) {
        this.data = data;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("ledger")
    public String ledger;

    public ActivityCreateTransaction withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    public ActivityCreateTransaction(){}
}
