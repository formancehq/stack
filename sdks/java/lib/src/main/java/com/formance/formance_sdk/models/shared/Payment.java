/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import com.formance.formance_sdk.utils.DateTimeDeserializer;
import com.formance.formance_sdk.utils.DateTimeSerializer;
import java.time.OffsetDateTime;

public class Payment {
    @JsonProperty("accountID")
    public String accountID;

    public Payment withAccountID(String accountID) {
        this.accountID = accountID;
        return this;
    }
    
    @JsonProperty("adjustments")
    public PaymentAdjustment[] adjustments;

    public Payment withAdjustments(PaymentAdjustment[] adjustments) {
        this.adjustments = adjustments;
        return this;
    }
    
    @JsonProperty("asset")
    public String asset;

    public Payment withAsset(String asset) {
        this.asset = asset;
        return this;
    }
    
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("createdAt")
    public OffsetDateTime createdAt;

    public Payment withCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
        return this;
    }
    
    @JsonProperty("id")
    public String id;

    public Payment withId(String id) {
        this.id = id;
        return this;
    }
    
    @JsonProperty("initialAmount")
    public Long initialAmount;

    public Payment withInitialAmount(Long initialAmount) {
        this.initialAmount = initialAmount;
        return this;
    }
    
    @JsonProperty("metadata")
    public PaymentMetadata metadata;

    public Payment withMetadata(PaymentMetadata metadata) {
        this.metadata = metadata;
        return this;
    }
    
    @JsonProperty("provider")
    public Connector provider;

    public Payment withProvider(Connector provider) {
        this.provider = provider;
        return this;
    }
    
    @JsonProperty("raw")
    public java.util.Map<String, Object> raw;

    public Payment withRaw(java.util.Map<String, Object> raw) {
        this.raw = raw;
        return this;
    }
    
    @JsonProperty("reference")
    public String reference;

    public Payment withReference(String reference) {
        this.reference = reference;
        return this;
    }
    
    @JsonProperty("scheme")
    public PaymentScheme scheme;

    public Payment withScheme(PaymentScheme scheme) {
        this.scheme = scheme;
        return this;
    }
    
    @JsonProperty("status")
    public PaymentStatus status;

    public Payment withStatus(PaymentStatus status) {
        this.status = status;
        return this;
    }
    
    @JsonProperty("type")
    public PaymentType type;

    public Payment withType(PaymentType type) {
        this.type = type;
        return this;
    }
    
    public Payment(@JsonProperty("accountID") String accountID, @JsonProperty("adjustments") PaymentAdjustment[] adjustments, @JsonProperty("asset") String asset, @JsonProperty("createdAt") OffsetDateTime createdAt, @JsonProperty("id") String id, @JsonProperty("initialAmount") Long initialAmount, @JsonProperty("metadata") PaymentMetadata metadata, @JsonProperty("provider") Connector provider, @JsonProperty("raw") java.util.Map<String, Object> raw, @JsonProperty("reference") String reference, @JsonProperty("scheme") PaymentScheme scheme, @JsonProperty("status") PaymentStatus status, @JsonProperty("type") PaymentType type) {
        this.accountID = accountID;
        this.adjustments = adjustments;
        this.asset = asset;
        this.createdAt = createdAt;
        this.id = id;
        this.initialAmount = initialAmount;
        this.metadata = metadata;
        this.provider = provider;
        this.raw = raw;
        this.reference = reference;
        this.scheme = scheme;
        this.status = status;
        this.type = type;
  }
}
