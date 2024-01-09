# OrganizationUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Role** | [**Role**](Role.md) |  | [default to EMPTY]
**Email** | **string** |  | 
**Id** | **string** | User ID | 

## Methods

### NewOrganizationUser

`func NewOrganizationUser(role Role, email string, id string, ) *OrganizationUser`

NewOrganizationUser instantiates a new OrganizationUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationUserWithDefaults

`func NewOrganizationUserWithDefaults() *OrganizationUser`

NewOrganizationUserWithDefaults instantiates a new OrganizationUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRole

`func (o *OrganizationUser) GetRole() Role`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *OrganizationUser) GetRoleOk() (*Role, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *OrganizationUser) SetRole(v Role)`

SetRole sets Role field to given value.


### GetEmail

`func (o *OrganizationUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *OrganizationUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *OrganizationUser) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetId

`func (o *OrganizationUser) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OrganizationUser) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OrganizationUser) SetId(v string)`

SetId sets Id field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


