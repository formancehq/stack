# ConnectorsResponseDataInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Provider** | Pointer to [**Connector**](Connector.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Disabled** | Pointer to **bool** |  | [optional] 

## Methods

### NewConnectorsResponseDataInner

`func NewConnectorsResponseDataInner() *ConnectorsResponseDataInner`

NewConnectorsResponseDataInner instantiates a new ConnectorsResponseDataInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorsResponseDataInnerWithDefaults

`func NewConnectorsResponseDataInnerWithDefaults() *ConnectorsResponseDataInner`

NewConnectorsResponseDataInnerWithDefaults instantiates a new ConnectorsResponseDataInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProvider

`func (o *ConnectorsResponseDataInner) GetProvider() Connector`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *ConnectorsResponseDataInner) GetProviderOk() (*Connector, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *ConnectorsResponseDataInner) SetProvider(v Connector)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *ConnectorsResponseDataInner) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetEnabled

`func (o *ConnectorsResponseDataInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ConnectorsResponseDataInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ConnectorsResponseDataInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ConnectorsResponseDataInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetDisabled

`func (o *ConnectorsResponseDataInner) GetDisabled() bool`

GetDisabled returns the Disabled field if non-nil, zero value otherwise.

### GetDisabledOk

`func (o *ConnectorsResponseDataInner) GetDisabledOk() (*bool, bool)`

GetDisabledOk returns a tuple with the Disabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabled

`func (o *ConnectorsResponseDataInner) SetDisabled(v bool)`

SetDisabled sets Disabled field to given value.

### HasDisabled

`func (o *ConnectorsResponseDataInner) HasDisabled() bool`

HasDisabled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


