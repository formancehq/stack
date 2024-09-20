# InvitationClaim

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Role** | Pointer to [**Role**](Role.md) |  | [optional] 
**StackClaims** | Pointer to [**[]StackClaim**](StackClaim.md) |  | [optional] 

## Methods

### NewInvitationClaim

`func NewInvitationClaim() *InvitationClaim`

NewInvitationClaim instantiates a new InvitationClaim object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInvitationClaimWithDefaults

`func NewInvitationClaimWithDefaults() *InvitationClaim`

NewInvitationClaimWithDefaults instantiates a new InvitationClaim object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRole

`func (o *InvitationClaim) GetRole() Role`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *InvitationClaim) GetRoleOk() (*Role, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *InvitationClaim) SetRole(v Role)`

SetRole sets Role field to given value.

### HasRole

`func (o *InvitationClaim) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetStackClaims

`func (o *InvitationClaim) GetStackClaims() []StackClaim`

GetStackClaims returns the StackClaims field if non-nil, zero value otherwise.

### GetStackClaimsOk

`func (o *InvitationClaim) GetStackClaimsOk() (*[]StackClaim, bool)`

GetStackClaimsOk returns a tuple with the StackClaims field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackClaims

`func (o *InvitationClaim) SetStackClaims(v []StackClaim)`

SetStackClaims sets StackClaims field to given value.

### HasStackClaims

`func (o *InvitationClaim) HasStackClaims() bool`

HasStackClaims returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


