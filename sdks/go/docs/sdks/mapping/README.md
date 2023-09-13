# Mapping

### Available Operations

* [GetMapping](#getmapping) - Get the mapping of a ledger
* [UpdateMapping](#updatemapping) - Update the mapping of a ledger

## GetMapping

Get the mapping of a ledger

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
    res, err := s.Mapping.GetMapping(ctx, operations.GetMappingRequest{
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.MappingResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.GetMappingRequest](../../models/operations/getmappingrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.GetMappingResponse](../../models/operations/getmappingresponse.md), error**


## UpdateMapping

Update the mapping of a ledger

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
    res, err := s.Mapping.UpdateMapping(ctx, operations.UpdateMappingRequest{
        Mapping: shared.Mapping{
            Contracts: []shared.Contract{
                shared.Contract{
                    Account: formance.String("users:001"),
                    Expr: shared.ContractExpr{},
                },
                shared.Contract{
                    Account: formance.String("users:001"),
                    Expr: shared.ContractExpr{},
                },
            },
        },
        Ledger: "ledger001",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.MappingResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                          | Type                                                                               | Required                                                                           | Description                                                                        |
| ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ctx`                                                                              | [context.Context](https://pkg.go.dev/context#Context)                              | :heavy_check_mark:                                                                 | The context to use for the request.                                                |
| `request`                                                                          | [operations.UpdateMappingRequest](../../models/operations/updatemappingrequest.md) | :heavy_check_mark:                                                                 | The request object to use for the request.                                         |


### Response

**[*operations.UpdateMappingResponse](../../models/operations/updatemappingresponse.md), error**

