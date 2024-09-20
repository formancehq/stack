# ServerInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | **string** |  | 
**Capabilities** | Pointer to [**[]Capability**](Capability.md) |  | [optional] 
**ConsoleURL** | Pointer to **string** |  | [optional] 

## Methods

### NewServerInfo

`func NewServerInfo(version string, ) *ServerInfo`

NewServerInfo instantiates a new ServerInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerInfoWithDefaults

`func NewServerInfoWithDefaults() *ServerInfo`

NewServerInfoWithDefaults instantiates a new ServerInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *ServerInfo) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ServerInfo) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ServerInfo) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetCapabilities

`func (o *ServerInfo) GetCapabilities() []Capability`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *ServerInfo) GetCapabilitiesOk() (*[]Capability, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *ServerInfo) SetCapabilities(v []Capability)`

SetCapabilities sets Capabilities field to given value.

### HasCapabilities

`func (o *ServerInfo) HasCapabilities() bool`

HasCapabilities returns a boolean if a field has been set.

### GetConsoleURL

`func (o *ServerInfo) GetConsoleURL() string`

GetConsoleURL returns the ConsoleURL field if non-nil, zero value otherwise.

### GetConsoleURLOk

`func (o *ServerInfo) GetConsoleURLOk() (*string, bool)`

GetConsoleURLOk returns a tuple with the ConsoleURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsoleURL

`func (o *ServerInfo) SetConsoleURL(v string)`

SetConsoleURL sets ConsoleURL field to given value.

### HasConsoleURL

`func (o *ServerInfo) HasConsoleURL() bool`

HasConsoleURL returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


