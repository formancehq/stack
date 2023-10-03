# Stack

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Stack name | 
**Metadata** | **map[string]string** |  | 
**Id** | **string** | Stack ID | 
**OrganizationId** | **string** | Organization ID | 
**Uri** | **string** | Base stack uri | 
**RegionID** | **string** | The region where the stack is installed | 
**StargateEnabled** | **bool** |  | 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**DisabledAt** | Pointer to **time.Time** |  | [optional] 
**Status** | **string** |  | 

## Methods

### NewStack

`func NewStack(name string, metadata map[string]string, id string, organizationId string, uri string, regionID string, stargateEnabled bool, status string, ) *Stack`

NewStack instantiates a new Stack object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackWithDefaults

`func NewStackWithDefaults() *Stack`

NewStackWithDefaults instantiates a new Stack object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Stack) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Stack) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Stack) SetName(v string)`

SetName sets Name field to given value.


### GetMetadata

`func (o *Stack) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *Stack) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *Stack) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.


### GetId

`func (o *Stack) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Stack) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Stack) SetId(v string)`

SetId sets Id field to given value.


### GetOrganizationId

`func (o *Stack) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *Stack) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *Stack) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.


### GetUri

`func (o *Stack) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *Stack) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *Stack) SetUri(v string)`

SetUri sets Uri field to given value.


### GetRegionID

`func (o *Stack) GetRegionID() string`

GetRegionID returns the RegionID field if non-nil, zero value otherwise.

### GetRegionIDOk

`func (o *Stack) GetRegionIDOk() (*string, bool)`

GetRegionIDOk returns a tuple with the RegionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionID

`func (o *Stack) SetRegionID(v string)`

SetRegionID sets RegionID field to given value.


### GetStargateEnabled

`func (o *Stack) GetStargateEnabled() bool`

GetStargateEnabled returns the StargateEnabled field if non-nil, zero value otherwise.

### GetStargateEnabledOk

`func (o *Stack) GetStargateEnabledOk() (*bool, bool)`

GetStargateEnabledOk returns a tuple with the StargateEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStargateEnabled

`func (o *Stack) SetStargateEnabled(v bool)`

SetStargateEnabled sets StargateEnabled field to given value.


### GetCreatedAt

`func (o *Stack) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Stack) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Stack) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Stack) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDeletedAt

`func (o *Stack) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *Stack) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *Stack) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *Stack) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.

### GetDisabledAt

`func (o *Stack) GetDisabledAt() time.Time`

GetDisabledAt returns the DisabledAt field if non-nil, zero value otherwise.

### GetDisabledAtOk

`func (o *Stack) GetDisabledAtOk() (*time.Time, bool)`

GetDisabledAtOk returns a tuple with the DisabledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledAt

`func (o *Stack) SetDisabledAt(v time.Time)`

SetDisabledAt sets DisabledAt field to given value.

### HasDisabledAt

`func (o *Stack) HasDisabledAt() bool`

HasDisabledAt returns a boolean if a field has been set.

### GetStatus

`func (o *Stack) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Stack) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Stack) SetStatus(v string)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


