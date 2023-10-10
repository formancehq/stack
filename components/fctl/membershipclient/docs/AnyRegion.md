# AnyRegion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**BaseUrl** | **string** |  | 
**CreatedAt** | **string** |  | 
**Active** | **bool** |  | 
**LastPing** | Pointer to **time.Time** |  | [optional] 
**Name** | **string** |  | 
**ClientID** | Pointer to **string** |  | [optional] 
**OrganizationID** | Pointer to **string** |  | [optional] 
**Creator** | Pointer to [**User**](User.md) |  | [optional] 
**Production** | Pointer to **bool** |  | [optional] 
**Public** | **bool** |  | 
**Version** | Pointer to **string** |  | [optional] 

## Methods

### NewAnyRegion

`func NewAnyRegion(id string, baseUrl string, createdAt string, active bool, name string, public bool, ) *AnyRegion`

NewAnyRegion instantiates a new AnyRegion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnyRegionWithDefaults

`func NewAnyRegionWithDefaults() *AnyRegion`

NewAnyRegionWithDefaults instantiates a new AnyRegion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AnyRegion) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AnyRegion) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AnyRegion) SetId(v string)`

SetId sets Id field to given value.


### GetBaseUrl

`func (o *AnyRegion) GetBaseUrl() string`

GetBaseUrl returns the BaseUrl field if non-nil, zero value otherwise.

### GetBaseUrlOk

`func (o *AnyRegion) GetBaseUrlOk() (*string, bool)`

GetBaseUrlOk returns a tuple with the BaseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseUrl

`func (o *AnyRegion) SetBaseUrl(v string)`

SetBaseUrl sets BaseUrl field to given value.


### GetCreatedAt

`func (o *AnyRegion) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AnyRegion) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AnyRegion) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetActive

`func (o *AnyRegion) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *AnyRegion) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *AnyRegion) SetActive(v bool)`

SetActive sets Active field to given value.


### GetLastPing

`func (o *AnyRegion) GetLastPing() time.Time`

GetLastPing returns the LastPing field if non-nil, zero value otherwise.

### GetLastPingOk

`func (o *AnyRegion) GetLastPingOk() (*time.Time, bool)`

GetLastPingOk returns a tuple with the LastPing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPing

`func (o *AnyRegion) SetLastPing(v time.Time)`

SetLastPing sets LastPing field to given value.

### HasLastPing

`func (o *AnyRegion) HasLastPing() bool`

HasLastPing returns a boolean if a field has been set.

### GetName

`func (o *AnyRegion) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AnyRegion) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AnyRegion) SetName(v string)`

SetName sets Name field to given value.


### GetClientID

`func (o *AnyRegion) GetClientID() string`

GetClientID returns the ClientID field if non-nil, zero value otherwise.

### GetClientIDOk

`func (o *AnyRegion) GetClientIDOk() (*string, bool)`

GetClientIDOk returns a tuple with the ClientID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientID

`func (o *AnyRegion) SetClientID(v string)`

SetClientID sets ClientID field to given value.

### HasClientID

`func (o *AnyRegion) HasClientID() bool`

HasClientID returns a boolean if a field has been set.

### GetOrganizationID

`func (o *AnyRegion) GetOrganizationID() string`

GetOrganizationID returns the OrganizationID field if non-nil, zero value otherwise.

### GetOrganizationIDOk

`func (o *AnyRegion) GetOrganizationIDOk() (*string, bool)`

GetOrganizationIDOk returns a tuple with the OrganizationID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationID

`func (o *AnyRegion) SetOrganizationID(v string)`

SetOrganizationID sets OrganizationID field to given value.

### HasOrganizationID

`func (o *AnyRegion) HasOrganizationID() bool`

HasOrganizationID returns a boolean if a field has been set.

### GetCreator

`func (o *AnyRegion) GetCreator() User`

GetCreator returns the Creator field if non-nil, zero value otherwise.

### GetCreatorOk

`func (o *AnyRegion) GetCreatorOk() (*User, bool)`

GetCreatorOk returns a tuple with the Creator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreator

`func (o *AnyRegion) SetCreator(v User)`

SetCreator sets Creator field to given value.

### HasCreator

`func (o *AnyRegion) HasCreator() bool`

HasCreator returns a boolean if a field has been set.

### GetProduction

`func (o *AnyRegion) GetProduction() bool`

GetProduction returns the Production field if non-nil, zero value otherwise.

### GetProductionOk

`func (o *AnyRegion) GetProductionOk() (*bool, bool)`

GetProductionOk returns a tuple with the Production field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProduction

`func (o *AnyRegion) SetProduction(v bool)`

SetProduction sets Production field to given value.

### HasProduction

`func (o *AnyRegion) HasProduction() bool`

HasProduction returns a boolean if a field has been set.

### GetPublic

`func (o *AnyRegion) GetPublic() bool`

GetPublic returns the Public field if non-nil, zero value otherwise.

### GetPublicOk

`func (o *AnyRegion) GetPublicOk() (*bool, bool)`

GetPublicOk returns a tuple with the Public field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublic

`func (o *AnyRegion) SetPublic(v bool)`

SetPublic sets Public field to given value.


### GetVersion

`func (o *AnyRegion) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *AnyRegion) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *AnyRegion) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *AnyRegion) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


