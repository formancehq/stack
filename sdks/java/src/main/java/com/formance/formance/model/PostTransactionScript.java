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
 * PostTransactionScript
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2023-01-10T17:14:15.510105Z[Etc/UTC]")
public class PostTransactionScript {
  public static final String SERIALIZED_NAME_PLAIN = "plain";
  @SerializedName(SERIALIZED_NAME_PLAIN)
  private String plain;

  public static final String SERIALIZED_NAME_VARS = "vars";
  @SerializedName(SERIALIZED_NAME_VARS)
  private Object vars;

  public PostTransactionScript() {
  }

  public PostTransactionScript plain(String plain) {
    
    this.plain = plain;
    return this;
  }

   /**
   * Get plain
   * @return plain
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(example = "vars { account $user } send [COIN 10] (  source = @world  destination = $user ) ", required = true, value = "")

  public String getPlain() {
    return plain;
  }


  public void setPlain(String plain) {
    this.plain = plain;
  }


  public PostTransactionScript vars(Object vars) {
    
    this.vars = vars;
    return this;
  }

   /**
   * Get vars
   * @return vars
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(example = "{\"vars\":{\"user\":\"users:042\"}}", value = "")

  public Object getVars() {
    return vars;
  }


  public void setVars(Object vars) {
    this.vars = vars;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    PostTransactionScript postTransactionScript = (PostTransactionScript) o;
    return Objects.equals(this.plain, postTransactionScript.plain) &&
        Objects.equals(this.vars, postTransactionScript.vars);
  }

  @Override
  public int hashCode() {
    return Objects.hash(plain, vars);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class PostTransactionScript {\n");
    sb.append("    plain: ").append(toIndentedString(plain)).append("\n");
    sb.append("    vars: ").append(toIndentedString(vars)).append("\n");
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

