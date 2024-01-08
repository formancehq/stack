# Formance SDK


## Overview

Formance Stack API: Open, modular foundation for unique payments flows

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
* [GetAPIAuthWellKnownOpenidConfiguration](#getapiauthwellknownopenidconfiguration)

## GetVersions

Show stack version information

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
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

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.GetVersionsResponse](../../pkg/models/operations/getversionsresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 400-600            | */*                |

## GetAPIAuthWellKnownOpenidConfiguration

### Example Usage

```go
package main

import(
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"context"
	"log"
	"net/http"
)

func main() {
    s := formancesdkgo.New(
        formancesdkgo.WithSecurity(shared.Security{
            Authorization: "Bearer <YOUR_ACCESS_TOKEN_HERE>",
        }),
    )

    ctx := context.Background()
    res, err := s.GetAPIAuthWellKnownOpenidConfiguration(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

### Parameters

| Parameter                                             | Type                                                  | Required                                              | Description                                           |
| ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- | ----------------------------------------------------- |
| `ctx`                                                 | [context.Context](https://pkg.go.dev/context#Context) | :heavy_check_mark:                                    | The context to use for the request.                   |


### Response

**[*operations.GetAPIAuthWellKnownOpenidConfigurationResponse](../../pkg/models/operations/getapiauthwellknownopenidconfigurationresponse.md), error**
| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 400-600            | */*                |
