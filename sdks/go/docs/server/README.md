# Server

### Available Operations

* [GetInfo](#getinfo) - Show server information

## GetInfo

Show server information

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Server.GetInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigInfoResponse != nil {
        // handle response
    }
}
```
