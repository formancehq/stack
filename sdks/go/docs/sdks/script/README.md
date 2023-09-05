# Script

### Available Operations

* [~~RunScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**

## ~~RunScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "",
        }),
    )

    ctx := context.Background()
    res, err := s.Script.RunScript(ctx, operations.RunScriptRequest{
        Script: shared.Script{
            Metadata: map[string]interface{}{
                "rem": "voluptates",
                "quasi": "repudiandae",
                "sint": "veritatis",
            },
            Plain: "vars {
        account $user
        }
        send [COIN 10] (
        	source = @world
        	destination = $user
        )
        ",
            Reference: formance.String("order_1234"),
            Vars: &shared.ScriptVars{},
        },
        Ledger: "ledger001",
        Preview: formance.Bool(true),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ScriptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `ctx`                                                                      | [context.Context](https://pkg.go.dev/context#Context)                      | :heavy_check_mark:                                                         | The context to use for the request.                                        |
| `request`                                                                  | [operations.RunScriptRequest](../../models/operations/runscriptrequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |


### Response

**[*operations.RunScriptResponse](../../models/operations/runscriptresponse.md), error**

