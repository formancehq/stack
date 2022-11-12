# Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to **interface{}** | The payload | [optional] 
**Cursor** | Pointer to [**ResponseCursor**](ResponseCursor.md) |  | [optional] 
**Kind** | Pointer to **interface{}** | The kind of the object, either \&quot;TRANSACTION\&quot; or \&quot;META\&quot; | [optional] 
**Ledger** | Pointer to **interface{}** | The ledger | [optional] 

## Methods

### NewResponse

`func NewResponse() *Response`

NewResponse instantiates a new Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResponseWithDefaults

`func NewResponseWithDefaults() *Response`

NewResponseWithDefaults instantiates a new Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Response) GetData() interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Response) GetDataOk() (*interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Response) SetData(v interface{})`

SetData sets Data field to given value.

### HasData

`func (o *Response) HasData() bool`

HasData returns a boolean if a field has been set.

### SetDataNil

`func (o *Response) SetDataNil(b bool)`

 SetDataNil sets the value for Data to be an explicit nil

### UnsetData
`func (o *Response) UnsetData()`

UnsetData ensures that no value is present for Data, not even an explicit nil
### GetCursor

`func (o *Response) GetCursor() ResponseCursor`

GetCursor returns the Cursor field if non-nil, zero value otherwise.

### GetCursorOk

`func (o *Response) GetCursorOk() (*ResponseCursor, bool)`

GetCursorOk returns a tuple with the Cursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCursor

`func (o *Response) SetCursor(v ResponseCursor)`

SetCursor sets Cursor field to given value.

### HasCursor

`func (o *Response) HasCursor() bool`

HasCursor returns a boolean if a field has been set.

### GetKind

`func (o *Response) GetKind() interface{}`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *Response) GetKindOk() (*interface{}, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *Response) SetKind(v interface{})`

SetKind sets Kind field to given value.

### HasKind

`func (o *Response) HasKind() bool`

HasKind returns a boolean if a field has been set.

### SetKindNil

`func (o *Response) SetKindNil(b bool)`

 SetKindNil sets the value for Kind to be an explicit nil

### UnsetKind
`func (o *Response) UnsetKind()`

UnsetKind ensures that no value is present for Kind, not even an explicit nil
### GetLedger

`func (o *Response) GetLedger() interface{}`

GetLedger returns the Ledger field if non-nil, zero value otherwise.

### GetLedgerOk

`func (o *Response) GetLedgerOk() (*interface{}, bool)`

GetLedgerOk returns a tuple with the Ledger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedger

`func (o *Response) SetLedger(v interface{})`

SetLedger sets Ledger field to given value.

### HasLedger

`func (o *Response) HasLedger() bool`

HasLedger returns a boolean if a field has been set.

### SetLedgerNil

`func (o *Response) SetLedgerNil(b bool)`

 SetLedgerNil sets the value for Ledger to be an explicit nil

### UnsetLedger
`func (o *Response) UnsetLedger()`

UnsetLedger ensures that no value is present for Ledger, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


