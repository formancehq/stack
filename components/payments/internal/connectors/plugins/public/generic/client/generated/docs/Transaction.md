# Transaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**RelatedTransactionID** | Pointer to **string** |  | [optional] 
**CreatedAt** | **time.Time** |  | 
**UpdatedAt** | **time.Time** |  | 
**Currency** | **string** |  | 
**Scheme** | Pointer to **string** |  | [optional] 
**Type** | [**TransactionType**](TransactionType.md) |  | 
**Status** | [**TransactionStatus**](TransactionStatus.md) |  | 
**Amount** | **string** |  | 
**SourceAccountID** | Pointer to **string** |  | [optional] 
**DestinationAccountID** | Pointer to **string** |  | [optional] 
**Metadata** | Pointer to **map[string]string** |  | [optional] 

## Methods

### NewTransaction

`func NewTransaction(id string, createdAt time.Time, updatedAt time.Time, currency string, type_ TransactionType, status TransactionStatus, amount string, ) *Transaction`

NewTransaction instantiates a new Transaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionWithDefaults

`func NewTransactionWithDefaults() *Transaction`

NewTransactionWithDefaults instantiates a new Transaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Transaction) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Transaction) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Transaction) SetId(v string)`

SetId sets Id field to given value.


### GetRelatedTransactionID

`func (o *Transaction) GetRelatedTransactionID() string`

GetRelatedTransactionID returns the RelatedTransactionID field if non-nil, zero value otherwise.

### GetRelatedTransactionIDOk

`func (o *Transaction) GetRelatedTransactionIDOk() (*string, bool)`

GetRelatedTransactionIDOk returns a tuple with the RelatedTransactionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelatedTransactionID

`func (o *Transaction) SetRelatedTransactionID(v string)`

SetRelatedTransactionID sets RelatedTransactionID field to given value.

### HasRelatedTransactionID

`func (o *Transaction) HasRelatedTransactionID() bool`

HasRelatedTransactionID returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Transaction) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Transaction) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Transaction) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *Transaction) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Transaction) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Transaction) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetCurrency

`func (o *Transaction) GetCurrency() string`

GetCurrency returns the Currency field if non-nil, zero value otherwise.

### GetCurrencyOk

`func (o *Transaction) GetCurrencyOk() (*string, bool)`

GetCurrencyOk returns a tuple with the Currency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrency

`func (o *Transaction) SetCurrency(v string)`

SetCurrency sets Currency field to given value.


### GetScheme

`func (o *Transaction) GetScheme() string`

GetScheme returns the Scheme field if non-nil, zero value otherwise.

### GetSchemeOk

`func (o *Transaction) GetSchemeOk() (*string, bool)`

GetSchemeOk returns a tuple with the Scheme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheme

`func (o *Transaction) SetScheme(v string)`

SetScheme sets Scheme field to given value.

### HasScheme

`func (o *Transaction) HasScheme() bool`

HasScheme returns a boolean if a field has been set.

### GetType

`func (o *Transaction) GetType() TransactionType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Transaction) GetTypeOk() (*TransactionType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Transaction) SetType(v TransactionType)`

SetType sets Type field to given value.


### GetStatus

`func (o *Transaction) GetStatus() TransactionStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Transaction) GetStatusOk() (*TransactionStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Transaction) SetStatus(v TransactionStatus)`

SetStatus sets Status field to given value.


### GetAmount

`func (o *Transaction) GetAmount() string`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *Transaction) GetAmountOk() (*string, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *Transaction) SetAmount(v string)`

SetAmount sets Amount field to given value.


### GetSourceAccountID

`func (o *Transaction) GetSourceAccountID() string`

GetSourceAccountID returns the SourceAccountID field if non-nil, zero value otherwise.

### GetSourceAccountIDOk

`func (o *Transaction) GetSourceAccountIDOk() (*string, bool)`

GetSourceAccountIDOk returns a tuple with the SourceAccountID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceAccountID

`func (o *Transaction) SetSourceAccountID(v string)`

SetSourceAccountID sets SourceAccountID field to given value.

### HasSourceAccountID

`func (o *Transaction) HasSourceAccountID() bool`

HasSourceAccountID returns a boolean if a field has been set.

### GetDestinationAccountID

`func (o *Transaction) GetDestinationAccountID() string`

GetDestinationAccountID returns the DestinationAccountID field if non-nil, zero value otherwise.

### GetDestinationAccountIDOk

`func (o *Transaction) GetDestinationAccountIDOk() (*string, bool)`

GetDestinationAccountIDOk returns a tuple with the DestinationAccountID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationAccountID

`func (o *Transaction) SetDestinationAccountID(v string)`

SetDestinationAccountID sets DestinationAccountID field to given value.

### HasDestinationAccountID

`func (o *Transaction) HasDestinationAccountID() bool`

HasDestinationAccountID returns a boolean if a field has been set.

### GetMetadata

`func (o *Transaction) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *Transaction) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *Transaction) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *Transaction) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### SetMetadataNil

`func (o *Transaction) SetMetadataNil(b bool)`

 SetMetadataNil sets the value for Metadata to be an explicit nil

### UnsetMetadata
`func (o *Transaction) UnsetMetadata()`

UnsetMetadata ensures that no value is present for Metadata, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


