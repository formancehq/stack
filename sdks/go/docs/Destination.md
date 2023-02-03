# Destination

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Account** | Pointer to [**LedgerAccountSource**](LedgerAccountSource.md) |  | [optional] 
**Payment** | Pointer to [**PaymentDestination**](PaymentDestination.md) |  | [optional] 
**Wallet** | Pointer to [**WalletSource**](WalletSource.md) |  | [optional] 

## Methods

### NewDestination

`func NewDestination() *Destination`

NewDestination instantiates a new Destination object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDestinationWithDefaults

`func NewDestinationWithDefaults() *Destination`

NewDestinationWithDefaults instantiates a new Destination object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccount

`func (o *Destination) GetAccount() LedgerAccountSource`

GetAccount returns the Account field if non-nil, zero value otherwise.

### GetAccountOk

`func (o *Destination) GetAccountOk() (*LedgerAccountSource, bool)`

GetAccountOk returns a tuple with the Account field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccount

`func (o *Destination) SetAccount(v LedgerAccountSource)`

SetAccount sets Account field to given value.

### HasAccount

`func (o *Destination) HasAccount() bool`

HasAccount returns a boolean if a field has been set.

### GetPayment

`func (o *Destination) GetPayment() PaymentDestination`

GetPayment returns the Payment field if non-nil, zero value otherwise.

### GetPaymentOk

`func (o *Destination) GetPaymentOk() (*PaymentDestination, bool)`

GetPaymentOk returns a tuple with the Payment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayment

`func (o *Destination) SetPayment(v PaymentDestination)`

SetPayment sets Payment field to given value.

### HasPayment

`func (o *Destination) HasPayment() bool`

HasPayment returns a boolean if a field has been set.

### GetWallet

`func (o *Destination) GetWallet() WalletSource`

GetWallet returns the Wallet field if non-nil, zero value otherwise.

### GetWalletOk

`func (o *Destination) GetWalletOk() (*WalletSource, bool)`

GetWalletOk returns a tuple with the Wallet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWallet

`func (o *Destination) SetWallet(v WalletSource)`

SetWallet sets Wallet field to given value.

### HasWallet

`func (o *Destination) HasWallet() bool`

HasWallet returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


