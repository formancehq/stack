/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;


public class DeleteTransactionMetadataRequest {
    /**
     * Transaction ID.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=id")
    public Long id;

    public DeleteTransactionMetadataRequest withId(Long id) {
        this.id = id;
        return this;
    }
    
    /**
     * The key to remove.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=key")
    public String key;

    public DeleteTransactionMetadataRequest withKey(String key) {
        this.key = key;
        return this;
    }
    
    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=ledger")
    public String ledger;

    public DeleteTransactionMetadataRequest withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    public DeleteTransactionMetadataRequest(@JsonProperty("id") Long id, @JsonProperty("key") String key, @JsonProperty("ledger") String ledger) {
        this.id = id;
        this.key = key;
        this.ledger = ledger;
  }
}
