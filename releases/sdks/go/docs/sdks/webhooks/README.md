# Webhooks
(*Webhooks*)

### Available Operations

* [ActivateConfig](#activateconfig) - Activate one config
* [ChangeConfigSecret](#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](#deactivateconfig) - Deactivate one config
* [DeleteConfig](#deleteconfig) - Delete one config
* [GetManyConfigs](#getmanyconfigs) - Get many configs
* [InsertConfig](#insertconfig) - Insert a new config
* [TestConfig](#testconfig) - Test one config

## ActivateConfig

Activate a webhooks config by ID, to start receiving webhooks to its endpoint.

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.ActivateConfig(ctx, operations.ActivateConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.ActivateConfigRequest](../../models/operations/activateconfigrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.ActivateConfigResponse](../../models/operations/activateconfigresponse.md), error**


## ChangeConfigSecret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.ChangeConfigSecret(ctx, operations.ChangeConfigSecretRequest{
        ConfigChangeSecret: &shared.ConfigChangeSecret{
            Secret: "V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3",
        },
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.ChangeConfigSecretRequest](../../models/operations/changeconfigsecretrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |


### Response

**[*operations.ChangeConfigSecretResponse](../../models/operations/changeconfigsecretresponse.md), error**


## DeactivateConfig

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.DeactivateConfig(ctx, operations.DeactivateConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [operations.DeactivateConfigRequest](../../models/operations/deactivateconfigrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |


### Response

**[*operations.DeactivateConfigResponse](../../models/operations/deactivateconfigresponse.md), error**


## DeleteConfig

Delete a webhooks config by ID.

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.DeleteConfig(ctx, operations.DeleteConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DeleteConfigRequest](../../models/operations/deleteconfigrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |


### Response

**[*operations.DeleteConfigResponse](../../models/operations/deleteconfigresponse.md), error**


## GetManyConfigs

Sorted by updated date descending

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.GetManyConfigs(ctx, operations.GetManyConfigsRequest{
        Endpoint: formancesdkgo.String("https://example.com"),
        ID: formancesdkgo.String("4997257d-dfb6-445b-929c-cbe2ab182818"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                            | Type                                                                                 | Required                                                                             | Description                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| `ctx`                                                                                | [context.Context](https://pkg.go.dev/context#Context)                                | :heavy_check_mark:                                                                   | The context to use for the request.                                                  |
| `request`                                                                            | [operations.GetManyConfigsRequest](../../models/operations/getmanyconfigsrequest.md) | :heavy_check_mark:                                                                   | The request object to use for the request.                                           |


### Response

**[*operations.GetManyConfigsResponse](../../models/operations/getmanyconfigsresponse.md), error**


## InsertConfig

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.InsertConfig(ctx, shared.ConfigUser{
        Endpoint: "https://example.com",
        EventTypes: []string{
            "TYPE1",
            "TYPE2",
        },
        Secret: formancesdkgo.String("V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ConfigResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                              | Type                                                   | Required                                               | Description                                            |
| ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------ |
| `ctx`                                                  | [context.Context](https://pkg.go.dev/context#Context)  | :heavy_check_mark:                                     | The context to use for the request.                    |
| `request`                                              | [shared.ConfigUser](../../models/shared/configuser.md) | :heavy_check_mark:                                     | The request object to use for the request.             |


### Response

**[*operations.InsertConfigResponse](../../models/operations/insertconfigresponse.md), error**


## TestConfig

Test a config by sending a webhook to its endpoint.

### Example Usage

```go
package main

import(
	"context"
	"log"
	formancesdkgo "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formancesdkgo.New()

    ctx := context.Background()
    res, err := s.Webhooks.TestConfig(ctx, operations.TestConfigRequest{
        ID: "4997257d-dfb6-445b-929c-cbe2ab182818",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.AttemptResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.TestConfigRequest](../../models/operations/testconfigrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |


### Response

**[*operations.TestConfigResponse](../../models/operations/testconfigresponse.md), error**

