# Stack

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Stack name | 
**Metadata** | Pointer to **map[string]string** |  | [optional] 
**Version** | Pointer to **string** | Supported only with agent version &gt;&#x3D; v0.7.0 | [optional] 
**Status** | **string** |  | 
**State** | **string** |  | 
**ExpectedStatus** | **string** |  | 
**LastStateUpdate** | **time.Time** |  | 
**LastExpectedStatusUpdate** | **time.Time** |  | 
**LastStatusUpdate** | **time.Time** |  | 
**Reachable** | **bool** | Stack is reachable through Stargate | 
**LastReachableUpdate** | Pointer to **time.Time** | Last time the stack was reachable | [optional] 
**Id** | **string** | Stack ID | 
**OrganizationId** | **string** | Organization ID | 
**Uri** | **string** | Base stack uri | 
**RegionID** | **string** | The region where the stack is installed | 
**StargateEnabled** | **bool** |  | 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DeletedAt** | Pointer to **time.Time** |  | [optional] 
**DisabledAt** | Pointer to **time.Time** |  | [optional] 
**AuditEnabled** | Pointer to **bool** |  | [optional] 
**Synchronised** | **bool** |  | 

## Methods

### NewStack

`func NewStack(name string, status string, state string, expectedStatus string, lastStateUpdate time.Time, lastExpectedStatusUpdate time.Time, lastStatusUpdate time.Time, reachable bool, id string, organizationId string, uri string, regionID string, stargateEnabled bool, synchronised bool, ) *Stack`

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

### HasMetadata

`func (o *Stack) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetVersion

`func (o *Stack) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *Stack) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *Stack) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *Stack) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

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


### GetState

`func (o *Stack) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Stack) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Stack) SetState(v string)`

SetState sets State field to given value.


### GetExpectedStatus

`func (o *Stack) GetExpectedStatus() string`

GetExpectedStatus returns the ExpectedStatus field if non-nil, zero value otherwise.

### GetExpectedStatusOk

`func (o *Stack) GetExpectedStatusOk() (*string, bool)`

GetExpectedStatusOk returns a tuple with the ExpectedStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedStatus

`func (o *Stack) SetExpectedStatus(v string)`

SetExpectedStatus sets ExpectedStatus field to given value.


### GetLastStateUpdate

`func (o *Stack) GetLastStateUpdate() time.Time`

GetLastStateUpdate returns the LastStateUpdate field if non-nil, zero value otherwise.

### GetLastStateUpdateOk

`func (o *Stack) GetLastStateUpdateOk() (*time.Time, bool)`

GetLastStateUpdateOk returns a tuple with the LastStateUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStateUpdate

`func (o *Stack) SetLastStateUpdate(v time.Time)`

SetLastStateUpdate sets LastStateUpdate field to given value.


### GetLastExpectedStatusUpdate

`func (o *Stack) GetLastExpectedStatusUpdate() time.Time`

GetLastExpectedStatusUpdate returns the LastExpectedStatusUpdate field if non-nil, zero value otherwise.

### GetLastExpectedStatusUpdateOk

`func (o *Stack) GetLastExpectedStatusUpdateOk() (*time.Time, bool)`

GetLastExpectedStatusUpdateOk returns a tuple with the LastExpectedStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastExpectedStatusUpdate

`func (o *Stack) SetLastExpectedStatusUpdate(v time.Time)`

SetLastExpectedStatusUpdate sets LastExpectedStatusUpdate field to given value.


### GetLastStatusUpdate

`func (o *Stack) GetLastStatusUpdate() time.Time`

GetLastStatusUpdate returns the LastStatusUpdate field if non-nil, zero value otherwise.

### GetLastStatusUpdateOk

`func (o *Stack) GetLastStatusUpdateOk() (*time.Time, bool)`

GetLastStatusUpdateOk returns a tuple with the LastStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStatusUpdate

`func (o *Stack) SetLastStatusUpdate(v time.Time)`

SetLastStatusUpdate sets LastStatusUpdate field to given value.


### GetReachable

`func (o *Stack) GetReachable() bool`

GetReachable returns the Reachable field if non-nil, zero value otherwise.

### GetReachableOk

`func (o *Stack) GetReachableOk() (*bool, bool)`

GetReachableOk returns a tuple with the Reachable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReachable

`func (o *Stack) SetReachable(v bool)`

SetReachable sets Reachable field to given value.


### GetLastReachableUpdate

`func (o *Stack) GetLastReachableUpdate() time.Time`

GetLastReachableUpdate returns the LastReachableUpdate field if non-nil, zero value otherwise.

### GetLastReachableUpdateOk

`func (o *Stack) GetLastReachableUpdateOk() (*time.Time, bool)`

GetLastReachableUpdateOk returns a tuple with the LastReachableUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastReachableUpdate

`func (o *Stack) SetLastReachableUpdate(v time.Time)`

SetLastReachableUpdate sets LastReachableUpdate field to given value.

### HasLastReachableUpdate

`func (o *Stack) HasLastReachableUpdate() bool`

HasLastReachableUpdate returns a boolean if a field has been set.

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

### GetAuditEnabled

`func (o *Stack) GetAuditEnabled() bool`

GetAuditEnabled returns the AuditEnabled field if non-nil, zero value otherwise.

### GetAuditEnabledOk

`func (o *Stack) GetAuditEnabledOk() (*bool, bool)`

GetAuditEnabledOk returns a tuple with the AuditEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuditEnabled

`func (o *Stack) SetAuditEnabled(v bool)`

SetAuditEnabled sets AuditEnabled field to given value.

### HasAuditEnabled

`func (o *Stack) HasAuditEnabled() bool`

HasAuditEnabled returns a boolean if a field has been set.

### GetSynchronised

`func (o *Stack) GetSynchronised() bool`

GetSynchronised returns the Synchronised field if non-nil, zero value otherwise.

### GetSynchronisedOk

`func (o *Stack) GetSynchronisedOk() (*bool, bool)`

GetSynchronisedOk returns a tuple with the Synchronised field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynchronised

`func (o *Stack) SetSynchronised(v bool)`

SetSynchronised sets Synchronised field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


