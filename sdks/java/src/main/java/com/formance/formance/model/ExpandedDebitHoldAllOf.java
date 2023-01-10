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

/**
 * ExpandedDebitHoldAllOf
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2023-01-10T18:47:48.987377Z[Etc/UTC]")
public class ExpandedDebitHoldAllOf {
  public static final String SERIALIZED_NAME_REMAINING = "remaining";
  @SerializedName(SERIALIZED_NAME_REMAINING)
  private Long remaining;

  public static final String SERIALIZED_NAME_ORIGINAL_AMOUNT = "originalAmount";
  @SerializedName(SERIALIZED_NAME_ORIGINAL_AMOUNT)
  private Long originalAmount;

  public ExpandedDebitHoldAllOf() {
  }

  public ExpandedDebitHoldAllOf remaining(Long remaining) {
    
    this.remaining = remaining;
    return this;
  }

   /**
   * Remaining amount on hold
   * @return remaining
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(example = "10", required = true, value = "Remaining amount on hold")

  public Long getRemaining() {
    return remaining;
  }


  public void setRemaining(Long remaining) {
    this.remaining = remaining;
  }


  public ExpandedDebitHoldAllOf originalAmount(Long originalAmount) {
    
    this.originalAmount = originalAmount;
    return this;
  }

   /**
   * Original amount on hold
   * @return originalAmount
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(example = "100", required = true, value = "Original amount on hold")

  public Long getOriginalAmount() {
    return originalAmount;
  }


  public void setOriginalAmount(Long originalAmount) {
    this.originalAmount = originalAmount;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ExpandedDebitHoldAllOf expandedDebitHoldAllOf = (ExpandedDebitHoldAllOf) o;
    return Objects.equals(this.remaining, expandedDebitHoldAllOf.remaining) &&
        Objects.equals(this.originalAmount, expandedDebitHoldAllOf.originalAmount);
  }

  @Override
  public int hashCode() {
    return Objects.hash(remaining, originalAmount);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ExpandedDebitHoldAllOf {\n");
    sb.append("    remaining: ").append(toIndentedString(remaining)).append("\n");
    sb.append("    originalAmount: ").append(toIndentedString(originalAmount)).append("\n");
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

