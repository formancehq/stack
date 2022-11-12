# ResponseCursor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PageSize** | Pointer to **float32** |  | [optional] 
**HasMore** | Pointer to **bool** |  | [optional] 
**Total** | Pointer to [**ResponseCursorTotal**](ResponseCursorTotal.md) |  | [optional] 
**Next** | Pointer to **string** |  | [optional] 
**Previous** | Pointer to **string** |  | [optional] 
**Data** | Pointer to **[]map[string]interface{}** |  | [optional] 

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

`func (o *ResponseCursor) GetPageSize() float32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *ResponseCursor) GetPageSizeOk() (*float32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *ResponseCursor) SetPageSize(v float32)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *ResponseCursor) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### GetHasMore

`func (o *ResponseCursor) GetHasMore() bool`

GetHasMore returns the HasMore field if non-nil, zero value otherwise.

### GetHasMoreOk

`func (o *ResponseCursor) GetHasMoreOk() (*bool, bool)`

GetHasMoreOk returns a tuple with the HasMore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasMore

`func (o *ResponseCursor) SetHasMore(v bool)`

SetHasMore sets HasMore field to given value.

### HasHasMore

`func (o *ResponseCursor) HasHasMore() bool`

HasHasMore returns a boolean if a field has been set.

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

`func (o *ResponseCursor) GetNext() string`

GetNext returns the Next field if non-nil, zero value otherwise.

### GetNextOk

`func (o *ResponseCursor) GetNextOk() (*string, bool)`

GetNextOk returns a tuple with the Next field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNext

`func (o *ResponseCursor) SetNext(v string)`

SetNext sets Next field to given value.

### HasNext

`func (o *ResponseCursor) HasNext() bool`

HasNext returns a boolean if a field has been set.

### GetPrevious

`func (o *ResponseCursor) GetPrevious() string`

GetPrevious returns the Previous field if non-nil, zero value otherwise.

### GetPreviousOk

`func (o *ResponseCursor) GetPreviousOk() (*string, bool)`

GetPreviousOk returns a tuple with the Previous field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrevious

`func (o *ResponseCursor) SetPrevious(v string)`

SetPrevious sets Previous field to given value.

### HasPrevious

`func (o *ResponseCursor) HasPrevious() bool`

HasPrevious returns a boolean if a field has been set.

### GetData

`func (o *ResponseCursor) GetData() []map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ResponseCursor) GetDataOk() (*[]map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ResponseCursor) SetData(v []map[string]interface{})`

SetData sets Data field to given value.

### HasData

`func (o *ResponseCursor) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


