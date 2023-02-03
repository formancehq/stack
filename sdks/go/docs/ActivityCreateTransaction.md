# ActivityCreateTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ledger** | Pointer to **string** |  | [optional] 
**Data** | Pointer to [**PostTransaction**](PostTransaction.md) |  | [optional] 

## Methods

### NewActivityCreateTransaction

`func NewActivityCreateTransaction() *ActivityCreateTransaction`

NewActivityCreateTransaction instantiates a new ActivityCreateTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActivityCreateTransactionWithDefaults

`func NewActivityCreateTransactionWithDefaults() *ActivityCreateTransaction`

NewActivityCreateTransactionWithDefaults instantiates a new ActivityCreateTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLedger

`func (o *ActivityCreateTransaction) GetLedger() string`

GetLedger returns the Ledger field if non-nil, zero value otherwise.

### GetLedgerOk

`func (o *ActivityCreateTransaction) GetLedgerOk() (*string, bool)`

GetLedgerOk returns a tuple with the Ledger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedger

`func (o *ActivityCreateTransaction) SetLedger(v string)`

SetLedger sets Ledger field to given value.

### HasLedger

`func (o *ActivityCreateTransaction) HasLedger() bool`

HasLedger returns a boolean if a field has been set.

### GetData

`func (o *ActivityCreateTransaction) GetData() PostTransaction`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ActivityCreateTransaction) GetDataOk() (*PostTransaction, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ActivityCreateTransaction) SetData(v PostTransaction)`

SetData sets Data field to given value.

### HasData

`func (o *ActivityCreateTransaction) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


