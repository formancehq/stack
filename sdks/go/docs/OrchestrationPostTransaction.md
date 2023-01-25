# OrchestrationPostTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**Postings** | Pointer to [**[]Posting**](Posting.md) |  | [optional] 
**Script** | Pointer to [**PostTransactionScript**](PostTransactionScript.md) |  | [optional] 
**Reference** | Pointer to **string** |  | [optional] 
**Metadata** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewOrchestrationPostTransaction

`func NewOrchestrationPostTransaction() *OrchestrationPostTransaction`

NewOrchestrationPostTransaction instantiates a new OrchestrationPostTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrchestrationPostTransactionWithDefaults

`func NewOrchestrationPostTransactionWithDefaults() *OrchestrationPostTransaction`

NewOrchestrationPostTransactionWithDefaults instantiates a new OrchestrationPostTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *OrchestrationPostTransaction) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *OrchestrationPostTransaction) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *OrchestrationPostTransaction) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *OrchestrationPostTransaction) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetPostings

`func (o *OrchestrationPostTransaction) GetPostings() []Posting`

GetPostings returns the Postings field if non-nil, zero value otherwise.

### GetPostingsOk

`func (o *OrchestrationPostTransaction) GetPostingsOk() (*[]Posting, bool)`

GetPostingsOk returns a tuple with the Postings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostings

`func (o *OrchestrationPostTransaction) SetPostings(v []Posting)`

SetPostings sets Postings field to given value.

### HasPostings

`func (o *OrchestrationPostTransaction) HasPostings() bool`

HasPostings returns a boolean if a field has been set.

### GetScript

`func (o *OrchestrationPostTransaction) GetScript() PostTransactionScript`

GetScript returns the Script field if non-nil, zero value otherwise.

### GetScriptOk

`func (o *OrchestrationPostTransaction) GetScriptOk() (*PostTransactionScript, bool)`

GetScriptOk returns a tuple with the Script field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScript

`func (o *OrchestrationPostTransaction) SetScript(v PostTransactionScript)`

SetScript sets Script field to given value.

### HasScript

`func (o *OrchestrationPostTransaction) HasScript() bool`

HasScript returns a boolean if a field has been set.

### GetReference

`func (o *OrchestrationPostTransaction) GetReference() string`

GetReference returns the Reference field if non-nil, zero value otherwise.

### GetReferenceOk

`func (o *OrchestrationPostTransaction) GetReferenceOk() (*string, bool)`

GetReferenceOk returns a tuple with the Reference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReference

`func (o *OrchestrationPostTransaction) SetReference(v string)`

SetReference sets Reference field to given value.

### HasReference

`func (o *OrchestrationPostTransaction) HasReference() bool`

HasReference returns a boolean if a field has been set.

### GetMetadata

`func (o *OrchestrationPostTransaction) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *OrchestrationPostTransaction) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *OrchestrationPostTransaction) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *OrchestrationPostTransaction) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### SetMetadataNil

`func (o *OrchestrationPostTransaction) SetMetadataNil(b bool)`

 SetMetadataNil sets the value for Metadata to be an explicit nil

### UnsetMetadata
`func (o *OrchestrationPostTransaction) UnsetMetadata()`

UnsetMetadata ensures that no value is present for Metadata, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


