# StackVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | Pointer to **string** | Supported only with agent version &gt;&#x3D; v0.7.0 | [optional] 

## Methods

### NewStackVersion

`func NewStackVersion() *StackVersion`

NewStackVersion instantiates a new StackVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackVersionWithDefaults

`func NewStackVersionWithDefaults() *StackVersion`

NewStackVersionWithDefaults instantiates a new StackVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *StackVersion) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *StackVersion) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *StackVersion) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *StackVersion) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


