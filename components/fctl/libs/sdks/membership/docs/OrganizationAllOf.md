# OrganizationAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Organization ID | 
**OwnerId** | **string** | Owner ID | 
**AvailableStacks** | Pointer to **int32** | Number of available stacks | [optional] 
**AvailableSandboxes** | Pointer to **int32** | Number of available sandboxes | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewOrganizationAllOf

`func NewOrganizationAllOf(id string, ownerId string, ) *OrganizationAllOf`

NewOrganizationAllOf instantiates a new OrganizationAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationAllOfWithDefaults

`func NewOrganizationAllOfWithDefaults() *OrganizationAllOf`

NewOrganizationAllOfWithDefaults instantiates a new OrganizationAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OrganizationAllOf) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OrganizationAllOf) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OrganizationAllOf) SetId(v string)`

SetId sets Id field to given value.


### GetOwnerId

`func (o *OrganizationAllOf) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *OrganizationAllOf) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *OrganizationAllOf) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.


### GetAvailableStacks

`func (o *OrganizationAllOf) GetAvailableStacks() int32`

GetAvailableStacks returns the AvailableStacks field if non-nil, zero value otherwise.

### GetAvailableStacksOk

`func (o *OrganizationAllOf) GetAvailableStacksOk() (*int32, bool)`

GetAvailableStacksOk returns a tuple with the AvailableStacks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableStacks

`func (o *OrganizationAllOf) SetAvailableStacks(v int32)`

SetAvailableStacks sets AvailableStacks field to given value.

### HasAvailableStacks

`func (o *OrganizationAllOf) HasAvailableStacks() bool`

HasAvailableStacks returns a boolean if a field has been set.

### GetAvailableSandboxes

`func (o *OrganizationAllOf) GetAvailableSandboxes() int32`

GetAvailableSandboxes returns the AvailableSandboxes field if non-nil, zero value otherwise.

### GetAvailableSandboxesOk

`func (o *OrganizationAllOf) GetAvailableSandboxesOk() (*int32, bool)`

GetAvailableSandboxesOk returns a tuple with the AvailableSandboxes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableSandboxes

`func (o *OrganizationAllOf) SetAvailableSandboxes(v int32)`

SetAvailableSandboxes sets AvailableSandboxes field to given value.

### HasAvailableSandboxes

`func (o *OrganizationAllOf) HasAvailableSandboxes() bool`

HasAvailableSandboxes returns a boolean if a field has been set.

### GetCreatedAt

`func (o *OrganizationAllOf) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OrganizationAllOf) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OrganizationAllOf) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *OrganizationAllOf) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *OrganizationAllOf) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *OrganizationAllOf) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *OrganizationAllOf) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *OrganizationAllOf) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


