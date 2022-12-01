# GetManyConfigs200ResponseCursor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**HasMore** | Pointer to **bool** |  | [optional] 
**Data** | [**[]ConfigActivated**](ConfigActivated.md) |  | 

## Methods

### NewGetManyConfigs200ResponseCursor

`func NewGetManyConfigs200ResponseCursor(data []ConfigActivated, ) *GetManyConfigs200ResponseCursor`

NewGetManyConfigs200ResponseCursor instantiates a new GetManyConfigs200ResponseCursor object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetManyConfigs200ResponseCursorWithDefaults

`func NewGetManyConfigs200ResponseCursorWithDefaults() *GetManyConfigs200ResponseCursor`

NewGetManyConfigs200ResponseCursorWithDefaults instantiates a new GetManyConfigs200ResponseCursor object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHasMore

`func (o *GetManyConfigs200ResponseCursor) GetHasMore() bool`

GetHasMore returns the HasMore field if non-nil, zero value otherwise.

### GetHasMoreOk

`func (o *GetManyConfigs200ResponseCursor) GetHasMoreOk() (*bool, bool)`

GetHasMoreOk returns a tuple with the HasMore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasMore

`func (o *GetManyConfigs200ResponseCursor) SetHasMore(v bool)`

SetHasMore sets HasMore field to given value.

### HasHasMore

`func (o *GetManyConfigs200ResponseCursor) HasHasMore() bool`

HasHasMore returns a boolean if a field has been set.

### GetData

`func (o *GetManyConfigs200ResponseCursor) GetData() []ConfigActivated`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *GetManyConfigs200ResponseCursor) GetDataOk() (*[]ConfigActivated, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *GetManyConfigs200ResponseCursor) SetData(v []ConfigActivated)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


