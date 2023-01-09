# LedgerInfoResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**LedgerInfo**](LedgerInfo.md) |  | [optional] 

## Methods

### NewLedgerInfoResponse

`func NewLedgerInfoResponse() *LedgerInfoResponse`

NewLedgerInfoResponse instantiates a new LedgerInfoResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLedgerInfoResponseWithDefaults

`func NewLedgerInfoResponseWithDefaults() *LedgerInfoResponse`

NewLedgerInfoResponseWithDefaults instantiates a new LedgerInfoResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *LedgerInfoResponse) GetData() LedgerInfo`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *LedgerInfoResponse) GetDataOk() (*LedgerInfo, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *LedgerInfoResponse) SetData(v LedgerInfo)`

SetData sets Data field to given value.

### HasData

`func (o *LedgerInfoResponse) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


