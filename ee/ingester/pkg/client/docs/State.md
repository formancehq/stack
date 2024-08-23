# State

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Label** | **string** |  | 
**Cursor** | Pointer to **string** |  | [optional] 
**PreviousState** | Pointer to [**State**](State.md) |  | [optional] 

## Methods

### NewState

`func NewState(label string, ) *State`

NewState instantiates a new State object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStateWithDefaults

`func NewStateWithDefaults() *State`

NewStateWithDefaults instantiates a new State object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabel

`func (o *State) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *State) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *State) SetLabel(v string)`

SetLabel sets Label field to given value.


### GetCursor

`func (o *State) GetCursor() string`

GetCursor returns the Cursor field if non-nil, zero value otherwise.

### GetCursorOk

`func (o *State) GetCursorOk() (*string, bool)`

GetCursorOk returns a tuple with the Cursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCursor

`func (o *State) SetCursor(v string)`

SetCursor sets Cursor field to given value.

### HasCursor

`func (o *State) HasCursor() bool`

HasCursor returns a boolean if a field has been set.

### GetPreviousState

`func (o *State) GetPreviousState() State`

GetPreviousState returns the PreviousState field if non-nil, zero value otherwise.

### GetPreviousStateOk

`func (o *State) GetPreviousStateOk() (*State, bool)`

GetPreviousStateOk returns a tuple with the PreviousState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreviousState

`func (o *State) SetPreviousState(v State)`

SetPreviousState sets PreviousState field to given value.

### HasPreviousState

`func (o *State) HasPreviousState() bool`

HasPreviousState returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


