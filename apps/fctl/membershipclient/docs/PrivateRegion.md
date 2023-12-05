# PrivateRegion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**BaseUrl** | **string** |  | 
**CreatedAt** | **string** |  | 
**Active** | **bool** |  | 
**LastPing** | Pointer to **time.Time** |  | [optional] 
**Name** | **string** |  | 
**OrganizationID** | **string** |  | 
**CreatorID** | **string** |  | 
**Secret** | Pointer to [**PrivateRegionAllOfSecret**](PrivateRegionAllOfSecret.md) |  | [optional] 

## Methods

### NewPrivateRegion

`func NewPrivateRegion(id string, baseUrl string, createdAt string, active bool, name string, organizationID string, creatorID string, ) *PrivateRegion`

NewPrivateRegion instantiates a new PrivateRegion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivateRegionWithDefaults

`func NewPrivateRegionWithDefaults() *PrivateRegion`

NewPrivateRegionWithDefaults instantiates a new PrivateRegion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PrivateRegion) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PrivateRegion) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PrivateRegion) SetId(v string)`

SetId sets Id field to given value.


### GetBaseUrl

`func (o *PrivateRegion) GetBaseUrl() string`

GetBaseUrl returns the BaseUrl field if non-nil, zero value otherwise.

### GetBaseUrlOk

`func (o *PrivateRegion) GetBaseUrlOk() (*string, bool)`

GetBaseUrlOk returns a tuple with the BaseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseUrl

`func (o *PrivateRegion) SetBaseUrl(v string)`

SetBaseUrl sets BaseUrl field to given value.


### GetCreatedAt

`func (o *PrivateRegion) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PrivateRegion) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PrivateRegion) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetActive

`func (o *PrivateRegion) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *PrivateRegion) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *PrivateRegion) SetActive(v bool)`

SetActive sets Active field to given value.


### GetLastPing

`func (o *PrivateRegion) GetLastPing() time.Time`

GetLastPing returns the LastPing field if non-nil, zero value otherwise.

### GetLastPingOk

`func (o *PrivateRegion) GetLastPingOk() (*time.Time, bool)`

GetLastPingOk returns a tuple with the LastPing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPing

`func (o *PrivateRegion) SetLastPing(v time.Time)`

SetLastPing sets LastPing field to given value.

### HasLastPing

`func (o *PrivateRegion) HasLastPing() bool`

HasLastPing returns a boolean if a field has been set.

### GetName

`func (o *PrivateRegion) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PrivateRegion) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PrivateRegion) SetName(v string)`

SetName sets Name field to given value.


### GetOrganizationID

`func (o *PrivateRegion) GetOrganizationID() string`

GetOrganizationID returns the OrganizationID field if non-nil, zero value otherwise.

### GetOrganizationIDOk

`func (o *PrivateRegion) GetOrganizationIDOk() (*string, bool)`

GetOrganizationIDOk returns a tuple with the OrganizationID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationID

`func (o *PrivateRegion) SetOrganizationID(v string)`

SetOrganizationID sets OrganizationID field to given value.


### GetCreatorID

`func (o *PrivateRegion) GetCreatorID() string`

GetCreatorID returns the CreatorID field if non-nil, zero value otherwise.

### GetCreatorIDOk

`func (o *PrivateRegion) GetCreatorIDOk() (*string, bool)`

GetCreatorIDOk returns a tuple with the CreatorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatorID

`func (o *PrivateRegion) SetCreatorID(v string)`

SetCreatorID sets CreatorID field to given value.


### GetSecret

`func (o *PrivateRegion) GetSecret() PrivateRegionAllOfSecret`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *PrivateRegion) GetSecretOk() (*PrivateRegionAllOfSecret, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *PrivateRegion) SetSecret(v PrivateRegionAllOfSecret)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *PrivateRegion) HasSecret() bool`

HasSecret returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


