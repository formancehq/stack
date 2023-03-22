/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

public class BankingCircleConfig {
    @JsonProperty("authorizationEndpoint")
    public String authorizationEndpoint;

    public BankingCircleConfig withAuthorizationEndpoint(String authorizationEndpoint) {
        this.authorizationEndpoint = authorizationEndpoint;
        return this;
    }
    
    @JsonProperty("endpoint")
    public String endpoint;

    public BankingCircleConfig withEndpoint(String endpoint) {
        this.endpoint = endpoint;
        return this;
    }
    
    @JsonProperty("password")
    public String password;

    public BankingCircleConfig withPassword(String password) {
        this.password = password;
        return this;
    }
    
    @JsonProperty("username")
    public String username;

    public BankingCircleConfig withUsername(String username) {
        this.username = username;
        return this;
    }
    
    public BankingCircleConfig(@JsonProperty("username") String username, @JsonProperty("password") String password, @JsonProperty("authorizationEndpoint") String authorizationEndpoint, @JsonProperty("endpoint") String endpoint) {
        this.username = username;
        this.password = password;
        this.authorizationEndpoint = authorizationEndpoint;
        this.endpoint = endpoint;
  }
}
