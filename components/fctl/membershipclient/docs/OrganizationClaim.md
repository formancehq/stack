# OrganizationClaim

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**DisplayName** | **string** |  | 
**Stacks** | [**[]StackClaim**](StackClaim.md) |  | 

## Methods

### NewOrganizationClaim

`func NewOrganizationClaim(id string, name string, displayName string, stacks []StackClaim, ) *OrganizationClaim`

NewOrganizationClaim instantiates a new OrganizationClaim object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationClaimWithDefaults

`func NewOrganizationClaimWithDefaults() *OrganizationClaim`

NewOrganizationClaimWithDefaults instantiates a new OrganizationClaim object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OrganizationClaim) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OrganizationClaim) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OrganizationClaim) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *OrganizationClaim) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrganizationClaim) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrganizationClaim) SetName(v string)`

SetName sets Name field to given value.


### GetDisplayName

`func (o *OrganizationClaim) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *OrganizationClaim) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *OrganizationClaim) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetStacks

`func (o *OrganizationClaim) GetStacks() []StackClaim`

GetStacks returns the Stacks field if non-nil, zero value otherwise.

### GetStacksOk

`func (o *OrganizationClaim) GetStacksOk() (*[]StackClaim, bool)`

GetStacksOk returns a tuple with the Stacks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStacks

`func (o *OrganizationClaim) SetStacks(v []StackClaim)`

SetStacks sets Stacks field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


