# Search

### Available Operations

* [Search](#search) - Search
* [SearchgetServerInfo](#searchgetserverinfo) - Get server info

## Search

ElasticSearch query engine

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Search.Search(ctx, shared.Query{
        After: []string{
            "users:002",
            "users:002",
            "users:002",
        },
        Cursor: formance.String("YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol="),
        Ledgers: []string{
            "quickstart",
            "quickstart",
            "quickstart",
            "quickstart",
        },
        PageSize: formance.Int64(249796),
        Policy: formance.String("OR"),
        Raw: map[string]interface{}{
            "enim": "accusamus",
            "delectus": "quidem",
            "provident": "nam",
        },
        Sort: formance.String("txid:asc"),
        Target: formance.String("id"),
        Terms: []string{
            "destination=central_bank1",
            "destination=central_bank1",
            "destination=central_bank1",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.Response != nil {
        // handle response
    }
}
```

## SearchgetServerInfo

Get server info

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
    res, err := s.Search.SearchgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ServerInfo != nil {
        // handle response
    }
}
```
