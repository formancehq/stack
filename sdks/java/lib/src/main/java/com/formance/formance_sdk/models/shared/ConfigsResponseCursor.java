/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;


public class ConfigsResponseCursor {
    @JsonProperty("data")
    public WebhooksConfig[] data;

    public ConfigsResponseCursor withData(WebhooksConfig[] data) {
        this.data = data;
        return this;
    }
    
    @JsonProperty("hasMore")
    public Boolean hasMore;

    public ConfigsResponseCursor withHasMore(Boolean hasMore) {
        this.hasMore = hasMore;
        return this;
    }
    
    public ConfigsResponseCursor(@JsonProperty("data") WebhooksConfig[] data, @JsonProperty("hasMore") Boolean hasMore) {
        this.data = data;
        this.hasMore = hasMore;
  }
}
