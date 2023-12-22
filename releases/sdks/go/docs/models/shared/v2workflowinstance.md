# V2WorkflowInstance


## Fields

| Field                                                   | Type                                                    | Required                                                | Description                                             |
| ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| `CreatedAt`                                             | [time.Time](https://pkg.go.dev/time#Time)               | :heavy_check_mark:                                      | N/A                                                     |
| `Error`                                                 | **string*                                               | :heavy_minus_sign:                                      | N/A                                                     |
| `ID`                                                    | *string*                                                | :heavy_check_mark:                                      | N/A                                                     |
| `Status`                                                | [][V2StageStatus](../../models/shared/v2stagestatus.md) | :heavy_minus_sign:                                      | N/A                                                     |
| `Terminated`                                            | *bool*                                                  | :heavy_check_mark:                                      | N/A                                                     |
| `TerminatedAt`                                          | [*time.Time](https://pkg.go.dev/time#Time)              | :heavy_minus_sign:                                      | N/A                                                     |
| `UpdatedAt`                                             | [time.Time](https://pkg.go.dev/time#Time)               | :heavy_check_mark:                                      | N/A                                                     |
| `WorkflowID`                                            | *string*                                                | :heavy_check_mark:                                      | N/A                                                     |