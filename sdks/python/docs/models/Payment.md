# Formance.model.payment.Payment

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict,  | frozendict.frozendict,  |  | 

### Dictionary Keys
Key | Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | ------------- | -------------
**date** | str, datetime,  | str,  |  | value must conform to RFC-3339 date-time
**amount** | decimal.Decimal, int,  | decimal.Decimal,  |  | 
**scheme** | str,  | str,  |  | must be one of ["visa", "mastercard", "apple pay", "google pay", "sepa debit", "sepa credit", "sepa", "a2a", "ach debit", "ach", "rtp", "other", ] 
**provider** | str,  | str,  |  | 
**id** | str,  | str,  |  | 
**asset** | str,  | str,  |  | 
**type** | str,  | str,  |  | must be one of ["pay-in", "payout", "other", ] 
**status** | str,  | str,  |  | 
**reference** | str,  | str,  |  | [optional] 
**raw** | dict, frozendict.frozendict, str, date, datetime, uuid.UUID, int, float, decimal.Decimal, bool, None, list, tuple, bytes, io.FileIO, io.BufferedReader,  | frozendict.frozendict, str, decimal.Decimal, BoolClass, NoneClass, tuple, bytes, FileIO |  | [optional] 
**any_string_name** | dict, frozendict.frozendict, str, date, datetime, int, float, bool, decimal.Decimal, None, list, tuple, bytes, io.FileIO, io.BufferedReader | frozendict.frozendict, str, BoolClass, decimal.Decimal, NoneClass, tuple, bytes, FileIO | any string name can be used but the value must be the correct type | [optional]

[[Back to Model list]](../../README.md#documentation-for-models) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to README]](../../README.md)

