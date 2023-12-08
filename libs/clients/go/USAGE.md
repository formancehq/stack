<!-- Start SDK Example Usage -->


```go
package main

import (
	"context"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"log"
)

func main() {
	s := formancesdkgo.New()

	ctx := context.Background()
	res, err := s.Formance.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage -->