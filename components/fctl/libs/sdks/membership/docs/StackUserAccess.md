# StackUserAccess

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StackId** | **string** | Stack ID | 
**UserId** | **string** | User ID | 
**Role** | [**Role**](Role.md) |  | 

## Methods

### NewStackUserAccess

`func NewStackUserAccess(stackId string, userId string, role Role, ) *StackUserAccess`

NewStackUserAccess instantiates a new StackUserAccess object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackUserAccessWithDefaults

`func NewStackUserAccessWithDefaults() *StackUserAccess`

NewStackUserAccessWithDefaults instantiates a new StackUserAccess object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

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


### GetRole

`func (o *StackUserAccess) GetRole() Role`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *StackUserAccess) GetRoleOk() (*Role, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *StackUserAccess) SetRole(v Role)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


