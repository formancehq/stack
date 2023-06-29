/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import org.openapis.openapi.utils.SpeakeasyMetadata;

public class GetBalancesAggregatedRequest {
    /**
     * Filter balances involving given account, either as source or destination.
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=address")
    public String address;

    public GetBalancesAggregatedRequest withAddress(String address) {
        this.address = address;
        return this;
    }
    
    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=ledger")
    public String ledger;

    public GetBalancesAggregatedRequest withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    public GetBalancesAggregatedRequest(@JsonProperty("ledger") String ledger) {
        this.ledger = ledger;
  }
}
