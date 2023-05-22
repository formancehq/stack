# GetVersionsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Region** | **string** |  | 
**Env** | **string** |  | 
**Versions** | [**[]Version**](Version.md) |  | 

## Methods

### NewGetVersionsResponse

`func NewGetVersionsResponse(region string, env string, versions []Version, ) *GetVersionsResponse`

NewGetVersionsResponse instantiates a new GetVersionsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetVersionsResponseWithDefaults

`func NewGetVersionsResponseWithDefaults() *GetVersionsResponse`

NewGetVersionsResponseWithDefaults instantiates a new GetVersionsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegion

`func (o *GetVersionsResponse) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *GetVersionsResponse) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *GetVersionsResponse) SetRegion(v string)`

SetRegion sets Region field to given value.


### GetEnv

`func (o *GetVersionsResponse) GetEnv() string`

GetEnv returns the Env field if non-nil, zero value otherwise.

### GetEnvOk

`func (o *GetVersionsResponse) GetEnvOk() (*string, bool)`

GetEnvOk returns a tuple with the Env field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnv

`func (o *GetVersionsResponse) SetEnv(v string)`

SetEnv sets Env field to given value.


### GetVersions

`func (o *GetVersionsResponse) GetVersions() []Version`

GetVersions returns the Versions field if non-nil, zero value otherwise.

### GetVersionsOk

`func (o *GetVersionsResponse) GetVersionsOk() (*[]Version, bool)`

GetVersionsOk returns a tuple with the Versions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersions

`func (o *GetVersionsResponse) SetVersions(v []Version)`

SetVersions sets Versions field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


