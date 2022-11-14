# Cursor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**HasMore** | Pointer to **bool** |  | [optional] 

## Methods

### NewCursor

`func NewCursor() *Cursor`

NewCursor instantiates a new Cursor object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCursorWithDefaults

`func NewCursorWithDefaults() *Cursor`

NewCursorWithDefaults instantiates a new Cursor object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHasMore

`func (o *Cursor) GetHasMore() bool`

GetHasMore returns the HasMore field if non-nil, zero value otherwise.

### GetHasMoreOk

`func (o *Cursor) GetHasMoreOk() (*bool, bool)`

GetHasMoreOk returns a tuple with the HasMore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasMore

`func (o *Cursor) SetHasMore(v bool)`

SetHasMore sets HasMore field to given value.

### HasHasMore

`func (o *Cursor) HasHasMore() bool`

HasHasMore returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


