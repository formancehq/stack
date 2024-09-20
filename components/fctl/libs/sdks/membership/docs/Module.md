# Module

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**State** | **string** |  | 
**Status** | **string** |  | 
**LastStatusUpdate** | **time.Time** |  | 
**LastStateUpdate** | **time.Time** |  | 

## Methods

### NewModule

`func NewModule(name string, state string, status string, lastStatusUpdate time.Time, lastStateUpdate time.Time, ) *Module`

NewModule instantiates a new Module object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewModuleWithDefaults

`func NewModuleWithDefaults() *Module`

NewModuleWithDefaults instantiates a new Module object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Module) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Module) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Module) SetName(v string)`

SetName sets Name field to given value.


### GetState

`func (o *Module) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Module) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Module) SetState(v string)`

SetState sets State field to given value.


### GetStatus

`func (o *Module) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Module) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Module) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetLastStatusUpdate

`func (o *Module) GetLastStatusUpdate() time.Time`

GetLastStatusUpdate returns the LastStatusUpdate field if non-nil, zero value otherwise.

### GetLastStatusUpdateOk

`func (o *Module) GetLastStatusUpdateOk() (*time.Time, bool)`

GetLastStatusUpdateOk returns a tuple with the LastStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStatusUpdate

`func (o *Module) SetLastStatusUpdate(v time.Time)`

SetLastStatusUpdate sets LastStatusUpdate field to given value.


### GetLastStateUpdate

`func (o *Module) GetLastStateUpdate() time.Time`

GetLastStateUpdate returns the LastStateUpdate field if non-nil, zero value otherwise.

### GetLastStateUpdateOk

`func (o *Module) GetLastStateUpdateOk() (*time.Time, bool)`

GetLastStateUpdateOk returns a tuple with the LastStateUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStateUpdate

`func (o *Module) SetLastStateUpdate(v time.Time)`

SetLastStateUpdate sets LastStateUpdate field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


