# ActivityDebitWallet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Data** | Pointer to [**DebitWalletRequest**](DebitWalletRequest.md) |  | [optional] 

## Methods

### NewActivityDebitWallet

`func NewActivityDebitWallet() *ActivityDebitWallet`

NewActivityDebitWallet instantiates a new ActivityDebitWallet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActivityDebitWalletWithDefaults

`func NewActivityDebitWalletWithDefaults() *ActivityDebitWallet`

NewActivityDebitWalletWithDefaults instantiates a new ActivityDebitWallet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ActivityDebitWallet) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ActivityDebitWallet) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ActivityDebitWallet) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ActivityDebitWallet) HasId() bool`

HasId returns a boolean if a field has been set.

### GetData

`func (o *ActivityDebitWallet) GetData() DebitWalletRequest`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ActivityDebitWallet) GetDataOk() (*DebitWalletRequest, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ActivityDebitWallet) SetData(v DebitWalletRequest)`

SetData sets Data field to given value.

### HasData

`func (o *ActivityDebitWallet) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


