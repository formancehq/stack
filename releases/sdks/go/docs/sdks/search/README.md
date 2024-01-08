# Search
(*Search*)

### Available Operations

* [Search](#search) - Search
* [SearchgetServerInfo](#searchgetserverinfo) - Get server info

## Search

ElasticSearch query engine

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.Search.Search(ctx, shared.Query{
        After: []string{
            "users:002",
        },
        Cursor: formancesdkgo.String("YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol="),
        Ledgers: []string{
            "quickstart",
        },
        Policy: formancesdkgo.String("OR"),
        Raw: &shared.QueryRaw{},
        Sort: formancesdkgo.String("id:asc"),
        Terms: []string{
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

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |
| `request`                                             | [shared.Query](../../pkg/models/shared/query.md)      | :heavy_check_mark:                                    | The request object to use for the request.            |


### Response

**[*operations.SearchResponse](../../pkg/models/operations/searchresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 400-600            | */*                |

## SearchgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.SearchgetServerInfoResponse](../../pkg/models/operations/searchgetserverinforesponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 400-600            | */*                |
