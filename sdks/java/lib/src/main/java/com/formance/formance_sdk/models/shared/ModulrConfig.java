/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

public class ModulrConfig {
    @JsonProperty("apiKey")
    public String apiKey;

    public ModulrConfig withApiKey(String apiKey) {
        this.apiKey = apiKey;
        return this;
    }
    
    @JsonProperty("apiSecret")
    public String apiSecret;

    public ModulrConfig withApiSecret(String apiSecret) {
        this.apiSecret = apiSecret;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("endpoint")
    public String endpoint;

    public ModulrConfig withEndpoint(String endpoint) {
        this.endpoint = endpoint;
        return this;
    }
    
    @JsonProperty("name")
    public String name;

    public ModulrConfig withName(String name) {
        this.name = name;
        return this;
    }
    
    /**
     * The frequency at which the connector will try to fetch new BalanceTransaction objects from Modulr API.
     * 
     */
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("pollingPeriod")
    public String pollingPeriod;

    public ModulrConfig withPollingPeriod(String pollingPeriod) {
        this.pollingPeriod = pollingPeriod;
        return this;
    }
    
    public ModulrConfig(@JsonProperty("apiKey") String apiKey, @JsonProperty("apiSecret") String apiSecret, @JsonProperty("name") String name) {
        this.apiKey = apiKey;
        this.apiSecret = apiSecret;
        this.name = name;
  }
}
