# Scope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Label** | **string** |  | 
**Id** | **string** |  | 
**Transient** | Pointer to **[]string** |  | [optional] 

## Methods

### NewScope

`func NewScope(label string, id string, ) *Scope`

NewScope instantiates a new Scope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewScopeWithDefaults

`func NewScopeWithDefaults() *Scope`

NewScopeWithDefaults instantiates a new Scope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabel

`func (o *Scope) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *Scope) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *Scope) SetLabel(v string)`

SetLabel sets Label field to given value.


### GetId

`func (o *Scope) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Scope) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Scope) SetId(v string)`

SetId sets Id field to given value.


### GetTransient

`func (o *Scope) GetTransient() []string`

GetTransient returns the Transient field if non-nil, zero value otherwise.

### GetTransientOk

`func (o *Scope) GetTransientOk() (*[]string, bool)`

GetTransientOk returns a tuple with the Transient field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransient

`func (o *Scope) SetTransient(v []string)`

SetTransient sets Transient field to given value.

### HasTransient

`func (o *Scope) HasTransient() bool`

HasTransient returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


