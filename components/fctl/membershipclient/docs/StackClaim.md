# StackClaim

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StackId** | **string** |  | 
**Roles** | **[]string** | User roles | 

## Methods

### NewStackClaim

`func NewStackClaim(stackId string, roles []string, ) *StackClaim`

NewStackClaim instantiates a new StackClaim object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackClaimWithDefaults

`func NewStackClaimWithDefaults() *StackClaim`

NewStackClaimWithDefaults instantiates a new StackClaim object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStackId

`func (o *StackClaim) GetStackId() string`

GetStackId returns the StackId field if non-nil, zero value otherwise.

### GetStackIdOk

`func (o *StackClaim) GetStackIdOk() (*string, bool)`

GetStackIdOk returns a tuple with the StackId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackId

`func (o *StackClaim) SetStackId(v string)`

SetStackId sets StackId field to given value.


### GetRoles

`func (o *StackClaim) GetRoles() []string`

GetRoles returns the Roles field if non-nil, zero value otherwise.

### GetRolesOk

`func (o *StackClaim) GetRolesOk() (*[]string, bool)`

GetRolesOk returns a tuple with the Roles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoles

`func (o *StackClaim) SetRoles(v []string)`

SetRoles sets Roles field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


