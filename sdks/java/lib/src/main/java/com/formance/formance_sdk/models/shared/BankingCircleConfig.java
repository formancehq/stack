/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
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
    
    /**
     * The frequency at which the connector will try to fetch new BalanceTransaction objects from Banking Circle API.
     * 
     */
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("pollingPeriod")
    public String pollingPeriod;

    public BankingCircleConfig withPollingPeriod(String pollingPeriod) {
        this.pollingPeriod = pollingPeriod;
        return this;
    }
    
    @JsonProperty("userCertificate")
    public String userCertificate;

    public BankingCircleConfig withUserCertificate(String userCertificate) {
        this.userCertificate = userCertificate;
        return this;
    }
    
    @JsonProperty("userCertificateKey")
    public String userCertificateKey;

    public BankingCircleConfig withUserCertificateKey(String userCertificateKey) {
        this.userCertificateKey = userCertificateKey;
        return this;
    }
    
    @JsonProperty("username")
    public String username;

    public BankingCircleConfig withUsername(String username) {
        this.username = username;
        return this;
    }
    
    public BankingCircleConfig(@JsonProperty("userCertificate") String userCertificate, @JsonProperty("password") String password, @JsonProperty("authorizationEndpoint") String authorizationEndpoint, @JsonProperty("endpoint") String endpoint) {
        this.userCertificate = userCertificate;
        this.password = password;
        this.authorizationEndpoint = authorizationEndpoint;
        this.endpoint = endpoint;
  }
}
