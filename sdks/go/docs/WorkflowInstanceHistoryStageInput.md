# WorkflowInstanceHistoryStageInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GetAccount** | Pointer to [**ActivityGetAccount**](ActivityGetAccount.md) |  | [optional] 
**CreateTransaction** | Pointer to [**ActivityCreateTransaction**](ActivityCreateTransaction.md) |  | [optional] 
**RevertTransaction** | Pointer to [**ActivityRevertTransaction**](ActivityRevertTransaction.md) |  | [optional] 
**StripeTransfer** | Pointer to [**StripeTransferRequest**](StripeTransferRequest.md) |  | [optional] 
**GetPayment** | Pointer to [**ActivityGetPayment**](ActivityGetPayment.md) |  | [optional] 
**ConfirmHold** | Pointer to [**ActivityConfirmHold**](ActivityConfirmHold.md) |  | [optional] 
**CreditWallet** | Pointer to [**ActivityCreditWallet**](ActivityCreditWallet.md) |  | [optional] 
**DebitWallet** | Pointer to [**ActivityDebitWallet**](ActivityDebitWallet.md) |  | [optional] 
**GetWallet** | Pointer to [**ActivityGetWallet**](ActivityGetWallet.md) |  | [optional] 
**VoidHold** | Pointer to [**ActivityVoidHold**](ActivityVoidHold.md) |  | [optional] 

## Methods

### NewWorkflowInstanceHistoryStageInput

`func NewWorkflowInstanceHistoryStageInput() *WorkflowInstanceHistoryStageInput`

NewWorkflowInstanceHistoryStageInput instantiates a new WorkflowInstanceHistoryStageInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowInstanceHistoryStageInputWithDefaults

`func NewWorkflowInstanceHistoryStageInputWithDefaults() *WorkflowInstanceHistoryStageInput`

NewWorkflowInstanceHistoryStageInputWithDefaults instantiates a new WorkflowInstanceHistoryStageInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGetAccount

`func (o *WorkflowInstanceHistoryStageInput) GetGetAccount() ActivityGetAccount`

GetGetAccount returns the GetAccount field if non-nil, zero value otherwise.

### GetGetAccountOk

`func (o *WorkflowInstanceHistoryStageInput) GetGetAccountOk() (*ActivityGetAccount, bool)`

GetGetAccountOk returns a tuple with the GetAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGetAccount

`func (o *WorkflowInstanceHistoryStageInput) SetGetAccount(v ActivityGetAccount)`

SetGetAccount sets GetAccount field to given value.

### HasGetAccount

`func (o *WorkflowInstanceHistoryStageInput) HasGetAccount() bool`

HasGetAccount returns a boolean if a field has been set.

### GetCreateTransaction

`func (o *WorkflowInstanceHistoryStageInput) GetCreateTransaction() ActivityCreateTransaction`

GetCreateTransaction returns the CreateTransaction field if non-nil, zero value otherwise.

### GetCreateTransactionOk

`func (o *WorkflowInstanceHistoryStageInput) GetCreateTransactionOk() (*ActivityCreateTransaction, bool)`

GetCreateTransactionOk returns a tuple with the CreateTransaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreateTransaction

`func (o *WorkflowInstanceHistoryStageInput) SetCreateTransaction(v ActivityCreateTransaction)`

SetCreateTransaction sets CreateTransaction field to given value.

### HasCreateTransaction

`func (o *WorkflowInstanceHistoryStageInput) HasCreateTransaction() bool`

HasCreateTransaction returns a boolean if a field has been set.

### GetRevertTransaction

`func (o *WorkflowInstanceHistoryStageInput) GetRevertTransaction() ActivityRevertTransaction`

GetRevertTransaction returns the RevertTransaction field if non-nil, zero value otherwise.

### GetRevertTransactionOk

`func (o *WorkflowInstanceHistoryStageInput) GetRevertTransactionOk() (*ActivityRevertTransaction, bool)`

GetRevertTransactionOk returns a tuple with the RevertTransaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevertTransaction

`func (o *WorkflowInstanceHistoryStageInput) SetRevertTransaction(v ActivityRevertTransaction)`

SetRevertTransaction sets RevertTransaction field to given value.

### HasRevertTransaction

`func (o *WorkflowInstanceHistoryStageInput) HasRevertTransaction() bool`

HasRevertTransaction returns a boolean if a field has been set.

### GetStripeTransfer

`func (o *WorkflowInstanceHistoryStageInput) GetStripeTransfer() StripeTransferRequest`

GetStripeTransfer returns the StripeTransfer field if non-nil, zero value otherwise.

### GetStripeTransferOk

`func (o *WorkflowInstanceHistoryStageInput) GetStripeTransferOk() (*StripeTransferRequest, bool)`

GetStripeTransferOk returns a tuple with the StripeTransfer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripeTransfer

`func (o *WorkflowInstanceHistoryStageInput) SetStripeTransfer(v StripeTransferRequest)`

