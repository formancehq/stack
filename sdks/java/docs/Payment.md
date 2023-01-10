

# Payment


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**provider** | **String** |  |  |
|**reference** | **String** |  |  [optional] |
|**scheme** | [**SchemeEnum**](#SchemeEnum) |  |  |
|**status** | **String** |  |  |
|**type** | [**TypeEnum**](#TypeEnum) |  |  |
|**id** | **String** |  |  |
|**amount** | **Integer** |  |  |
|**asset** | **String** |  |  |
|**date** | **OffsetDateTime** |  |  |
|**raw** | **Object** |  |  [optional] |



## Enum: SchemeEnum

| Name | Value |
|---- | -----|
| VISA | &quot;visa&quot; |
| MASTERCARD | &quot;mastercard&quot; |
| APPLE_PAY | &quot;apple pay&quot; |
| GOOGLE_PAY | &quot;google pay&quot; |
| SEPA_DEBIT | &quot;sepa debit&quot; |
| SEPA_CREDIT | &quot;sepa credit&quot; |
| SEPA | &quot;sepa&quot; |
| A2A | &quot;a2a&quot; |
| ACH_DEBIT | &quot;ach debit&quot; |
| ACH | &quot;ach&quot; |
| RTP | &quot;rtp&quot; |
| OTHER | &quot;other&quot; |



## Enum: TypeEnum

| Name | Value |
|---- | -----|
| PAY_IN | &quot;pay-in&quot; |
| PAYOUT | &quot;payout&quot; |
| OTHER | &quot;other&quot; |



