# OrganizationExpanded

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Organization name | 
**DefaultOrganizationAccess** | Pointer to [**Role**](Role.md) |  | [optional] [default to EMPTY]
**DefaultStackAccess** | Pointer to [**Role**](Role.md) |  | [optional] [default to EMPTY]
**Domain** | Pointer to **string** | Organization domain | [optional] 
**Id** | **string** | Organization ID | 
**OwnerId** | **string** | Owner ID | 
**AvailableStacks** | Pointer to **int32** | Number of available stacks | [optional] 
**AvailableSandboxes** | Pointer to **int32** | Number of available sandboxes | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**TotalStacks** | Pointer to **int32** |  | [optional] 
**TotalUsers** | Pointer to **int32** |  | [optional] 
**Owner** | Pointer to [**User**](User.md) |  | [optional] 

## Methods

### NewOrganizationExpanded

`func NewOrganizationExpanded(name string, id string, ownerId string, ) *OrganizationExpanded`

NewOrganizationExpanded instantiates a new OrganizationExpanded object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationExpandedWithDefaults

`func NewOrganizationExpandedWithDefaults() *OrganizationExpanded`

NewOrganizationExpandedWithDefaults instantiates a new OrganizationExpanded object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *OrganizationExpanded) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrganizationExpanded) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrganizationExpanded) SetName(v string)`

SetName sets Name field to given value.


### GetDefaultOrganizationAccess

`func (o *OrganizationExpanded) GetDefaultOrganizationAccess() Role`

GetDefaultOrganizationAccess returns the DefaultOrganizationAccess field if non-nil, zero value otherwise.

### GetDefaultOrganizationAccessOk

`func (o *OrganizationExpanded) GetDefaultOrganizationAccessOk() (*Role, bool)`

GetDefaultOrganizationAccessOk returns a tuple with the DefaultOrganizationAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultOrganizationAccess

`func (o *OrganizationExpanded) SetDefaultOrganizationAccess(v Role)`

SetDefaultOrganizationAccess sets DefaultOrganizationAccess field to given value.

### HasDefaultOrganizationAccess

`func (o *OrganizationExpanded) HasDefaultOrganizationAccess() bool`

HasDefaultOrganizationAccess returns a boolean if a field has been set.

### GetDefaultStackAccess

`func (o *OrganizationExpanded) GetDefaultStackAccess() Role`

GetDefaultStackAccess returns the DefaultStackAccess field if non-nil, zero value otherwise.

### GetDefaultStackAccessOk

`func (o *OrganizationExpanded) GetDefaultStackAccessOk() (*Role, bool)`

GetDefaultStackAccessOk returns a tuple with the DefaultStackAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultStackAccess

`func (o *OrganizationExpanded) SetDefaultStackAccess(v Role)`

SetDefaultStackAccess sets DefaultStackAccess field to given value.

### HasDefaultStackAccess

`func (o *OrganizationExpanded) HasDefaultStackAccess() bool`

HasDefaultStackAccess returns a boolean if a field has been set.

### GetDomain

`func (o *OrganizationExpanded) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *OrganizationExpanded) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *OrganizationExpanded) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *OrganizationExpanded) HasDomain() bool`

HasDomain returns a boolean if a field has been set.

### GetId

`func (o *OrganizationExpanded) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OrganizationExpanded) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OrganizationExpanded) SetId(v string)`

SetId sets Id field to given value.


### GetOwnerId

`func (o *OrganizationExpanded) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *OrganizationExpanded) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *OrganizationExpanded) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.


### GetAvailableStacks

`func (o *OrganizationExpanded) GetAvailableStacks() int32`

GetAvailableStacks returns the AvailableStacks field if non-nil, zero value otherwise.

### GetAvailableStacksOk

`func (o *OrganizationExpanded) GetAvailableStacksOk() (*int32, bool)`

GetAvailableStacksOk returns a tuple with the AvailableStacks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableStacks

`func (o *OrganizationExpanded) SetAvailableStacks(v int32)`

SetAvailableStacks sets AvailableStacks field to given value.

### HasAvailableStacks

`func (o *OrganizationExpanded) HasAvailableStacks() bool`

HasAvailableStacks returns a boolean if a field has been set.

### GetAvailableSandboxes

`func (o *OrganizationExpanded) GetAvailableSandboxes() int32`

GetAvailableSandboxes returns the AvailableSandboxes field if non-nil, zero value otherwise.

### GetAvailableSandboxesOk

`func (o *OrganizationExpanded) GetAvailableSandboxesOk() (*int32, bool)`

GetAvailableSandboxesOk returns a tuple with the AvailableSandboxes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableSandboxes

`func (o *OrganizationExpanded) SetAvailableSandboxes(v int32)`

SetAvailableSandboxes sets AvailableSandboxes field to given value.

### HasAvailableSandboxes

`func (o *OrganizationExpanded) HasAvailableSandboxes() bool`

HasAvailableSandboxes returns a boolean if a field has been set.

### GetCreatedAt

`func (o *OrganizationExpanded) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OrganizationExpanded) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OrganizationExpanded) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *OrganizationExpanded) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *OrganizationExpanded) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *OrganizationExpanded) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *OrganizationExpanded) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *OrganizationExpanded) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetTotalStacks

`func (o *OrganizationExpanded) GetTotalStacks() int32`

GetTotalStacks returns the TotalStacks field if non-nil, zero value otherwise.

### GetTotalStacksOk

`func (o *OrganizationExpanded) GetTotalStacksOk() (*int32, bool)`

GetTotalStacksOk returns a tuple with the TotalStacks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalStacks

`func (o *OrganizationExpanded) SetTotalStacks(v int32)`

SetTotalStacks sets TotalStacks field to given value.

### HasTotalStacks

`func (o *OrganizationExpanded) HasTotalStacks() bool`

HasTotalStacks returns a boolean if a field has been set.

### GetTotalUsers

`func (o *OrganizationExpanded) GetTotalUsers() int32`

GetTotalUsers returns the TotalUsers field if non-nil, zero value otherwise.

### GetTotalUsersOk

`func (o *OrganizationExpanded) GetTotalUsersOk() (*int32, bool)`

GetTotalUsersOk returns a tuple with the TotalUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalUsers

`func (o *OrganizationExpanded) SetTotalUsers(v int32)`

SetTotalUsers sets TotalUsers field to given value.

### HasTotalUsers

`func (o *OrganizationExpanded) HasTotalUsers() bool`

HasTotalUsers returns a boolean if a field has been set.

### GetOwner

`func (o *OrganizationExpanded) GetOwner() User`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *OrganizationExpanded) GetOwnerOk() (*User, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *OrganizationExpanded) SetOwner(v User)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *OrganizationExpanded) HasOwner() bool`

HasOwner returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


