# OrganizationData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Organization name | 
**DefaultOrganizationAccess** | Pointer to **[]string** |  | [optional] 
**DefaultStackAccess** | Pointer to **[]string** |  | [optional] 

## Methods

### NewOrganizationData

`func NewOrganizationData(name string, ) *OrganizationData`

NewOrganizationData instantiates a new OrganizationData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationDataWithDefaults

`func NewOrganizationDataWithDefaults() *OrganizationData`

NewOrganizationDataWithDefaults instantiates a new OrganizationData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *OrganizationData) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrganizationData) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrganizationData) SetName(v string)`

SetName sets Name field to given value.


### GetDefaultOrganizationAccess

`func (o *OrganizationData) GetDefaultOrganizationAccess() []string`

GetDefaultOrganizationAccess returns the DefaultOrganizationAccess field if non-nil, zero value otherwise.

### GetDefaultOrganizationAccessOk

`func (o *OrganizationData) GetDefaultOrganizationAccessOk() (*[]string, bool)`

GetDefaultOrganizationAccessOk returns a tuple with the DefaultOrganizationAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultOrganizationAccess

`func (o *OrganizationData) SetDefaultOrganizationAccess(v []string)`

SetDefaultOrganizationAccess sets DefaultOrganizationAccess field to given value.

### HasDefaultOrganizationAccess

`func (o *OrganizationData) HasDefaultOrganizationAccess() bool`

HasDefaultOrganizationAccess returns a boolean if a field has been set.

### GetDefaultStackAccess

`func (o *OrganizationData) GetDefaultStackAccess() []string`

GetDefaultStackAccess returns the DefaultStackAccess field if non-nil, zero value otherwise.

### GetDefaultStackAccessOk

`func (o *OrganizationData) GetDefaultStackAccessOk() (*[]string, bool)`

GetDefaultStackAccessOk returns a tuple with the DefaultStackAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultStackAccess

`func (o *OrganizationData) SetDefaultStackAccess(v []string)`

SetDefaultStackAccess sets DefaultStackAccess field to given value.

### HasDefaultStackAccess

`func (o *OrganizationData) HasDefaultStackAccess() bool`

HasDefaultStackAccess returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


