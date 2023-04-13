# ExpandedTransactionAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PreCommitVolumes** | [**map[string]map[string]Volume**](map.md) |  | 
**PostCommitVolumes** | [**map[string]map[string]Volume**](map.md) |  | 

## Methods

### NewExpandedTransactionAllOf

`func NewExpandedTransactionAllOf(preCommitVolumes map[string]map[string]Volume, postCommitVolumes map[string]map[string]Volume, ) *ExpandedTransactionAllOf`

NewExpandedTransactionAllOf instantiates a new ExpandedTransactionAllOf object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExpandedTransactionAllOfWithDefaults

`func NewExpandedTransactionAllOfWithDefaults() *ExpandedTransactionAllOf`

NewExpandedTransactionAllOfWithDefaults instantiates a new ExpandedTransactionAllOf object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPreCommitVolumes

`func (o *ExpandedTransactionAllOf) GetPreCommitVolumes() map[string]map[string]Volume`

GetPreCommitVolumes returns the PreCommitVolumes field if non-nil, zero value otherwise.

### GetPreCommitVolumesOk

`func (o *ExpandedTransactionAllOf) GetPreCommitVolumesOk() (*map[string]map[string]Volume, bool)`

GetPreCommitVolumesOk returns a tuple with the PreCommitVolumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreCommitVolumes

`func (o *ExpandedTransactionAllOf) SetPreCommitVolumes(v map[string]map[string]Volume)`

SetPreCommitVolumes sets PreCommitVolumes field to given value.


### GetPostCommitVolumes

`func (o *ExpandedTransactionAllOf) GetPostCommitVolumes() map[string]map[string]Volume`

GetPostCommitVolumes returns the PostCommitVolumes field if non-nil, zero value otherwise.

### GetPostCommitVolumesOk

`func (o *ExpandedTransactionAllOf) GetPostCommitVolumesOk() (*map[string]map[string]Volume, bool)`

GetPostCommitVolumesOk returns a tuple with the PostCommitVolumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostCommitVolumes

`func (o *ExpandedTransactionAllOf) SetPostCommitVolumes(v map[string]map[string]Volume)`

SetPostCommitVolumes sets PostCommitVolumes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


