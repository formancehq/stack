/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;


public class StageSendSource {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("account")
    public StageSendSourceAccount account;

    public StageSendSource withAccount(StageSendSourceAccount account) {
        this.account = account;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("payment")
    public StageSendSourcePayment payment;

    public StageSendSource withPayment(StageSendSourcePayment payment) {
        this.payment = payment;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("wallet")
    public StageSendSourceWallet wallet;

    public StageSendSource withWallet(StageSendSourceWallet wallet) {
        this.wallet = wallet;
        return this;
    }
    
    public StageSendSource(){}
}
