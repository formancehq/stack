/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;

public class GetTransactionRequest {
    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=ledger")
    public String ledger;

    public GetTransactionRequest withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    /**
     * Transaction ID.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=txid")
    public Long txid;

    public GetTransactionRequest withTxid(Long txid) {
        this.txid = txid;
        return this;
    }
    
    public GetTransactionRequest(@JsonProperty("ledger") String ledger, @JsonProperty("txid") Long txid) {
        this.ledger = ledger;
        this.txid = txid;
  }
}
