# Logs

### Available Operations

* [ListLogs](#listlogs) - List the logs from a ledger

## ListLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/types"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Logs.ListLogs(ctx, operations.ListLogsRequest{
        Cursor: formance.String("aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ=="),
        EndTime: types.MustTimeFromString("2022-09-04T08:35:09.957Z"),
        Ledger: "ledger001",
        PageSize: formance.Int64(570197),
        StartTime: types.MustTimeFromString("2022-07-24T21:51:02.112Z"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.LogsCursorResponse != nil {
        // handle response
    }
}
```
