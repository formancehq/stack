# ExpandedTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | **time.Time** |  | 
**Postings** | [**[]Posting**](Posting.md) |  | 
**Reference** | Pointer to **string** |  | [optional] 
**Metadata** | **map[string]string** |  | 
**Txid** | **int64** |  | 
**PreCommitVolumes** | [**map[string]map[string]Volume**](map.md) |  | 
**PostCommitVolumes** | [**map[string]map[string]Volume**](map.md) |  | 

## Methods

### NewExpandedTransaction

`func NewExpandedTransaction(timestamp time.Time, postings []Posting, metadata map[string]string, txid int64, preCommitVolumes map[string]map[string]Volume, postCommitVolumes map[string]map[string]Volume, ) *ExpandedTransaction`

NewExpandedTransaction instantiates a new ExpandedTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExpandedTransactionWithDefaults

`func NewExpandedTransactionWithDefaults() *ExpandedTransaction`

NewExpandedTransactionWithDefaults instantiates a new ExpandedTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *ExpandedTransaction) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *ExpandedTransaction) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *ExpandedTransaction) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.


### GetPostings

`func (o *ExpandedTransaction) GetPostings() []Posting`

GetPostings returns the Postings field if non-nil, zero value otherwise.

### GetPostingsOk

`func (o *ExpandedTransaction) GetPostingsOk() (*[]Posting, bool)`

GetPostingsOk returns a tuple with the Postings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostings

`func (o *ExpandedTransaction) SetPostings(v []Posting)`

SetPostings sets Postings field to given value.


### GetReference

`func (o *ExpandedTransaction) GetReference() string`

GetReference returns the Reference field if non-nil, zero value otherwise.

### GetReferenceOk

`func (o *ExpandedTransaction) GetReferenceOk() (*string, bool)`

GetReferenceOk returns a tuple with the Reference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReference

`func (o *ExpandedTransaction) SetReference(v string)`

SetReference sets Reference field to given value.

### HasReference

`func (o *ExpandedTransaction) HasReference() bool`

HasReference returns a boolean if a field has been set.

### GetMetadata

`func (o *ExpandedTransaction) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ExpandedTransaction) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ExpandedTransaction) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.


### GetTxid

`func (o *ExpandedTransaction) GetTxid() int64`

GetTxid returns the Txid field if non-nil, zero value otherwise.

### GetTxidOk

`func (o *ExpandedTransaction) GetTxidOk() (*int64, bool)`

GetTxidOk returns a tuple with the Txid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxid

`func (o *ExpandedTransaction) SetTxid(v int64)`

SetTxid sets Txid field to given value.


### GetPreCommitVolumes

`func (o *ExpandedTransaction) GetPreCommitVolumes() map[string]map[string]Volume`

GetPreCommitVolumes returns the PreCommitVolumes field if non-nil, zero value otherwise.

### GetPreCommitVolumesOk

`func (o *ExpandedTransaction) GetPreCommitVolumesOk() (*map[string]map[string]Volume, bool)`

GetPreCommitVolumesOk returns a tuple with the PreCommitVolumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreCommitVolumes

`func (o *ExpandedTransaction) SetPreCommitVolumes(v map[string]map[string]Volume)`

SetPreCommitVolumes sets PreCommitVolumes field to given value.


### GetPostCommitVolumes

`func (o *ExpandedTransaction) GetPostCommitVolumes() map[string]map[string]Volume`

GetPostCommitVolumes returns the PostCommitVolumes field if non-nil, zero value otherwise.

### GetPostCommitVolumesOk

`func (o *ExpandedTransaction) GetPostCommitVolumesOk() (*map[string]map[string]Volume, bool)`

GetPostCommitVolumesOk returns a tuple with the PostCommitVolumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostCommitVolumes

`func (o *ExpandedTransaction) SetPostCommitVolumes(v map[string]map[string]Volume)`

SetPostCommitVolumes sets PostCommitVolumes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


