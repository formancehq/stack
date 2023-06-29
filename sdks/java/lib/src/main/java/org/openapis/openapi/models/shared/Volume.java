/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

public class Volume {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("balance")
    public Long balance;

    public Volume withBalance(Long balance) {
        this.balance = balance;
        return this;
    }
    
    @JsonProperty("input")
    public Long input;

    public Volume withInput(Long input) {
        this.input = input;
        return this;
    }
    
    @JsonProperty("output")
    public Long output;

    public Volume withOutput(Long output) {
        this.output = output;
        return this;
    }
    
    public Volume(@JsonProperty("input") Long input, @JsonProperty("output") Long output) {
        this.input = input;
        this.output = output;
  }
}
