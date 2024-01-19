<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"log"
)

func main() {
	s := v2.New(
		v2.WithSecurity(shared.Security{
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