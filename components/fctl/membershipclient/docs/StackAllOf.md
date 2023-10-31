# StackAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Stack ID | 
**OrganizationId** | **string** | Organization ID | 
**Uri** | **string** | Base stack uri | 
**RegionID** | **string** | The region where the stack is installed | 
**StargateEnabled** | **bool** |  | 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**DisabledAt** | Pointer to **time.Time** |  | [optional] 
**Status** | **string** |  | 
**State** | **string** |  | 
**ExpectedStatus** | **string** |  | 
**LastStateUpdate** | **time.Time** |  | 
**LastExpectedStatusUpdate** | **time.Time** |  | 
**LastStatusUpdate** | **time.Time** |  | 

## Methods

### NewStackAllOf

`func NewStackAllOf(id string, organizationId string, uri string, regionID string, stargateEnabled bool, status string, state string, expectedStatus string, lastStateUpdate time.Time, lastExpectedStatusUpdate time.Time, lastStatusUpdate time.Time, ) *StackAllOf`

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


### GetCreatedAt

`func (o *StackAllOf) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *StackAllOf) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *StackAllOf) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *StackAllOf) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

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

### GetDisabledAt

`func (o *StackAllOf) GetDisabledAt() time.Time`

GetDisabledAt returns the DisabledAt field if non-nil, zero value otherwise.

### GetDisabledAtOk

`func (o *StackAllOf) GetDisabledAtOk() (*time.Time, bool)`

GetDisabledAtOk returns a tuple with the DisabledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledAt

`func (o *StackAllOf) SetDisabledAt(v time.Time)`

SetDisabledAt sets DisabledAt field to given value.

### HasDisabledAt

`func (o *StackAllOf) HasDisabledAt() bool`

HasDisabledAt returns a boolean if a field has been set.

### GetStatus

`func (o *StackAllOf) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *StackAllOf) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *StackAllOf) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetState

`func (o *StackAllOf) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *StackAllOf) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *StackAllOf) SetState(v string)`

SetState sets State field to given value.


### GetExpectedStatus

`func (o *StackAllOf) GetExpectedStatus() string`

GetExpectedStatus returns the ExpectedStatus field if non-nil, zero value otherwise.

### GetExpectedStatusOk

`func (o *StackAllOf) GetExpectedStatusOk() (*string, bool)`

GetExpectedStatusOk returns a tuple with the ExpectedStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedStatus

`func (o *StackAllOf) SetExpectedStatus(v string)`

SetExpectedStatus sets ExpectedStatus field to given value.


### GetLastStateUpdate

`func (o *StackAllOf) GetLastStateUpdate() time.Time`

GetLastStateUpdate returns the LastStateUpdate field if non-nil, zero value otherwise.

### GetLastStateUpdateOk

`func (o *StackAllOf) GetLastStateUpdateOk() (*time.Time, bool)`

GetLastStateUpdateOk returns a tuple with the LastStateUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStateUpdate

`func (o *StackAllOf) SetLastStateUpdate(v time.Time)`

SetLastStateUpdate sets LastStateUpdate field to given value.


### GetLastExpectedStatusUpdate

`func (o *StackAllOf) GetLastExpectedStatusUpdate() time.Time`

GetLastExpectedStatusUpdate returns the LastExpectedStatusUpdate field if non-nil, zero value otherwise.

### GetLastExpectedStatusUpdateOk

`func (o *StackAllOf) GetLastExpectedStatusUpdateOk() (*time.Time, bool)`

GetLastExpectedStatusUpdateOk returns a tuple with the LastExpectedStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastExpectedStatusUpdate

`func (o *StackAllOf) SetLastExpectedStatusUpdate(v time.Time)`

SetLastExpectedStatusUpdate sets LastExpectedStatusUpdate field to given value.


### GetLastStatusUpdate

`func (o *StackAllOf) GetLastStatusUpdate() time.Time`

GetLastStatusUpdate returns the LastStatusUpdate field if non-nil, zero value otherwise.

### GetLastStatusUpdateOk

`func (o *StackAllOf) GetLastStatusUpdateOk() (*time.Time, bool)`

GetLastStatusUpdateOk returns a tuple with the LastStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStatusUpdate

`func (o *StackAllOf) SetLastStatusUpdate(v time.Time)`

SetLastStatusUpdate sets LastStatusUpdate field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


