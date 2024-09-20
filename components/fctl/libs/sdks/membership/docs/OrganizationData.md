# OrganizationData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Organization name | 
**DefaultOrganizationAccess** | Pointer to [**Role**](Role.md) |  | [optional] 
**DefaultStackAccess** | Pointer to [**Role**](Role.md) |  | [optional] 
**Domain** | Pointer to **string** | Organization domain | [optional] 

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

`func (o *OrganizationData) GetDefaultOrganizationAccess() Role`

GetDefaultOrganizationAccess returns the DefaultOrganizationAccess field if non-nil, zero value otherwise.

### GetDefaultOrganizationAccessOk

`func (o *OrganizationData) GetDefaultOrganizationAccessOk() (*Role, bool)`

GetDefaultOrganizationAccessOk returns a tuple with the DefaultOrganizationAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultOrganizationAccess

`func (o *OrganizationData) SetDefaultOrganizationAccess(v Role)`

SetDefaultOrganizationAccess sets DefaultOrganizationAccess field to given value.

### HasDefaultOrganizationAccess

`func (o *OrganizationData) HasDefaultOrganizationAccess() bool`

HasDefaultOrganizationAccess returns a boolean if a field has been set.

### GetDefaultStackAccess

`func (o *OrganizationData) GetDefaultStackAccess() Role`

GetDefaultStackAccess returns the DefaultStackAccess field if non-nil, zero value otherwise.

### GetDefaultStackAccessOk

`func (o *OrganizationData) GetDefaultStackAccessOk() (*Role, bool)`

GetDefaultStackAccessOk returns a tuple with the DefaultStackAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultStackAccess

`func (o *OrganizationData) SetDefaultStackAccess(v Role)`

SetDefaultStackAccess sets DefaultStackAccess field to given value.

### HasDefaultStackAccess

`func (o *OrganizationData) HasDefaultStackAccess() bool`

HasDefaultStackAccess returns a boolean if a field has been set.

### GetDomain

`func (o *OrganizationData) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *OrganizationData) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *OrganizationData) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *OrganizationData) HasDomain() bool`

HasDomain returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


