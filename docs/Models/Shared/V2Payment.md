# V2Payment


## Fields

| Field                                                                                                  | Type                                                                                                   | Required                                                                                               | Description                                                                                            | Example                                                                                                |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `adjustments`                                                                                          | array<[\formance\stack\Models\Shared\V2PaymentAdjustment](../../Models/Shared/V2PaymentAdjustment.md)> | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `asset`                                                                                                | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    | USD                                                                                                    |
| `connectorID`                                                                                          | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `createdAt`                                                                                            | [\DateTime](https://www.php.net/manual/en/class.datetime.php)                                          | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `destinationAccountID`                                                                                 | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `id`                                                                                                   | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    | XXX                                                                                                    |
| `initialAmount`                                                                                        | *int*                                                                                                  | :heavy_check_mark:                                                                                     | N/A                                                                                                    | 100                                                                                                    |
| `metadata`                                                                                             | [\formance\stack\Models\Shared\V2PaymentMetadata](../../Models/Shared/V2PaymentMetadata.md)            | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `provider`                                                                                             | [?\formance\stack\Models\Shared\V2Connector](../../Models/Shared/V2Connector.md)                       | :heavy_minus_sign:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `raw`                                                                                                  | [\formance\stack\Models\Shared\V2PaymentRaw](../../Models/Shared/V2PaymentRaw.md)                      | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `reference`                                                                                            | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `scheme`                                                                                               | [\formance\stack\Models\Shared\Scheme](../../Models/Shared/Scheme.md)                                  | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `sourceAccountID`                                                                                      | *string*                                                                                               | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `status`                                                                                               | [\formance\stack\Models\Shared\V2PaymentStatus](../../Models/Shared/V2PaymentStatus.md)                | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |
| `type`                                                                                                 | [\formance\stack\Models\Shared\V2PaymentType](../../Models/Shared/V2PaymentType.md)                    | :heavy_check_mark:                                                                                     | N/A                                                                                                    |                                                                                                        |