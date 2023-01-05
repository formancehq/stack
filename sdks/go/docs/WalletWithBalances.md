# WalletWithBalances

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The unique ID of the wallet. | 
**Metadata** | **map[string]interface{}** | Metadata associated with the wallet. | 
**Name** | **string** |  | 
**Balances** | [**WalletWithBalancesAllOfBalances**](WalletWithBalancesAllOfBalances.md) |  | 

## Methods

### NewWalletWithBalances

`func NewWalletWithBalances(id string, metadata map[string]interface{}, name string, balances WalletWithBalancesAllOfBalances, ) *WalletWithBalances`

NewWalletWithBalances instantiates a new WalletWithBalances object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWalletWithBalancesWithDefaults

`func NewWalletWithBalancesWithDefaults() *WalletWithBalances`

NewWalletWithBalancesWithDefaults instantiates a new WalletWithBalances object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *WalletWithBalances) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WalletWithBalances) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WalletWithBalances) SetId(v string)`

SetId sets Id field to given value.


### GetMetadata

`func (o *WalletWithBalances) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *WalletWithBalances) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *WalletWithBalances) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.


### GetName

`func (o *WalletWithBalances) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WalletWithBalances) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WalletWithBalances) SetName(v string)`

SetName sets Name field to given value.


### GetBalances

`func (o *WalletWithBalances) GetBalances() WalletWithBalancesAllOfBalances`

GetBalances returns the Balances field if non-nil, zero value otherwise.

### GetBalancesOk

`func (o *WalletWithBalances) GetBalancesOk() (*WalletWithBalancesAllOfBalances, bool)`

GetBalancesOk returns a tuple with the Balances field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalances

`func (o *WalletWithBalances) SetBalances(v WalletWithBalancesAllOfBalances)`

SetBalances sets Balances field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


