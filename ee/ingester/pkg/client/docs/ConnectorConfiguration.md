# ConnectorConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Driver** | **string** |  | 
**Config** | **map[string]interface{}** |  | 

## Methods

### NewConnectorConfiguration

`func NewConnectorConfiguration(driver string, config map[string]interface{}, ) *ConnectorConfiguration`

NewConnectorConfiguration instantiates a new ConnectorConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorConfigurationWithDefaults

`func NewConnectorConfigurationWithDefaults() *ConnectorConfiguration`

NewConnectorConfigurationWithDefaults instantiates a new ConnectorConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDriver

`func (o *ConnectorConfiguration) GetDriver() string`

GetDriver returns the Driver field if non-nil, zero value otherwise.

### GetDriverOk

`func (o *ConnectorConfiguration) GetDriverOk() (*string, bool)`

GetDriverOk returns a tuple with the Driver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDriver

`func (o *ConnectorConfiguration) SetDriver(v string)`

SetDriver sets Driver field to given value.


### GetConfig

`func (o *ConnectorConfiguration) GetConfig() map[string]interface{}`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *ConnectorConfiguration) GetConfigOk() (*map[string]interface{}, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *ConnectorConfiguration) SetConfig(v map[string]interface{})`

SetConfig sets Config field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


