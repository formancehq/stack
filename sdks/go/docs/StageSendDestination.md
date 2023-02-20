# StageSendDestination

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Wallet** | Pointer to [**StageSendSourceWallet**](StageSendSourceWallet.md) |  | [optional] 
**Account** | Pointer to [**StageSendSourceAccount**](StageSendSourceAccount.md) |  | [optional] 
**Payment** | Pointer to [**StageSendDestinationPayment**](StageSendDestinationPayment.md) |  | [optional] 

## Methods

### NewStageSendDestination

`func NewStageSendDestination() *StageSendDestination`

NewStageSendDestination instantiates a new StageSendDestination object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStageSendDestinationWithDefaults

`func NewStageSendDestinationWithDefaults() *StageSendDestination`

NewStageSendDestinationWithDefaults instantiates a new StageSendDestination object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWallet

`func (o *StageSendDestination) GetWallet() StageSendSourceWallet`

GetWallet returns the Wallet field if non-nil, zero value otherwise.

### GetWalletOk

`func (o *StageSendDestination) GetWalletOk() (*StageSendSourceWallet, bool)`

GetWalletOk returns a tuple with the Wallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWallet

`func (o *StageSendDestination) SetWallet(v StageSendSourceWallet)`

SetWallet sets Wallet field to given value.

### HasWallet

`func (o *StageSendDestination) HasWallet() bool`

HasWallet returns a boolean if a field has been set.

### GetAccount

`func (o *StageSendDestination) GetAccount() StageSendSourceAccount`

GetAccount returns the Account field if non-nil, zero value otherwise.

### GetAccountOk

`func (o *StageSendDestination) GetAccountOk() (*StageSendSourceAccount, bool)`

GetAccountOk returns a tuple with the Account field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccount

`func (o *StageSendDestination) SetAccount(v StageSendSourceAccount)`

SetAccount sets Account field to given value.

### HasAccount

`func (o *StageSendDestination) HasAccount() bool`

HasAccount returns a boolean if a field has been set.

### GetPayment

`func (o *StageSendDestination) GetPayment() StageSendDestinationPayment`

GetPayment returns the Payment field if non-nil, zero value otherwise.

### GetPaymentOk

`func (o *StageSendDestination) GetPaymentOk() (*StageSendDestinationPayment, bool)`

GetPaymentOk returns a tuple with the Payment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayment

`func (o *StageSendDestination) SetPayment(v StageSendDestinationPayment)`

SetPayment sets Payment field to given value.

### HasPayment

`func (o *StageSendDestination) HasPayment() bool`

HasPayment returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


