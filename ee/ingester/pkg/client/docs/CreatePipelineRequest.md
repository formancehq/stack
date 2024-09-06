# CreatePipelineRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Module** | **string** |  | 
**ConnectorID** | **string** |  | 

## Methods

### NewCreatePipelineRequest

`func NewCreatePipelineRequest(module string, connectorID string, ) *CreatePipelineRequest`

NewCreatePipelineRequest instantiates a new CreatePipelineRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePipelineRequestWithDefaults

`func NewCreatePipelineRequestWithDefaults() *CreatePipelineRequest`

NewCreatePipelineRequestWithDefaults instantiates a new CreatePipelineRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetModule

`func (o *CreatePipelineRequest) GetModule() string`

GetModule returns the Module field if non-nil, zero value otherwise.

### GetModuleOk

`func (o *CreatePipelineRequest) GetModuleOk() (*string, bool)`

GetModuleOk returns a tuple with the Module field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModule

`func (o *CreatePipelineRequest) SetModule(v string)`

SetModule sets Module field to given value.


### GetConnectorID

`func (o *CreatePipelineRequest) GetConnectorID() string`

GetConnectorID returns the ConnectorID field if non-nil, zero value otherwise.

### GetConnectorIDOk

`func (o *CreatePipelineRequest) GetConnectorIDOk() (*string, bool)`

GetConnectorIDOk returns a tuple with the ConnectorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorID

`func (o *CreatePipelineRequest) SetConnectorID(v string)`

SetConnectorID sets ConnectorID field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


