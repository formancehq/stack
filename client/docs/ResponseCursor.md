# ResponseCursor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PageSize** | Pointer to **interface{}** |  | [optional] 
**HasMore** | Pointer to **interface{}** |  | [optional] 
**Total** | Pointer to [**ResponseCursorTotal**](ResponseCursorTotal.md) |  | [optional] 
**Next** | Pointer to **interface{}** |  | [optional] 
**Previous** | Pointer to **interface{}** |  | [optional] 
**Data** | Pointer to **interface{}** |  | [optional] 

## Methods

### NewResponseCursor

`func NewResponseCursor() *ResponseCursor`

NewResponseCursor instantiates a new ResponseCursor object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResponseCursorWithDefaults

`func NewResponseCursorWithDefaults() *ResponseCursor`

NewResponseCursorWithDefaults instantiates a new ResponseCursor object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPageSize

`func (o *ResponseCursor) GetPageSize() interface{}`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *ResponseCursor) GetPageSizeOk() (*interface{}, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *ResponseCursor) SetPageSize(v interface{})`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *ResponseCursor) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### SetPageSizeNil

`func (o *ResponseCursor) SetPageSizeNil(b bool)`

 SetPageSizeNil sets the value for PageSize to be an explicit nil

### UnsetPageSize
`func (o *ResponseCursor) UnsetPageSize()`

UnsetPageSize ensures that no value is present for PageSize, not even an explicit nil
### GetHasMore

`func (o *ResponseCursor) GetHasMore() interface{}`

GetHasMore returns the HasMore field if non-nil, zero value otherwise.

### GetHasMoreOk

`func (o *ResponseCursor) GetHasMoreOk() (*interface{}, bool)`

GetHasMoreOk returns a tuple with the HasMore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasMore

`func (o *ResponseCursor) SetHasMore(v interface{})`

SetHasMore sets HasMore field to given value.

### HasHasMore

`func (o *ResponseCursor) HasHasMore() bool`

HasHasMore returns a boolean if a field has been set.

### SetHasMoreNil

`func (o *ResponseCursor) SetHasMoreNil(b bool)`

 SetHasMoreNil sets the value for HasMore to be an explicit nil

### UnsetHasMore
`func (o *ResponseCursor) UnsetHasMore()`

UnsetHasMore ensures that no value is present for HasMore, not even an explicit nil
### GetTotal

`func (o *ResponseCursor) GetTotal() ResponseCursorTotal`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *ResponseCursor) GetTotalOk() (*ResponseCursorTotal, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *ResponseCursor) SetTotal(v ResponseCursorTotal)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *ResponseCursor) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetNext

`func (o *ResponseCursor) GetNext() interface{}`

GetNext returns the Next field if non-nil, zero value otherwise.

### GetNextOk

`func (o *ResponseCursor) GetNextOk() (*interface{}, bool)`

GetNextOk returns a tuple with the Next field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNext

`func (o *ResponseCursor) SetNext(v interface{})`

SetNext sets Next field to given value.

### HasNext

`func (o *ResponseCursor) HasNext() bool`

HasNext returns a boolean if a field has been set.

### SetNextNil

`func (o *ResponseCursor) SetNextNil(b bool)`

 SetNextNil sets the value for Next to be an explicit nil

### UnsetNext
`func (o *ResponseCursor) UnsetNext()`

UnsetNext ensures that no value is present for Next, not even an explicit nil
### GetPrevious

`func (o *ResponseCursor) GetPrevious() interface{}`

GetPrevious returns the Previous field if non-nil, zero value otherwise.

### GetPreviousOk

`func (o *ResponseCursor) GetPreviousOk() (*interface{}, bool)`

GetPreviousOk returns a tuple with the Previous field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrevious

`func (o *ResponseCursor) SetPrevious(v interface{})`

SetPrevious sets Previous field to given value.

### HasPrevious

`func (o *ResponseCursor) HasPrevious() bool`

HasPrevious returns a boolean if a field has been set.

### SetPreviousNil

`func (o *ResponseCursor) SetPreviousNil(b bool)`

 SetPreviousNil sets the value for Previous to be an explicit nil

### UnsetPrevious
`func (o *ResponseCursor) UnsetPrevious()`

UnsetPrevious ensures that no value is present for Previous, not even an explicit nil
### GetData

`func (o *ResponseCursor) GetData() interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ResponseCursor) GetDataOk() (*interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ResponseCursor) SetData(v interface{})`

SetData sets Data field to given value.

### HasData

`func (o *ResponseCursor) HasData() bool`

HasData returns a boolean if a field has been set.

### SetDataNil

`func (o *ResponseCursor) SetDataNil(b bool)`

 SetDataNil sets the value for Data to be an explicit nil

### UnsetData
`func (o *ResponseCursor) UnsetData()`

UnsetData ensures that no value is present for Data, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


