# GetWalletSummaryResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Balances** | [**[]BalanceWithAssets**](BalanceWithAssets.md) |  | 
**AvailableFunds** | **map[string]int32** |  | 
**ExpiredFunds** | **map[string]int32** |  | 
**ExpirableFunds** | **map[string]int32** |  | 
**HoldFunds** | **map[string]int32** |  | 

## Methods

### NewGetWalletSummaryResponse

`func NewGetWalletSummaryResponse(balances []BalanceWithAssets, availableFunds map[string]int32, expiredFunds map[string]int32, expirableFunds map[string]int32, holdFunds map[string]int32, ) *GetWalletSummaryResponse`

NewGetWalletSummaryResponse instantiates a new GetWalletSummaryResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetWalletSummaryResponseWithDefaults

`func NewGetWalletSummaryResponseWithDefaults() *GetWalletSummaryResponse`

NewGetWalletSummaryResponseWithDefaults instantiates a new GetWalletSummaryResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBalances

`func (o *GetWalletSummaryResponse) GetBalances() []BalanceWithAssets`

GetBalances returns the Balances field if non-nil, zero value otherwise.

### GetBalancesOk

`func (o *GetWalletSummaryResponse) GetBalancesOk() (*[]BalanceWithAssets, bool)`

GetBalancesOk returns a tuple with the Balances field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalances

`func (o *GetWalletSummaryResponse) SetBalances(v []BalanceWithAssets)`

SetBalances sets Balances field to given value.


### GetAvailableFunds

`func (o *GetWalletSummaryResponse) GetAvailableFunds() map[string]int32`

GetAvailableFunds returns the AvailableFunds field if non-nil, zero value otherwise.

### GetAvailableFundsOk

`func (o *GetWalletSummaryResponse) GetAvailableFundsOk() (*map[string]int32, bool)`

GetAvailableFundsOk returns a tuple with the AvailableFunds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableFunds

`func (o *GetWalletSummaryResponse) SetAvailableFunds(v map[string]int32)`

SetAvailableFunds sets AvailableFunds field to given value.


### GetExpiredFunds

`func (o *GetWalletSummaryResponse) GetExpiredFunds() map[string]int32`

GetExpiredFunds returns the ExpiredFunds field if non-nil, zero value otherwise.

### GetExpiredFundsOk

`func (o *GetWalletSummaryResponse) GetExpiredFundsOk() (*map[string]int32, bool)`

GetExpiredFundsOk returns a tuple with the ExpiredFunds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiredFunds

`func (o *GetWalletSummaryResponse) SetExpiredFunds(v map[string]int32)`

SetExpiredFunds sets ExpiredFunds field to given value.


### GetExpirableFunds

`func (o *GetWalletSummaryResponse) GetExpirableFunds() map[string]int32`

GetExpirableFunds returns the ExpirableFunds field if non-nil, zero value otherwise.

### GetExpirableFundsOk

`func (o *GetWalletSummaryResponse) GetExpirableFundsOk() (*map[string]int32, bool)`

GetExpirableFundsOk returns a tuple with the ExpirableFunds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpirableFunds

`func (o *GetWalletSummaryResponse) SetExpirableFunds(v map[string]int32)`

SetExpirableFunds sets ExpirableFunds field to given value.


### GetHoldFunds

`func (o *GetWalletSummaryResponse) GetHoldFunds() map[string]int32`

GetHoldFunds returns the HoldFunds field if non-nil, zero value otherwise.

### GetHoldFundsOk

`func (o *GetWalletSummaryResponse) GetHoldFundsOk() (*map[string]int32, bool)`

GetHoldFundsOk returns a tuple with the HoldFunds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHoldFunds

`func (o *GetWalletSummaryResponse) SetHoldFunds(v map[string]int32)`

SetHoldFunds sets HoldFunds field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


