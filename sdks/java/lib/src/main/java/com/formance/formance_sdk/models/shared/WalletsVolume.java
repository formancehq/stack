/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

public class WalletsVolume {
    @JsonProperty("balance")
    public Long balance;

    public WalletsVolume withBalance(Long balance) {
        this.balance = balance;
        return this;
    }
    
    @JsonProperty("input")
    public Long input;

    public WalletsVolume withInput(Long input) {
        this.input = input;
        return this;
    }
    
    @JsonProperty("output")
    public Long output;

    public WalletsVolume withOutput(Long output) {
        this.output = output;
        return this;
    }
    
    public WalletsVolume(@JsonProperty("balance") Long balance, @JsonProperty("input") Long input, @JsonProperty("output") Long output) {
        this.balance = balance;
        this.input = input;
        this.output = output;
  }
}
