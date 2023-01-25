# Send

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | [**Monetary1**](Monetary1.md) |  | 
**Destination** | [**Destination**](Destination.md) |  | 
**Source** | [**Source**](Source.md) |  | 

## Methods

### NewSend

`func NewSend(amount Monetary1, destination Destination, source Source, ) *Send`

NewSend instantiates a new Send object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSendWithDefaults

`func NewSendWithDefaults() *Send`

NewSendWithDefaults instantiates a new Send object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *Send) GetAmount() Monetary1`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *Send) GetAmountOk() (*Monetary1, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *Send) SetAmount(v Monetary1)`

SetAmount sets Amount field to given value.


### GetDestination

`func (o *Send) GetDestination() Destination`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *Send) GetDestinationOk() (*Destination, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *Send) SetDestination(v Destination)`

SetDestination sets Destination field to given value.


### GetSource

`func (o *Send) GetSource() Source`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *Send) GetSourceOk() (*Source, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *Send) SetSource(v Source)`

SetSource sets Source field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


