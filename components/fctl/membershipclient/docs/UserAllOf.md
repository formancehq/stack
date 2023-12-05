# UserAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | User ID | 
**Role** | Pointer to **string** | User role | [optional] 
**IsAdmin** | **bool** | Is the user an admin of the organization | 

## Methods

### NewUserAllOf

`func NewUserAllOf(id string, isAdmin bool, ) *UserAllOf`

NewUserAllOf instantiates a new UserAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserAllOfWithDefaults

`func NewUserAllOfWithDefaults() *UserAllOf`

NewUserAllOfWithDefaults instantiates a new UserAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UserAllOf) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UserAllOf) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UserAllOf) SetId(v string)`

SetId sets Id field to given value.


### GetRole

`func (o *UserAllOf) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *UserAllOf) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *UserAllOf) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *UserAllOf) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetIsAdmin

`func (o *UserAllOf) GetIsAdmin() bool`

GetIsAdmin returns the IsAdmin field if non-nil, zero value otherwise.

### GetIsAdminOk

`func (o *UserAllOf) GetIsAdminOk() (*bool, bool)`

GetIsAdminOk returns a tuple with the IsAdmin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsAdmin

`func (o *UserAllOf) SetIsAdmin(v bool)`

SetIsAdmin sets IsAdmin field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


