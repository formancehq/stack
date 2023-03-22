# Formance SDK

## Overview

Open, modular foundation for unique payments flows

# Introduction
This API is documented in **OpenAPI format**.

# Authentication
Formance Stack offers one forms of authentication:
  - OAuth2
OAuth2 - an open protocol to allow secure authorization in a simple
and standard method from web, mobile and desktop applications.
<SecurityDefinitions />


### Available Operations

* [GetVersions](#getversions) - Show stack version information

## GetVersions

Show stack version information

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

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
