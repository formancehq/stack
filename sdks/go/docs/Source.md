# Source

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Account** | Pointer to [**LedgerAccountSource**](LedgerAccountSource.md) |  | [optional] 
**Payment** | Pointer to [**PaymentSource**](PaymentSource.md) |  | [optional] 
**Wallet** | Pointer to [**WalletSource**](WalletSource.md) |  | [optional] 

## Methods

### NewSource

`func NewSource() *Source`

NewSource instantiates a new Source object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSourceWithDefaults

`func NewSourceWithDefaults() *Source`

NewSourceWithDefaults instantiates a new Source object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccount

`func (o *Source) GetAccount() LedgerAccountSource`

GetAccount returns the Account field if non-nil, zero value otherwise.

### GetAccountOk

`func (o *Source) GetAccountOk() (*LedgerAccountSource, bool)`

GetAccountOk returns a tuple with the Account field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccount

`func (o *Source) SetAccount(v LedgerAccountSource)`

SetAccount sets Account field to given value.

### HasAccount

`func (o *Source) HasAccount() bool`

HasAccount returns a boolean if a field has been set.

### GetPayment

`func (o *Source) GetPayment() PaymentSource`

GetPayment returns the Payment field if non-nil, zero value otherwise.

### GetPaymentOk

`func (o *Source) GetPaymentOk() (*PaymentSource, bool)`

GetPaymentOk returns a tuple with the Payment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayment

`func (o *Source) SetPayment(v PaymentSource)`

SetPayment sets Payment field to given value.

### HasPayment

`func (o *Source) HasPayment() bool`

HasPayment returns a boolean if a field has been set.

### GetWallet

`func (o *Source) GetWallet() WalletSource`

GetWallet returns the Wallet field if non-nil, zero value otherwise.

### GetWalletOk

`func (o *Source) GetWalletOk() (*WalletSource, bool)`

GetWalletOk returns a tuple with the Wallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWallet

`func (o *Source) SetWallet(v WalletSource)`

SetWallet sets Wallet field to given value.

### HasWallet

`func (o *Source) HasWallet() bool`

HasWallet returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


