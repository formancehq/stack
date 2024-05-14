# UserInfoResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**EmailVerified** | Pointer to **bool** |  | [optional] 
**Sub** | **string** |  | 
**Org** | [**[]OrganizationClaim**](OrganizationClaim.md) |  | 

## Methods

### NewUserInfoResponse

`func NewUserInfoResponse(email string, sub string, org []OrganizationClaim, ) *UserInfoResponse`

NewUserInfoResponse instantiates a new UserInfoResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserInfoResponseWithDefaults

`func NewUserInfoResponseWithDefaults() *UserInfoResponse`

NewUserInfoResponseWithDefaults instantiates a new UserInfoResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *UserInfoResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserInfoResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserInfoResponse) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetEmailVerified

`func (o *UserInfoResponse) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *UserInfoResponse) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *UserInfoResponse) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.

### HasEmailVerified

`func (o *UserInfoResponse) HasEmailVerified() bool`

HasEmailVerified returns a boolean if a field has been set.

### GetSub

`func (o *UserInfoResponse) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *UserInfoResponse) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSub

`func (o *UserInfoResponse) SetSub(v string)`

SetSub sets Sub field to given value.


### GetOrg

`func (o *UserInfoResponse) GetOrg() []OrganizationClaim`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *UserInfoResponse) GetOrgOk() (*[]OrganizationClaim, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *UserInfoResponse) SetOrg(v []OrganizationClaim)`

SetOrg sets Org field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


