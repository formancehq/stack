/*
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * The version of the OpenAPI document: develop
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.formance.formance.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

/**
 * AccountWithVolumesAndBalances
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2023-01-10T18:47:48.920061Z[Etc/UTC]")
public class AccountWithVolumesAndBalances {
  public static final String SERIALIZED_NAME_ADDRESS = "address";
  @SerializedName(SERIALIZED_NAME_ADDRESS)
  private String address;

  public static final String SERIALIZED_NAME_TYPE = "type";
  @SerializedName(SERIALIZED_NAME_TYPE)
  private String type;

  public static final String SERIALIZED_NAME_METADATA = "metadata";
  @SerializedName(SERIALIZED_NAME_METADATA)
  private Object metadata;

  public static final String SERIALIZED_NAME_VOLUMES = "volumes";
  @SerializedName(SERIALIZED_NAME_VOLUMES)
  private Map<String, Map<String, Long>> volumes = null;

  public static final String SERIALIZED_NAME_BALANCES = "balances";
  @SerializedName(SERIALIZED_NAME_BALANCES)
  private Map<String, Long> balances = null;

  public AccountWithVolumesAndBalances() {
  }

  public AccountWithVolumesAndBalances address(String address) {
    
    this.address = address;
    return this;
  }

   /**
   * Get address
   * @return address
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(example = "users:001", required = true, value = "")

  public String getAddress() {
    return address;
  }


  public void setAddress(String address) {
    this.address = address;
  }


  public AccountWithVolumesAndBalances type(String type) {
    
    this.type = type;
    return this;
  }

   /**
   * Get type
   * @return type
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(example = "virtual", value = "")

  public String getType() {
    return type;
  }


  public void setType(String type) {
    this.type = type;
  }


  public AccountWithVolumesAndBalances metadata(Object metadata) {
    
    this.metadata = metadata;
    return this;
  }

   /**
   * Get metadata
   * @return metadata
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(example = "{\"admin\":true,\"a\":{\"nested\":{\"key\":\"value\"}}}", value = "")

  public Object getMetadata() {
    return metadata;
  }


  public void setMetadata(Object metadata) {
    this.metadata = metadata;
  }


  public AccountWithVolumesAndBalances volumes(Map<String, Map<String, Long>> volumes) {
    
    this.volumes = volumes;
    return this;
  }

  public AccountWithVolumesAndBalances putVolumesItem(String key, Map<String, Long> volumesItem) {
    if (this.volumes == null) {
      this.volumes = new HashMap<>();
    }
    this.volumes.put(key, volumesItem);
    return this;
  }

   /**
   * Get volumes
   * @return volumes
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(example = "{\"COIN\":{\"input\":100,\"output\":0}}", value = "")

  public Map<String, Map<String, Long>> getVolumes() {
    return volumes;
  }


  public void setVolumes(Map<String, Map<String, Long>> volumes) {
    this.volumes = volumes;
  }


  public AccountWithVolumesAndBalances balances(Map<String, Long> balances) {
    
    this.balances = balances;
    return this;
  }

  public AccountWithVolumesAndBalances putBalancesItem(String key, Long balancesItem) {
    if (this.balances == null) {
      this.balances = new HashMap<>();
    }
    this.balances.put(key, balancesItem);
    return this;
  }

   /**
   * Get balances
   * @return balances
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(example = "{\"COIN\":100}", value = "")

  public Map<String, Long> getBalances() {
    return balances;
  }


  public void setBalances(Map<String, Long> balances) {
    this.balances = balances;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    AccountWithVolumesAndBalances accountWithVolumesAndBalances = (AccountWithVolumesAndBalances) o;
    return Objects.equals(this.address, accountWithVolumesAndBalances.address) &&
        Objects.equals(this.type, accountWithVolumesAndBalances.type) &&
        Objects.equals(this.metadata, accountWithVolumesAndBalances.metadata) &&
        Objects.equals(this.volumes, accountWithVolumesAndBalances.volumes) &&
        Objects.equals(this.balances, accountWithVolumesAndBalances.balances);
  }

  @Override
  public int hashCode() {
    return Objects.hash(address, type, metadata, volumes, balances);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AccountWithVolumesAndBalances {\n");
    sb.append("    address: ").append(toIndentedString(address)).append("\n");
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    metadata: ").append(toIndentedString(metadata)).append("\n");
    sb.append("    volumes: ").append(toIndentedString(volumes)).append("\n");
    sb.append("    balances: ").append(toIndentedString(balances)).append("\n");
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

