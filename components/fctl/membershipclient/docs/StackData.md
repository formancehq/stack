# StackData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Stack name | 
**Metadata** | **map[string]string** |  | 
**Version** | Pointer to **string** | Supported only with agent version &gt;&#x3D; v0.7.0 | [optional] 

## Methods

### NewStackData

`func NewStackData(name string, metadata map[string]string, ) *StackData`

NewStackData instantiates a new StackData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackDataWithDefaults

`func NewStackDataWithDefaults() *StackData`

NewStackDataWithDefaults instantiates a new StackData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *StackData) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *StackData) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *StackData) SetName(v string)`

SetName sets Name field to given value.


### GetMetadata

`func (o *StackData) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *StackData) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *StackData) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.


### GetVersion

`func (o *StackData) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *StackData) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *StackData) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *StackData) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


