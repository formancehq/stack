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
import com.formance.formance.model.PaymentStatus;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.time.OffsetDateTime;

/**
 * PaymentAdjustment
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class PaymentAdjustment {
  public static final String SERIALIZED_NAME_STATUS = "status";
  @SerializedName(SERIALIZED_NAME_STATUS)
  private PaymentStatus status;

  public static final String SERIALIZED_NAME_AMOUNT = "amount";
  @SerializedName(SERIALIZED_NAME_AMOUNT)
  private Long amount;

  public static final String SERIALIZED_NAME_DATE = "date";
  @SerializedName(SERIALIZED_NAME_DATE)
  private OffsetDateTime date;

  public static final String SERIALIZED_NAME_RAW = "raw";
  @SerializedName(SERIALIZED_NAME_RAW)
  private Object raw;

  public static final String SERIALIZED_NAME_ABSOLUTE = "absolute";
  @SerializedName(SERIALIZED_NAME_ABSOLUTE)
  private Boolean absolute;

  public PaymentAdjustment() {
  }

  public PaymentAdjustment status(PaymentStatus status) {
    
    this.status = status;
    return this;
  }

   /**
   * Get status
   * @return status
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public PaymentStatus getStatus() {
    return status;
  }


  public void setStatus(PaymentStatus status) {
    this.status = status;
  }


  public PaymentAdjustment amount(Long amount) {
    
    this.amount = amount;
    return this;
  }

   /**
   * Get amount
   * minimum: 0
   * @return amount
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(example = "100", required = true, value = "")

  public Long getAmount() {
    return amount;
  }


  public void setAmount(Long amount) {
    this.amount = amount;
  }


  public PaymentAdjustment date(OffsetDateTime date) {
    
    this.date = date;
    return this;
  }

   /**
   * Get date
   * @return date
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public OffsetDateTime getDate() {
    return date;
  }


  public void setDate(OffsetDateTime date) {
    this.date = date;
  }


  public PaymentAdjustment raw(Object raw) {
    
    this.raw = raw;
    return this;
  }

   /**
   * Get raw
   * @return raw
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Object getRaw() {
    return raw;
  }


  public void setRaw(Object raw) {
    this.raw = raw;
  }


  public PaymentAdjustment absolute(Boolean absolute) {
    
    this.absolute = absolute;
    return this;
  }

   /**
   * Get absolute
   * @return absolute
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Boolean getAbsolute() {
    return absolute;
  }


  public void setAbsolute(Boolean absolute) {
    this.absolute = absolute;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    PaymentAdjustment paymentAdjustment = (PaymentAdjustment) o;
    return Objects.equals(this.status, paymentAdjustment.status) &&
        Objects.equals(this.amount, paymentAdjustment.amount) &&
        Objects.equals(this.date, paymentAdjustment.date) &&
        Objects.equals(this.raw, paymentAdjustment.raw) &&
        Objects.equals(this.absolute, paymentAdjustment.absolute);
  }

  @Override
  public int hashCode() {
    return Objects.hash(status, amount, date, raw, absolute);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class PaymentAdjustment {\n");
    sb.append("    status: ").append(toIndentedString(status)).append("\n");
    sb.append("    amount: ").append(toIndentedString(amount)).append("\n");
    sb.append("    date: ").append(toIndentedString(date)).append("\n");
    sb.append("    raw: ").append(toIndentedString(raw)).append("\n");
    sb.append("    absolute: ").append(toIndentedString(absolute)).append("\n");
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

