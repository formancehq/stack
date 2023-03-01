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
import com.formance.formance.model.ActivityConfirmHold;
import com.formance.formance.model.ActivityCreateTransaction;
import com.formance.formance.model.ActivityCreditWallet;
import com.formance.formance.model.ActivityDebitWallet;
import com.formance.formance.model.ActivityGetAccount;
import com.formance.formance.model.ActivityGetPayment;
import com.formance.formance.model.ActivityGetWallet;
import com.formance.formance.model.ActivityRevertTransaction;
import com.formance.formance.model.ActivityVoidHold;
import com.formance.formance.model.StripeTransferRequest;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;

/**
 * WorkflowInstanceHistoryStageInput
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class WorkflowInstanceHistoryStageInput {
  public static final String SERIALIZED_NAME_GET_ACCOUNT = "GetAccount";
  @SerializedName(SERIALIZED_NAME_GET_ACCOUNT)
  private ActivityGetAccount getAccount;

  public static final String SERIALIZED_NAME_CREATE_TRANSACTION = "CreateTransaction";
  @SerializedName(SERIALIZED_NAME_CREATE_TRANSACTION)
  private ActivityCreateTransaction createTransaction;

  public static final String SERIALIZED_NAME_REVERT_TRANSACTION = "RevertTransaction";
  @SerializedName(SERIALIZED_NAME_REVERT_TRANSACTION)
  private ActivityRevertTransaction revertTransaction;

  public static final String SERIALIZED_NAME_STRIPE_TRANSFER = "StripeTransfer";
  @SerializedName(SERIALIZED_NAME_STRIPE_TRANSFER)
  private StripeTransferRequest stripeTransfer;

  public static final String SERIALIZED_NAME_GET_PAYMENT = "GetPayment";
  @SerializedName(SERIALIZED_NAME_GET_PAYMENT)
  private ActivityGetPayment getPayment;

  public static final String SERIALIZED_NAME_CONFIRM_HOLD = "ConfirmHold";
  @SerializedName(SERIALIZED_NAME_CONFIRM_HOLD)
  private ActivityConfirmHold confirmHold;

  public static final String SERIALIZED_NAME_CREDIT_WALLET = "CreditWallet";
  @SerializedName(SERIALIZED_NAME_CREDIT_WALLET)
  private ActivityCreditWallet creditWallet;

  public static final String SERIALIZED_NAME_DEBIT_WALLET = "DebitWallet";
  @SerializedName(SERIALIZED_NAME_DEBIT_WALLET)
  private ActivityDebitWallet debitWallet;

  public static final String SERIALIZED_NAME_GET_WALLET = "GetWallet";
  @SerializedName(SERIALIZED_NAME_GET_WALLET)
  private ActivityGetWallet getWallet;

  public static final String SERIALIZED_NAME_VOID_HOLD = "VoidHold";
  @SerializedName(SERIALIZED_NAME_VOID_HOLD)
  private ActivityVoidHold voidHold;

  public WorkflowInstanceHistoryStageInput() {
  }

  public WorkflowInstanceHistoryStageInput getAccount(ActivityGetAccount getAccount) {
    
    this.getAccount = getAccount;
    return this;
  }

   /**
   * Get getAccount
   * @return getAccount
  **/
  @javax.annotation.Nullable

  public ActivityGetAccount getGetAccount() {
    return getAccount;
  }


  public void setGetAccount(ActivityGetAccount getAccount) {
    this.getAccount = getAccount;
  }


  public WorkflowInstanceHistoryStageInput createTransaction(ActivityCreateTransaction createTransaction) {
    
    this.createTransaction = createTransaction;
    return this;
  }

   /**
   * Get createTransaction
   * @return createTransaction
  **/
  @javax.annotation.Nullable

  public ActivityCreateTransaction getCreateTransaction() {
    return createTransaction;
  }


  public void setCreateTransaction(ActivityCreateTransaction createTransaction) {
    this.createTransaction = createTransaction;
  }


  public WorkflowInstanceHistoryStageInput revertTransaction(ActivityRevertTransaction revertTransaction) {
    
    this.revertTransaction = revertTransaction;
    return this;
  }

   /**
   * Get revertTransaction
   * @return revertTransaction
  **/
  @javax.annotation.Nullable

  public ActivityRevertTransaction getRevertTransaction() {
    return revertTransaction;
  }


  public void setRevertTransaction(ActivityRevertTransaction revertTransaction) {
    this.revertTransaction = revertTransaction;
  }


  public WorkflowInstanceHistoryStageInput stripeTransfer(StripeTransferRequest stripeTransfer) {
    
    this.stripeTransfer = stripeTransfer;
    return this;
  }

   /**
   * Get stripeTransfer
   * @return stripeTransfer
  **/
  @javax.annotation.Nullable

  public StripeTransferRequest getStripeTransfer() {
    return stripeTransfer;
  }


  public void setStripeTransfer(StripeTransferRequest stripeTransfer) {
    this.stripeTransfer = stripeTransfer;
  }


  public WorkflowInstanceHistoryStageInput getPayment(ActivityGetPayment getPayment) {
    
    this.getPayment = getPayment;
    return this;
  }

   /**
   * Get getPayment
   * @return getPayment
  **/
  @javax.annotation.Nullable

  public ActivityGetPayment getGetPayment() {
    return getPayment;
  }


  public void setGetPayment(ActivityGetPayment getPayment) {
    this.getPayment = getPayment;
  }


  public WorkflowInstanceHistoryStageInput confirmHold(ActivityConfirmHold confirmHold) {
    
    this.confirmHold = confirmHold;
    return this;
  }

   /**
   * Get confirmHold
   * @return confirmHold
  **/
  @javax.annotation.Nullable

  public ActivityConfirmHold getConfirmHold() {
    return confirmHold;
  }


  public void setConfirmHold(ActivityConfirmHold confirmHold) {
    this.confirmHold = confirmHold;
  }


  public WorkflowInstanceHistoryStageInput creditWallet(ActivityCreditWallet creditWallet) {
    
    this.creditWallet = creditWallet;
    return this;
  }

   /**
   * Get creditWallet
   * @return creditWallet
  **/
  @javax.annotation.Nullable

  public ActivityCreditWallet getCreditWallet() {
    return creditWallet;
  }


  public void setCreditWallet(ActivityCreditWallet creditWallet) {
    this.creditWallet = creditWallet;
  }


  public WorkflowInstanceHistoryStageInput debitWallet(ActivityDebitWallet debitWallet) {
    
    this.debitWallet = debitWallet;
    return this;
  }

   /**
   * Get debitWallet
   * @return debitWallet
  **/
  @javax.annotation.Nullable

  public ActivityDebitWallet getDebitWallet() {
    return debitWallet;
  }


  public void setDebitWallet(ActivityDebitWallet debitWallet) {
    this.debitWallet = debitWallet;
  }


  public WorkflowInstanceHistoryStageInput getWallet(ActivityGetWallet getWallet) {
    
    this.getWallet = getWallet;
    return this;
  }

   /**
   * Get getWallet
   * @return getWallet
  **/
  @javax.annotation.Nullable

  public ActivityGetWallet getGetWallet() {
    return getWallet;
  }


  public void setGetWallet(ActivityGetWallet getWallet) {
    this.getWallet = getWallet;
  }


  public WorkflowInstanceHistoryStageInput voidHold(ActivityVoidHold voidHold) {
    
    this.voidHold = voidHold;
    return this;
  }

   /**
   * Get voidHold
   * @return voidHold
  **/
  @javax.annotation.Nullable

  public ActivityVoidHold getVoidHold() {
    return voidHold;
  }


  public void setVoidHold(ActivityVoidHold voidHold) {
    this.voidHold = voidHold;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    WorkflowInstanceHistoryStageInput workflowInstanceHistoryStageInput = (WorkflowInstanceHistoryStageInput) o;
    return Objects.equals(this.getAccount, workflowInstanceHistoryStageInput.getAccount) &&
        Objects.equals(this.createTransaction, workflowInstanceHistoryStageInput.createTransaction) &&
        Objects.equals(this.revertTransaction, workflowInstanceHistoryStageInput.revertTransaction) &&
        Objects.equals(this.stripeTransfer, workflowInstanceHistoryStageInput.stripeTransfer) &&
        Objects.equals(this.getPayment, workflowInstanceHistoryStageInput.getPayment) &&
        Objects.equals(this.confirmHold, workflowInstanceHistoryStageInput.confirmHold) &&
        Objects.equals(this.creditWallet, workflowInstanceHistoryStageInput.creditWallet) &&
        Objects.equals(this.debitWallet, workflowInstanceHistoryStageInput.debitWallet) &&
        Objects.equals(this.getWallet, workflowInstanceHistoryStageInput.getWallet) &&
        Objects.equals(this.voidHold, workflowInstanceHistoryStageInput.voidHold);
  }

  @Override
  public int hashCode() {
    return Objects.hash(getAccount, createTransaction, revertTransaction, stripeTransfer, getPayment, confirmHold, creditWallet, debitWallet, getWallet, voidHold);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class WorkflowInstanceHistoryStageInput {\n");
    sb.append("    getAccount: ").append(toIndentedString(getAccount)).append("\n");
    sb.append("    createTransaction: ").append(toIndentedString(createTransaction)).append("\n");
    sb.append("    revertTransaction: ").append(toIndentedString(revertTransaction)).append("\n");
    sb.append("    stripeTransfer: ").append(toIndentedString(stripeTransfer)).append("\n");
    sb.append("    getPayment: ").append(toIndentedString(getPayment)).append("\n");
    sb.append("    confirmHold: ").append(toIndentedString(confirmHold)).append("\n");
    sb.append("    creditWallet: ").append(toIndentedString(creditWallet)).append("\n");
    sb.append("    debitWallet: ").append(toIndentedString(debitWallet)).append("\n");
    sb.append("    getWallet: ").append(toIndentedString(getWallet)).append("\n");
    sb.append("    voidHold: ").append(toIndentedString(voidHold)).append("\n");
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

