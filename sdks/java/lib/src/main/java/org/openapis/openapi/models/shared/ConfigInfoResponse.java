/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * ConfigInfoResponse - OK
 */
public class ConfigInfoResponse {
    @JsonProperty("config")
    public Config config;

    public ConfigInfoResponse withConfig(Config config) {
        this.config = config;
        return this;
    }
    
    @JsonProperty("server")
    public String server;

    public ConfigInfoResponse withServer(String server) {
        this.server = server;
        return this;
    }
    
    @JsonProperty("version")
    public String version;

    public ConfigInfoResponse withVersion(String version) {
        this.version = version;
        return this;
    }
    
    public ConfigInfoResponse(@JsonProperty("config") Config config, @JsonProperty("server") String server, @JsonProperty("version") String version) {
        this.config = config;
        this.server = server;
        this.version = version;
  }
}
