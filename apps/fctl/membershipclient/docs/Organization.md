# Organization

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Organization name | 
**Id** | **string** | Organization ID | 
**OwnerId** | **string** | Owner ID | 
**AvailableStacks** | Pointer to **int32** | Number of available stacks | [optional] 
**AvailableSandboxes** | Pointer to **int32** | Number of available sandboxes | [optional] 

## Methods

### NewOrganization

`func NewOrganization(name string, id string, ownerId string, ) *Organization`

NewOrganization instantiates a new Organization object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationWithDefaults

`func NewOrganizationWithDefaults() *Organization`

NewOrganizationWithDefaults instantiates a new Organization object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Organization) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Organization) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Organization) SetName(v string)`

SetName sets Name field to given value.


### GetId

`func (o *Organization) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Organization) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Organization) SetId(v string)`

SetId sets Id field to given value.


### GetOwnerId

`func (o *Organization) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *Organization) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *Organization) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.


### GetAvailableStacks

`func (o *Organization) GetAvailableStacks() int32`

GetAvailableStacks returns the AvailableStacks field if non-nil, zero value otherwise.

### GetAvailableStacksOk

`func (o *Organization) GetAvailableStacksOk() (*int32, bool)`

GetAvailableStacksOk returns a tuple with the AvailableStacks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableStacks

`func (o *Organization) SetAvailableStacks(v int32)`

SetAvailableStacks sets AvailableStacks field to given value.

### HasAvailableStacks

`func (o *Organization) HasAvailableStacks() bool`

HasAvailableStacks returns a boolean if a field has been set.

### GetAvailableSandboxes

`func (o *Organization) GetAvailableSandboxes() int32`

GetAvailableSandboxes returns the AvailableSandboxes field if non-nil, zero value otherwise.

### GetAvailableSandboxesOk

`func (o *Organization) GetAvailableSandboxesOk() (*int32, bool)`

GetAvailableSandboxesOk returns a tuple with the AvailableSandboxes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableSandboxes

`func (o *Organization) SetAvailableSandboxes(v int32)`

SetAvailableSandboxes sets AvailableSandboxes field to given value.

### HasAvailableSandboxes

`func (o *Organization) HasAvailableSandboxes() bool`

HasAvailableSandboxes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


