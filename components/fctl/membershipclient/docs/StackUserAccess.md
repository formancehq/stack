# StackUserAccess

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Roles** | Pointer to **[]string** |  | [optional] 
**StackId** | **string** | Stack ID | 
**UserId** | **string** | User ID | 

## Methods

### NewStackUserAccess

`func NewStackUserAccess(stackId string, userId string, ) *StackUserAccess`

NewStackUserAccess instantiates a new StackUserAccess object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackUserAccessWithDefaults

`func NewStackUserAccessWithDefaults() *StackUserAccess`

NewStackUserAccessWithDefaults instantiates a new StackUserAccess object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRoles

`func (o *StackUserAccess) GetRoles() []string`

GetRoles returns the Roles field if non-nil, zero value otherwise.

### GetRolesOk

`func (o *StackUserAccess) GetRolesOk() (*[]string, bool)`

GetRolesOk returns a tuple with the Roles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoles

`func (o *StackUserAccess) SetRoles(v []string)`

SetRoles sets Roles field to given value.

### HasRoles

`func (o *StackUserAccess) HasRoles() bool`

HasRoles returns a boolean if a field has been set.

### GetStackId

`func (o *StackUserAccess) GetStackId() string`

GetStackId returns the StackId field if non-nil, zero value otherwise.

### GetStackIdOk

`func (o *StackUserAccess) GetStackIdOk() (*string, bool)`

GetStackIdOk returns a tuple with the StackId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackId

`func (o *StackUserAccess) SetStackId(v string)`

SetStackId sets StackId field to given value.


### GetUserId

`func (o *StackUserAccess) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *StackUserAccess) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *StackUserAccess) SetUserId(v string)`

SetUserId sets UserId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


