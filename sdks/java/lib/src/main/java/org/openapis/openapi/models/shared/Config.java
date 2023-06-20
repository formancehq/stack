/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

public class Config {
    @JsonProperty("storage")
    public LedgerStorage storage;

    public Config withStorage(LedgerStorage storage) {
        this.storage = storage;
        return this;
    }
    
    public Config(@JsonProperty("storage") LedgerStorage storage) {
        this.storage = storage;
  }
}
