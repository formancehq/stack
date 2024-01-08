<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"log"
)

func main() {
	s := formancesdkgo.New(
		formancesdkgo.WithSecurity(shared.Security{
			Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
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
<!-- End SDK Example Usage [usage] -->