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
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
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

## UpdateMapping

Update the mapping of a ledger

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Mapping.UpdateMapping(ctx, operations.UpdateMappingRequest{
        Mapping: shared.Mapping{
            Contracts: []shared.Contract{
                shared.Contract{
                    Account: formance.String("users:001"),
                    Expr: map[string]interface{}{
                        "numquam": "commodi",
                        "quam": "molestiae",
                        "velit": "error",
                    },
                },
                shared.Contract{
                    Account: formance.String("users:001"),
                    Expr: map[string]interface{}{
                        "quis": "vitae",
                    },
                },
                shared.Contract{
                    Account: formance.String("users:001"),
                    Expr: map[string]interface{}{
                        "animi": "enim",
                        "odit": "quo",
                        "sequi": "tenetur",
                    },
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
