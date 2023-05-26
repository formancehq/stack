# StackAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Stack ID | 
**OrganizationId** | **string** | Organization ID | 
**Uri** | **string** | Base stack uri | 
**RegionID** | **string** | The region where the stack is installed | 
**StargateEnabled** | **bool** |  | 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewStackAllOf

`func NewStackAllOf(id string, organizationId string, uri string, regionID string, stargateEnabled bool, ) *StackAllOf`

NewStackAllOf instantiates a new StackAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackAllOfWithDefaults

`func NewStackAllOfWithDefaults() *StackAllOf`

NewStackAllOfWithDefaults instantiates a new StackAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *StackAllOf) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *StackAllOf) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *StackAllOf) SetId(v string)`

SetId sets Id field to given value.


### GetOrganizationId

`func (o *StackAllOf) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *StackAllOf) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *StackAllOf) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.


### GetUri

`func (o *StackAllOf) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *StackAllOf) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *StackAllOf) SetUri(v string)`

SetUri sets Uri field to given value.


### GetRegionID

`func (o *StackAllOf) GetRegionID() string`

GetRegionID returns the RegionID field if non-nil, zero value otherwise.

### GetRegionIDOk

`func (o *StackAllOf) GetRegionIDOk() (*string, bool)`

GetRegionIDOk returns a tuple with the RegionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionID

`func (o *StackAllOf) SetRegionID(v string)`

SetRegionID sets RegionID field to given value.


### GetStargateEnabled

`func (o *StackAllOf) GetStargateEnabled() bool`

GetStargateEnabled returns the StargateEnabled field if non-nil, zero value otherwise.

### GetStargateEnabledOk

`func (o *StackAllOf) GetStargateEnabledOk() (*bool, bool)`

GetStargateEnabledOk returns a tuple with the StargateEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStargateEnabled

`func (o *StackAllOf) SetStargateEnabled(v bool)`

SetStargateEnabled sets StargateEnabled field to given value.


### GetDeletedAt

`func (o *StackAllOf) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *StackAllOf) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *StackAllOf) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *StackAllOf) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


