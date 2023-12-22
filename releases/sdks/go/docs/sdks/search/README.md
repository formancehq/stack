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
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Search.Search(ctx, shared.Query{
        After: []string{
            "u",
            "s",
            "e",
            "r",
            "s",
            ":",
            "0",
            "0",
            "2",
        },
        Cursor: formancesdkgo.String("YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol="),
        Ledgers: []string{
            "q",
            "u",
            "i",
            "c",
            "k",
            "s",
            "t",
            "a",
            "r",
            "t",
        },
        Policy: formancesdkgo.String("OR"),
        Raw: &shared.QueryRaw{},
        Sort: formancesdkgo.String("id:asc"),
        Terms: []string{
            "d",
            "e",
            "s",
            "t",
            "i",
            "n",
            "a",
            "t",
            "i",
            "o",
            "n",
            "=",
            "c",
            "e",
            "n",
            "t",
            "r",
            "a",
            "l",
            "_",
            "b",
            "a",
            "n",
            "k",
            "1",
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
| `request`                                             | [shared.Query](../../models/shared/query.md)          | :heavy_check_mark:                                    | The request object to use for the request.            |


### Response

**[*operations.SearchResponse](../../models/operations/searchresponse.md), error**


## SearchgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formancesdkgo.New()

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

**[*operations.SearchgetServerInfoResponse](../../models/operations/searchgetserverinforesponse.md), error**

