<!-- Start SDK Example Usage -->


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
            Authorization: "",
        }),
    )

    ctx := context.Background()
    res, err := s.GetVersions(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.GetVersionsResponse != nil {
        // handle response
    }
}
```
<!-- End SDK Example Usage -->