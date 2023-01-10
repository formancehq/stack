# Formance.model.task_stripe.TaskStripe

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict,  | frozendict.frozendict,  |  | 

### Dictionary Keys
Key | Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | ------------- | -------------
**id** | str, uuid.UUID,  | str,  |  | [optional] value must be a uuid
**connectorID** | str, uuid.UUID,  | str,  |  | [optional] value must be a uuid
**createdAt** | str, datetime,  | str,  |  | [optional] value must conform to RFC-3339 date-time
**updatedAt** | str, datetime,  | str,  |  | [optional] value must conform to RFC-3339 date-time
**[descriptor](#descriptor)** | dict, frozendict.frozendict,  | frozendict.frozendict,  |  | [optional] 
**status** | [**PaymentStatus**](PaymentStatus.md) | [**PaymentStatus**](PaymentStatus.md) |  | [optional] 
**[state](#state)** | dict, frozendict.frozendict,  | frozendict.frozendict,  |  | [optional] 
**error** | str,  | str,  |  | [optional] 
**any_string_name** | dict, frozendict.frozendict, str, date, datetime, int, float, bool, decimal.Decimal, None, list, tuple, bytes, io.FileIO, io.BufferedReader | frozendict.frozendict, str, BoolClass, decimal.Decimal, NoneClass, tuple, bytes, FileIO | any string name can be used but the value must be the correct type | [optional]

# descriptor

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict,  | frozendict.frozendict,  |  | 

### Dictionary Keys
Key | Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | ------------- | -------------
**name** | str,  | str,  |  | 
**main** | bool,  | BoolClass,  |  | [optional] 
**account** | str,  | str,  |  | [optional] 
**any_string_name** | dict, frozendict.frozendict, str, date, datetime, int, float, bool, decimal.Decimal, None, list, tuple, bytes, io.FileIO, io.BufferedReader | frozendict.frozendict, str, BoolClass, decimal.Decimal, NoneClass, tuple, bytes, FileIO | any string name can be used but the value must be the correct type | [optional]

# state

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict,  | frozendict.frozendict,  |  | 

[[Back to Model list]](../../README.md#documentation-for-models) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to README]](../../README.md)

