# FormanceHQ.model.attempt.Attempt

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict, str, date, datetime, uuid.UUID, int, float, decimal.Decimal, bool, None, list, tuple, bytes, io.FileIO, io.BufferedReader,  | frozendict.frozendict, str, decimal.Decimal, BoolClass, NoneClass, tuple, bytes, FileIO |  | 

### Dictionary Keys
Key | Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | ------------- | -------------
**id** | str, uuid.UUID,  | str,  |  | [optional] value must be a uuid
**webhookID** | str, uuid.UUID,  | str,  |  | [optional] value must be a uuid
**createdAt** | str, datetime,  | str,  |  | [optional] value must conform to RFC-3339 date-time
**updatedAt** | str, datetime,  | str,  |  | [optional] value must conform to RFC-3339 date-time
**config** | [**WebhooksConfig**](WebhooksConfig.md) | [**WebhooksConfig**](WebhooksConfig.md) |  | [optional] 
**payload** | str,  | str,  |  | [optional] 
**statusCode** | decimal.Decimal, int,  | decimal.Decimal,  |  | [optional] 
**retryAttempt** | decimal.Decimal, int,  | decimal.Decimal,  |  | [optional] 
**status** | str,  | str,  |  | [optional] 
**nextRetryAfter** | str, datetime,  | str,  |  | [optional] value must conform to RFC-3339 date-time
**any_string_name** | dict, frozendict.frozendict, str, date, datetime, int, float, bool, decimal.Decimal, None, list, tuple, bytes, io.FileIO, io.BufferedReader | frozendict.frozendict, str, BoolClass, decimal.Decimal, NoneClass, tuple, bytes, FileIO | any string name can be used but the value must be the correct type | [optional]

[[Back to Model list]](../../README.md#documentation-for-models) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to README]](../../README.md)

