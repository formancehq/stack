# PaymentResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**Payment**](Payment.md) |  | [optional] 

## Methods

### NewPaymentResponse

`func NewPaymentResponse() *PaymentResponse`

NewPaymentResponse instantiates a new PaymentResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaymentResponseWithDefaults

`func NewPaymentResponseWithDefaults() *PaymentResponse`

NewPaymentResponseWithDefaults instantiates a new PaymentResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *PaymentResponse) GetData() Payment`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PaymentResponse) GetDataOk() (*Payment, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PaymentResponse) SetData(v Payment)`

SetData sets Data field to given value.

### HasData

`func (o *PaymentResponse) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


