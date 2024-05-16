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
			Authorization: "<YOUR_AUTHORIZATION_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.GetOIDCWellKnowns(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage [usage] -->