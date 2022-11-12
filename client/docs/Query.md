# Query

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ledgers** | Pointer to **[]string** |  | [optional] 
**NextToken** | Pointer to **string** |  | [optional] 
**Size** | Pointer to **int32** |  | [optional] 
**Terms** | Pointer to **[]string** |  | [optional] 

## Methods

### NewQuery

`func NewQuery() *Query`

NewQuery instantiates a new Query object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueryWithDefaults

`func NewQueryWithDefaults() *Query`

NewQueryWithDefaults instantiates a new Query object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLedgers

`func (o *Query) GetLedgers() []string`

GetLedgers returns the Ledgers field if non-nil, zero value otherwise.

### GetLedgersOk

`func (o *Query) GetLedgersOk() (*[]string, bool)`

GetLedgersOk returns a tuple with the Ledgers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedgers

`func (o *Query) SetLedgers(v []string)`

SetLedgers sets Ledgers field to given value.

### HasLedgers

`func (o *Query) HasLedgers() bool`

HasLedgers returns a boolean if a field has been set.

### GetNextToken

`func (o *Query) GetNextToken() string`

GetNextToken returns the NextToken field if non-nil, zero value otherwise.

### GetNextTokenOk

`func (o *Query) GetNextTokenOk() (*string, bool)`

GetNextTokenOk returns a tuple with the NextToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextToken

`func (o *Query) SetNextToken(v string)`

SetNextToken sets NextToken field to given value.

### HasNextToken

`func (o *Query) HasNextToken() bool`

HasNextToken returns a boolean if a field has been set.

### GetSize

`func (o *Query) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Query) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Query) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *Query) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetTerms

`func (o *Query) GetTerms() []string`

GetTerms returns the Terms field if non-nil, zero value otherwise.

### GetTermsOk

`func (o *Query) GetTermsOk() (*[]string, bool)`

GetTermsOk returns a tuple with the Terms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerms

`func (o *Query) SetTerms(v []string)`

SetTerms sets Terms field to given value.

### HasTerms

`func (o *Query) HasTerms() bool`

HasTerms returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


