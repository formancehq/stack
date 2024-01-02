# InvitationClaim

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Roles** | **[]string** | User roles | 
**StackClaims** | [**[]StackClaimsInner**](StackClaimsInner.md) |  | 

## Methods

### NewInvitationClaim

`func NewInvitationClaim(roles []string, stackClaims []StackClaimsInner, ) *InvitationClaim`

NewInvitationClaim instantiates a new InvitationClaim object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInvitationClaimWithDefaults

`func NewInvitationClaimWithDefaults() *InvitationClaim`

NewInvitationClaimWithDefaults instantiates a new InvitationClaim object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRoles

`func (o *InvitationClaim) GetRoles() []string`

GetRoles returns the Roles field if non-nil, zero value otherwise.

### GetRolesOk

`func (o *InvitationClaim) GetRolesOk() (*[]string, bool)`

GetRolesOk returns a tuple with the Roles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoles

`func (o *InvitationClaim) SetRoles(v []string)`

SetRoles sets Roles field to given value.


### GetStackClaims

`func (o *InvitationClaim) GetStackClaims() []StackClaimsInner`

GetStackClaims returns the StackClaims field if non-nil, zero value otherwise.

### GetStackClaimsOk

`func (o *InvitationClaim) GetStackClaimsOk() (*[]StackClaimsInner, bool)`

GetStackClaimsOk returns a tuple with the StackClaims field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackClaims

`func (o *InvitationClaim) SetStackClaims(v []StackClaimsInner)`

SetStackClaims sets StackClaims field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


