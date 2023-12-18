# CreateStackRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Stack name | 
**Metadata** | **map[string]string** |  | 
**Version** | Pointer to **string** | Supported only with agent version &gt;&#x3D; v0.7.0 | [optional] 
**RegionID** | **string** |  | 

## Methods

### NewCreateStackRequest

`func NewCreateStackRequest(name string, metadata map[string]string, regionID string, ) *CreateStackRequest`

NewCreateStackRequest instantiates a new CreateStackRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateStackRequestWithDefaults

`func NewCreateStackRequestWithDefaults() *CreateStackRequest`

NewCreateStackRequestWithDefaults instantiates a new CreateStackRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateStackRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateStackRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateStackRequest) SetName(v string)`

SetName sets Name field to given value.


### GetMetadata

`func (o *CreateStackRequest) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *CreateStackRequest) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *CreateStackRequest) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.


### GetVersion

`func (o *CreateStackRequest) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *CreateStackRequest) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *CreateStackRequest) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *CreateStackRequest) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetRegionID

`func (o *CreateStackRequest) GetRegionID() string`

GetRegionID returns the RegionID field if non-nil, zero value otherwise.

### GetRegionIDOk

`func (o *CreateStackRequest) GetRegionIDOk() (*string, bool)`

GetRegionIDOk returns a tuple with the RegionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionID

`func (o *CreateStackRequest) SetRegionID(v string)`

SetRegionID sets RegionID field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


