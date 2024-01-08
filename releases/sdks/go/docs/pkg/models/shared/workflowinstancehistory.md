# WorkflowInstanceHistory


## Fields

| Field                                               | Type                                                | Required                                            | Description                                         |
| --------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- |
| `Error`                                             | **string*                                           | :heavy_minus_sign:                                  | N/A                                                 |
| `Input`                                             | [shared.Stage](../../../pkg/models/shared/stage.md) | :heavy_check_mark:                                  | N/A                                                 |
| `Name`                                              | *string*                                            | :heavy_check_mark:                                  | N/A                                                 |
| `StartedAt`                                         | [time.Time](https://pkg.go.dev/time#Time)           | :heavy_check_mark:                                  | N/A                                                 |
| `Terminated`                                        | *bool*                                              | :heavy_check_mark:                                  | N/A                                                 |
| `TerminatedAt`                                      | [*time.Time](https://pkg.go.dev/time#Time)          | :heavy_minus_sign:                                  | N/A                                                 |