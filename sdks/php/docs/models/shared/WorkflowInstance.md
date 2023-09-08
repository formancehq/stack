# WorkflowInstance


## Fields

| Field                                                         | Type                                                          | Required                                                      | Description                                                   |
| ------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------------------------------------- |
| `createdAt`                                                   | [\DateTime](https://www.php.net/manual/en/class.datetime.php) | :heavy_check_mark:                                            | N/A                                                           |
| `error`                                                       | *?string*                                                     | :heavy_minus_sign:                                            | N/A                                                           |
| `id`                                                          | *string*                                                      | :heavy_check_mark:                                            | N/A                                                           |
| `status`                                                      | array<[StageStatus](../../models/shared/StageStatus.md)>      | :heavy_minus_sign:                                            | N/A                                                           |
| `terminated`                                                  | *bool*                                                        | :heavy_check_mark:                                            | N/A                                                           |
| `terminatedAt`                                                | [\DateTime](https://www.php.net/manual/en/class.datetime.php) | :heavy_minus_sign:                                            | N/A                                                           |
| `updatedAt`                                                   | [\DateTime](https://www.php.net/manual/en/class.datetime.php) | :heavy_check_mark:                                            | N/A                                                           |
| `workflowID`                                                  | *string*                                                      | :heavy_check_mark:                                            | N/A                                                           |