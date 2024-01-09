# Invitation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**OrganizationId** | **string** |  | 
**UserEmail** | **string** |  | 
**Status** | **string** |  | 
**CreationDate** | **time.Time** |  | 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**Role** | [**Role**](Role.md) |  | [default to EMPTY]
**StackClaims** | Pointer to [**[]StackClaim**](StackClaim.md) |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 
**OrganizationAccess** | Pointer to [**OrganizationUser**](OrganizationUser.md) |  | [optional] 

## Methods

### NewInvitation

`func NewInvitation(id string, organizationId string, userEmail string, status string, creationDate time.Time, role Role, ) *Invitation`

NewInvitation instantiates a new Invitation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInvitationWithDefaults

`func NewInvitationWithDefaults() *Invitation`

NewInvitationWithDefaults instantiates a new Invitation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Invitation) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Invitation) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Invitation) SetId(v string)`

SetId sets Id field to given value.


### GetOrganizationId

`func (o *Invitation) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *Invitation) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *Invitation) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.


### GetUserEmail

`func (o *Invitation) GetUserEmail() string`

GetUserEmail returns the UserEmail field if non-nil, zero value otherwise.

### GetUserEmailOk

`func (o *Invitation) GetUserEmailOk() (*string, bool)`

GetUserEmailOk returns a tuple with the UserEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserEmail

`func (o *Invitation) SetUserEmail(v string)`

SetUserEmail sets UserEmail field to given value.


### GetStatus

`func (o *Invitation) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Invitation) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Invitation) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetCreationDate

`func (o *Invitation) GetCreationDate() time.Time`

GetCreationDate returns the CreationDate field if non-nil, zero value otherwise.

### GetCreationDateOk

`func (o *Invitation) GetCreationDateOk() (*time.Time, bool)`

GetCreationDateOk returns a tuple with the CreationDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreationDate

`func (o *Invitation) SetCreationDate(v time.Time)`

SetCreationDate sets CreationDate field to given value.


### GetUpdatedAt

`func (o *Invitation) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Invitation) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Invitation) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Invitation) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetRole

`func (o *Invitation) GetRole() Role`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *Invitation) GetRoleOk() (*Role, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *Invitation) SetRole(v Role)`

SetRole sets Role field to given value.


### GetStackClaims

`func (o *Invitation) GetStackClaims() []StackClaim`

GetStackClaims returns the StackClaims field if non-nil, zero value otherwise.

### GetStackClaimsOk

`func (o *Invitation) GetStackClaimsOk() (*[]StackClaim, bool)`

GetStackClaimsOk returns a tuple with the StackClaims field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackClaims

`func (o *Invitation) SetStackClaims(v []StackClaim)`

SetStackClaims sets StackClaims field to given value.

### HasStackClaims

`func (o *Invitation) HasStackClaims() bool`

HasStackClaims returns a boolean if a field has been set.

### GetUserId

`func (o *Invitation) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *Invitation) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *Invitation) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *Invitation) HasUserId() bool`

HasUserId returns a boolean if a field has been set.

### GetOrganizationAccess

`func (o *Invitation) GetOrganizationAccess() OrganizationUser`

GetOrganizationAccess returns the OrganizationAccess field if non-nil, zero value otherwise.

### GetOrganizationAccessOk

`func (o *Invitation) GetOrganizationAccessOk() (*OrganizationUser, bool)`

GetOrganizationAccessOk returns a tuple with the OrganizationAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationAccess

`func (o *Invitation) SetOrganizationAccess(v OrganizationUser)`

SetOrganizationAccess sets OrganizationAccess field to given value.

### HasOrganizationAccess

`func (o *Invitation) HasOrganizationAccess() bool`

HasOrganizationAccess returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


