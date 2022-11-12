# Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to **map[string]interface{}** | The payload | [optional] 
**Kind** | Pointer to **string** | The kind of the object, either \&quot;TRANSACTION\&quot; or \&quot;META\&quot; | [optional] 
**Ledger** | Pointer to **string** | The ledger | [optional] 

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

`func (o *Response) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Response) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Response) SetData(v map[string]interface{})`

SetData sets Data field to given value.

### HasData

`func (o *Response) HasData() bool`

HasData returns a boolean if a field has been set.

### GetKind

`func (o *Response) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *Response) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *Response) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *Response) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetLedger

`func (o *Response) GetLedger() string`

GetLedger returns the Ledger field if non-nil, zero value otherwise.

### GetLedgerOk

`func (o *Response) GetLedgerOk() (*string, bool)`

GetLedgerOk returns a tuple with the Ledger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedger

`func (o *Response) SetLedger(v string)`

SetLedger sets Ledger field to given value.

### HasLedger

`func (o *Response) HasLedger() bool`

HasLedger returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


