# BalanceWithAssets

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**ExpiresAt** | Pointer to **time.Time** |  | [optional] 
**Assets** | **map[string]int64** |  | 

## Methods

### NewBalanceWithAssets

`func NewBalanceWithAssets(name string, assets map[string]int64, ) *BalanceWithAssets`

NewBalanceWithAssets instantiates a new BalanceWithAssets object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBalanceWithAssetsWithDefaults

`func NewBalanceWithAssetsWithDefaults() *BalanceWithAssets`

NewBalanceWithAssetsWithDefaults instantiates a new BalanceWithAssets object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BalanceWithAssets) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BalanceWithAssets) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BalanceWithAssets) SetName(v string)`

SetName sets Name field to given value.


### GetExpiresAt

`func (o *BalanceWithAssets) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *BalanceWithAssets) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *BalanceWithAssets) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *BalanceWithAssets) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetAssets

`func (o *BalanceWithAssets) GetAssets() map[string]int64`

GetAssets returns the Assets field if non-nil, zero value otherwise.

### GetAssetsOk

`func (o *BalanceWithAssets) GetAssetsOk() (*map[string]int64, bool)`

GetAssetsOk returns a tuple with the Assets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssets

`func (o *BalanceWithAssets) SetAssets(v map[string]int64)`

SetAssets sets Assets field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


