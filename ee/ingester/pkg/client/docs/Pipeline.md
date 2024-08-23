# Pipeline

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Module** | **string** |  | 
**ConnectorID** | **string** |  | 
**Id** | **string** |  | 
**State** | [**State**](State.md) |  | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewPipeline

`func NewPipeline(module string, connectorID string, id string, state State, createdAt time.Time, ) *Pipeline`

NewPipeline instantiates a new Pipeline object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineWithDefaults

`func NewPipelineWithDefaults() *Pipeline`

NewPipelineWithDefaults instantiates a new Pipeline object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetModule

`func (o *Pipeline) GetModule() string`

GetModule returns the Module field if non-nil, zero value otherwise.

### GetModuleOk

`func (o *Pipeline) GetModuleOk() (*string, bool)`

GetModuleOk returns a tuple with the Module field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModule

`func (o *Pipeline) SetModule(v string)`

SetModule sets Module field to given value.


### GetConnectorID

`func (o *Pipeline) GetConnectorID() string`

GetConnectorID returns the ConnectorID field if non-nil, zero value otherwise.

### GetConnectorIDOk

`func (o *Pipeline) GetConnectorIDOk() (*string, bool)`

GetConnectorIDOk returns a tuple with the ConnectorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorID

`func (o *Pipeline) SetConnectorID(v string)`

SetConnectorID sets ConnectorID field to given value.


### GetId

`func (o *Pipeline) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Pipeline) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Pipeline) SetId(v string)`

SetId sets Id field to given value.


### GetState

`func (o *Pipeline) GetState() State`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Pipeline) GetStateOk() (*State, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Pipeline) SetState(v State)`

SetState sets State field to given value.


### GetCreatedAt

`func (o *Pipeline) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Pipeline) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Pipeline) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


