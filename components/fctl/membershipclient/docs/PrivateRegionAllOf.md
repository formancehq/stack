# PrivateRegionAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientID** | **string** |  |
**OrganizationID** | **string** |  |
**CreatorID** | **string** |  |
**Secret** | Pointer to [**PrivateRegionAllOfSecret**](PrivateRegionAllOfSecret.md) |  | [optional]

## Methods

### NewPrivateRegionAllOf

`func NewPrivateRegionAllOf(clientID string, organizationID string, creatorID string, ) *PrivateRegionAllOf`

NewPrivateRegionAllOf instantiates a new PrivateRegionAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivateRegionAllOfWithDefaults

`func NewPrivateRegionAllOfWithDefaults() *PrivateRegionAllOf`

NewPrivateRegionAllOfWithDefaults instantiates a new PrivateRegionAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientID

`func (o *PrivateRegionAllOf) GetClientID() string`

GetClientID returns the ClientID field if non-nil, zero value otherwise.

### GetClientIDOk

`func (o *PrivateRegionAllOf) GetClientIDOk() (*string, bool)`

GetClientIDOk returns a tuple with the ClientID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientID

`func (o *PrivateRegionAllOf) SetClientID(v string)`

SetClientID sets ClientID field to given value.


### GetOrganizationID

`func (o *PrivateRegionAllOf) GetOrganizationID() string`

GetOrganizationID returns the OrganizationID field if non-nil, zero value otherwise.

### GetOrganizationIDOk

`func (o *PrivateRegionAllOf) GetOrganizationIDOk() (*string, bool)`

GetOrganizationIDOk returns a tuple with the OrganizationID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationID

`func (o *PrivateRegionAllOf) SetOrganizationID(v string)`

SetOrganizationID sets OrganizationID field to given value.


### GetCreatorID

`func (o *PrivateRegionAllOf) GetCreatorID() string`

GetCreatorID returns the CreatorID field if non-nil, zero value otherwise.

### GetCreatorIDOk

`func (o *PrivateRegionAllOf) GetCreatorIDOk() (*string, bool)`

GetCreatorIDOk returns a tuple with the CreatorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatorID

`func (o *PrivateRegionAllOf) SetCreatorID(v string)`

SetCreatorID sets CreatorID field to given value.


### GetSecret

`func (o *PrivateRegionAllOf) GetSecret() PrivateRegionAllOfSecret`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *PrivateRegionAllOf) GetSecretOk() (*PrivateRegionAllOfSecret, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *PrivateRegionAllOf) SetSecret(v PrivateRegionAllOfSecret)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *PrivateRegionAllOf) HasSecret() bool`

HasSecret returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