SetStripeTransfer sets StripeTransfer field to given value.

### HasStripeTransfer

`func (o *WorkflowInstanceHistoryStageInput) HasStripeTransfer() bool`

HasStripeTransfer returns a boolean if a field has been set.

### GetGetPayment

`func (o *WorkflowInstanceHistoryStageInput) GetGetPayment() ActivityGetPayment`

GetGetPayment returns the GetPayment field if non-nil, zero value otherwise.

### GetGetPaymentOk

`func (o *WorkflowInstanceHistoryStageInput) GetGetPaymentOk() (*ActivityGetPayment, bool)`

GetGetPaymentOk returns a tuple with the GetPayment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGetPayment

`func (o *WorkflowInstanceHistoryStageInput) SetGetPayment(v ActivityGetPayment)`

SetGetPayment sets GetPayment field to given value.

### HasGetPayment

`func (o *WorkflowInstanceHistoryStageInput) HasGetPayment() bool`

HasGetPayment returns a boolean if a field has been set.

### GetConfirmHold

`func (o *WorkflowInstanceHistoryStageInput) GetConfirmHold() ActivityConfirmHold`

GetConfirmHold returns the ConfirmHold field if non-nil, zero value otherwise.

### GetConfirmHoldOk

`func (o *WorkflowInstanceHistoryStageInput) GetConfirmHoldOk() (*ActivityConfirmHold, bool)`

GetConfirmHoldOk returns a tuple with the ConfirmHold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfirmHold

`func (o *WorkflowInstanceHistoryStageInput) SetConfirmHold(v ActivityConfirmHold)`

SetConfirmHold sets ConfirmHold field to given value.

### HasConfirmHold

`func (o *WorkflowInstanceHistoryStageInput) HasConfirmHold() bool`

HasConfirmHold returns a boolean if a field has been set.

### GetCreditWallet

`func (o *WorkflowInstanceHistoryStageInput) GetCreditWallet() ActivityCreditWallet`

GetCreditWallet returns the CreditWallet field if non-nil, zero value otherwise.

### GetCreditWalletOk

`func (o *WorkflowInstanceHistoryStageInput) GetCreditWalletOk() (*ActivityCreditWallet, bool)`

GetCreditWalletOk returns a tuple with the CreditWallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreditWallet

`func (o *WorkflowInstanceHistoryStageInput) SetCreditWallet(v ActivityCreditWallet)`

SetCreditWallet sets CreditWallet field to given value.

### HasCreditWallet

`func (o *WorkflowInstanceHistoryStageInput) HasCreditWallet() bool`

HasCreditWallet returns a boolean if a field has been set.

### GetDebitWallet

`func (o *WorkflowInstanceHistoryStageInput) GetDebitWallet() ActivityDebitWallet`

GetDebitWallet returns the DebitWallet field if non-nil, zero value otherwise.

### GetDebitWalletOk

`func (o *WorkflowInstanceHistoryStageInput) GetDebitWalletOk() (*ActivityDebitWallet, bool)`

GetDebitWalletOk returns a tuple with the DebitWallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebitWallet

`func (o *WorkflowInstanceHistoryStageInput) SetDebitWallet(v ActivityDebitWallet)`

SetDebitWallet sets DebitWallet field to given value.

### HasDebitWallet

`func (o *WorkflowInstanceHistoryStageInput) HasDebitWallet() bool`

HasDebitWallet returns a boolean if a field has been set.

### GetGetWallet

`func (o *WorkflowInstanceHistoryStageInput) GetGetWallet() ActivityGetWallet`

GetGetWallet returns the GetWallet field if non-nil, zero value otherwise.

### GetGetWalletOk

`func (o *WorkflowInstanceHistoryStageInput) GetGetWalletOk() (*ActivityGetWallet, bool)`

GetGetWalletOk returns a tuple with the GetWallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGetWallet

`func (o *WorkflowInstanceHistoryStageInput) SetGetWallet(v ActivityGetWallet)`

SetGetWallet sets GetWallet field to given value.

### HasGetWallet

`func (o *WorkflowInstanceHistoryStageInput) HasGetWallet() bool`

HasGetWallet returns a boolean if a field has been set.

### GetVoidHold

`func (o *WorkflowInstanceHistoryStageInput) GetVoidHold() ActivityVoidHold`

GetVoidHold returns the VoidHold field if non-nil, zero value otherwise.

### GetVoidHoldOk

`func (o *WorkflowInstanceHistoryStageInput) GetVoidHoldOk() (*ActivityVoidHold, bool)`

GetVoidHoldOk returns a tuple with the VoidHold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoidHold

`func (o *WorkflowInstanceHistoryStageInput) SetVoidHold(v ActivityVoidHold)`

SetVoidHold sets VoidHold field to given value.

### HasVoidHold

`func (o *WorkflowInstanceHistoryStageInput) HasVoidHold() bool`

HasVoidHold returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


