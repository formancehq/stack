/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;
import java.time.OffsetDateTime;


public class ListTransactionsRequest {
    @SpeakeasyMetadata("request:mediaType=application/json")
    public java.util.Map<String, Object> requestBody;

    public ListTransactionsRequest withRequestBody(java.util.Map<String, Object> requestBody) {
        this.requestBody = requestBody;
        return this;
    }
    
    /**
     * Parameter used in pagination requests. Maximum page size is set to 15.
     * Set to the value of next for the next page of results.
     * Set to the value of previous for the previous page of results.
     * No other parameters can be set when this parameter is set.
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=cursor")
    public String cursor;

    public ListTransactionsRequest withCursor(String cursor) {
        this.cursor = cursor;
        return this;
    }
    
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=expand")
    public String expand;

    public ListTransactionsRequest withExpand(String expand) {
        this.expand = expand;
        return this;
    }
    
    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=ledger")
    public String ledger;

    public ListTransactionsRequest withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    /**
     * The maximum number of results to return per page.
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=pageSize")
    public Long pageSize;

    public ListTransactionsRequest withPageSize(Long pageSize) {
        this.pageSize = pageSize;
        return this;
    }
    
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=pit")
    public OffsetDateTime pit;

    public ListTransactionsRequest withPit(OffsetDateTime pit) {
        this.pit = pit;
        return this;
    }
    
    public ListTransactionsRequest(@JsonProperty("ledger") String ledger) {
        this.ledger = ledger;
  }
}
