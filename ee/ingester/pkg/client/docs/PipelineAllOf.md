# PipelineAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**State** | [**State**](State.md) |  | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewPipelineAllOf

`func NewPipelineAllOf(id string, state State, createdAt time.Time, ) *PipelineAllOf`

NewPipelineAllOf instantiates a new PipelineAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineAllOfWithDefaults

`func NewPipelineAllOfWithDefaults() *PipelineAllOf`

NewPipelineAllOfWithDefaults instantiates a new PipelineAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PipelineAllOf) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PipelineAllOf) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PipelineAllOf) SetId(v string)`

SetId sets Id field to given value.


### GetState

`func (o *PipelineAllOf) GetState() State`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *PipelineAllOf) GetStateOk() (*State, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *PipelineAllOf) SetState(v State)`

SetState sets State field to given value.


### GetCreatedAt

`func (o *PipelineAllOf) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PipelineAllOf) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PipelineAllOf) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


