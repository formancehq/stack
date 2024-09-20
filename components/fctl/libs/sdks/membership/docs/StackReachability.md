# StackReachability

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Reachable** | **bool** | Stack is reachable through Stargate | 
**LastReachableUpdate** | Pointer to **time.Time** | Last time the stack was reachable | [optional] 

## Methods

### NewStackReachability

`func NewStackReachability(reachable bool, ) *StackReachability`

NewStackReachability instantiates a new StackReachability object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackReachabilityWithDefaults

`func NewStackReachabilityWithDefaults() *StackReachability`

NewStackReachabilityWithDefaults instantiates a new StackReachability object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReachable

`func (o *StackReachability) GetReachable() bool`

GetReachable returns the Reachable field if non-nil, zero value otherwise.

### GetReachableOk

`func (o *StackReachability) GetReachableOk() (*bool, bool)`

GetReachableOk returns a tuple with the Reachable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReachable

`func (o *StackReachability) SetReachable(v bool)`

SetReachable sets Reachable field to given value.


### GetLastReachableUpdate

`func (o *StackReachability) GetLastReachableUpdate() time.Time`

GetLastReachableUpdate returns the LastReachableUpdate field if non-nil, zero value otherwise.

### GetLastReachableUpdateOk

`func (o *StackReachability) GetLastReachableUpdateOk() (*time.Time, bool)`

GetLastReachableUpdateOk returns a tuple with the LastReachableUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastReachableUpdate

`func (o *StackReachability) SetLastReachableUpdate(v time.Time)`

SetLastReachableUpdate sets LastReachableUpdate field to given value.

### HasLastReachableUpdate

`func (o *StackReachability) HasLastReachableUpdate() bool`

HasLastReachableUpdate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


