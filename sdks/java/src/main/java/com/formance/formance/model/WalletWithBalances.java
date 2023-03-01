/*
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * The version of the OpenAPI document: v1.0.20230301
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.formance.formance.model;

import java.util.Objects;
import java.util.Arrays;
import com.formance.formance.model.WalletWithBalancesBalances;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;
import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

/**
 * WalletWithBalances
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class WalletWithBalances {
  public static final String SERIALIZED_NAME_ID = "id";
  @SerializedName(SERIALIZED_NAME_ID)
  private UUID id;

  public static final String SERIALIZED_NAME_METADATA = "metadata";
  @SerializedName(SERIALIZED_NAME_METADATA)
  private Map<String, Object> metadata = new HashMap<>();

  public static final String SERIALIZED_NAME_NAME = "name";
  @SerializedName(SERIALIZED_NAME_NAME)
  private String name;

  public static final String SERIALIZED_NAME_CREATED_AT = "createdAt";
  @SerializedName(SERIALIZED_NAME_CREATED_AT)
  private OffsetDateTime createdAt;

  public static final String SERIALIZED_NAME_BALANCES = "balances";
  @SerializedName(SERIALIZED_NAME_BALANCES)
  private WalletWithBalancesBalances balances;

  public static final String SERIALIZED_NAME_LEDGER = "ledger";
  @SerializedName(SERIALIZED_NAME_LEDGER)
  private String ledger;

  public WalletWithBalances() {
  }

  public WalletWithBalances id(UUID id) {
    
    this.id = id;
    return this;
  }

   /**
   * The unique ID of the wallet.
   * @return id
  **/
  @javax.annotation.Nonnull

  public UUID getId() {
    return id;
  }


  public void setId(UUID id) {
    this.id = id;
  }


  public WalletWithBalances metadata(Map<String, Object> metadata) {
    
    this.metadata = metadata;
    return this;
  }

  public WalletWithBalances putMetadataItem(String key, Object metadataItem) {
    this.metadata.put(key, metadataItem);
    return this;
  }

   /**
   * Metadata associated with the wallet.
   * @return metadata
  **/
  @javax.annotation.Nonnull

  public Map<String, Object> getMetadata() {
    return metadata;
  }


  public void setMetadata(Map<String, Object> metadata) {
    this.metadata = metadata;
  }


  public WalletWithBalances name(String name) {
    
    this.name = name;
    return this;
  }

   /**
   * Get name
   * @return name
  **/
  @javax.annotation.Nonnull

  public String getName() {
    return name;
  }


  public void setName(String name) {
    this.name = name;
  }


  public WalletWithBalances createdAt(OffsetDateTime createdAt) {
    
    this.createdAt = createdAt;
    return this;
  }

   /**
   * Get createdAt
   * @return createdAt
  **/
  @javax.annotation.Nonnull

  public OffsetDateTime getCreatedAt() {
    return createdAt;
  }


  public void setCreatedAt(OffsetDateTime createdAt) {
    this.createdAt = createdAt;
  }


  public WalletWithBalances balances(WalletWithBalancesBalances balances) {
    
    this.balances = balances;
    return this;
  }

   /**
   * Get balances
   * @return balances
  **/
  @javax.annotation.Nonnull

  public WalletWithBalancesBalances getBalances() {
    return balances;
  }


  public void setBalances(WalletWithBalancesBalances balances) {
    this.balances = balances;
  }


  public WalletWithBalances ledger(String ledger) {
    
    this.ledger = ledger;
    return this;
  }

   /**
   * Get ledger
   * @return ledger
  **/
  @javax.annotation.Nonnull

  public String getLedger() {
    return ledger;
  }


  public void setLedger(String ledger) {
    this.ledger = ledger;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    WalletWithBalances walletWithBalances = (WalletWithBalances) o;
    return Objects.equals(this.id, walletWithBalances.id) &&
        Objects.equals(this.metadata, walletWithBalances.metadata) &&
        Objects.equals(this.name, walletWithBalances.name) &&
        Objects.equals(this.createdAt, walletWithBalances.createdAt) &&
        Objects.equals(this.balances, walletWithBalances.balances) &&
        Objects.equals(this.ledger, walletWithBalances.ledger);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, metadata, name, createdAt, balances, ledger);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class WalletWithBalances {\n");
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    metadata: ").append(toIndentedString(metadata)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    createdAt: ").append(toIndentedString(createdAt)).append("\n");
    sb.append("    balances: ").append(toIndentedString(balances)).append("\n");
    sb.append("    ledger: ").append(toIndentedString(ledger)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

