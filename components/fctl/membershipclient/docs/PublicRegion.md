# PublicRegion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**BaseUrl** | **string** |  | 
**CreatedAt** | **string** |  | 
**Active** | **bool** |  | 
**LastPing** | Pointer to **time.Time** |  | [optional] 
**Name** | **string** |  | 
**Production** | **bool** |  | 

## Methods

### NewPublicRegion

`func NewPublicRegion(id string, baseUrl string, createdAt string, active bool, name string, production bool, ) *PublicRegion`

NewPublicRegion instantiates a new PublicRegion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPublicRegionWithDefaults

`func NewPublicRegionWithDefaults() *PublicRegion`

NewPublicRegionWithDefaults instantiates a new PublicRegion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PublicRegion) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PublicRegion) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PublicRegion) SetId(v string)`

SetId sets Id field to given value.


### GetBaseUrl

`func (o *PublicRegion) GetBaseUrl() string`

GetBaseUrl returns the BaseUrl field if non-nil, zero value otherwise.

### GetBaseUrlOk

`func (o *PublicRegion) GetBaseUrlOk() (*string, bool)`

GetBaseUrlOk returns a tuple with the BaseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseUrl

`func (o *PublicRegion) SetBaseUrl(v string)`

SetBaseUrl sets BaseUrl field to given value.


### GetCreatedAt

`func (o *PublicRegion) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PublicRegion) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PublicRegion) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetActive

`func (o *PublicRegion) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *PublicRegion) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *PublicRegion) SetActive(v bool)`

SetActive sets Active field to given value.


### GetLastPing

`func (o *PublicRegion) GetLastPing() time.Time`

GetLastPing returns the LastPing field if non-nil, zero value otherwise.

### GetLastPingOk

`func (o *PublicRegion) GetLastPingOk() (*time.Time, bool)`

GetLastPingOk returns a tuple with the LastPing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPing

`func (o *PublicRegion) SetLastPing(v time.Time)`

SetLastPing sets LastPing field to given value.

### HasLastPing

`func (o *PublicRegion) HasLastPing() bool`

HasLastPing returns a boolean if a field has been set.

### GetName

`func (o *PublicRegion) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PublicRegion) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PublicRegion) SetName(v string)`

SetName sets Name field to given value.


### GetProduction

`func (o *PublicRegion) GetProduction() bool`

GetProduction returns the Production field if non-nil, zero value otherwise.

### GetProductionOk

`func (o *PublicRegion) GetProductionOk() (*bool, bool)`

GetProductionOk returns a tuple with the Production field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProduction

`func (o *PublicRegion) SetProduction(v bool)`

SetProduction sets Production field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


