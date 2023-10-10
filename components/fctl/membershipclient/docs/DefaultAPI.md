# \DefaultAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptInvitation**](DefaultAPI.md#AcceptInvitation) | **Post** /me/invitations/{invitationId}/accept | Accept invitation
[**BillingPortal**](DefaultAPI.md#BillingPortal) | **Get** /organizations/{organizationId}/billing/portal | Access to the billing portal
[**BillingSetup**](DefaultAPI.md#BillingSetup) | **Get** /organizations/{organizationId}/billing/setup | Create a billing setup
[**CreateInvitation**](DefaultAPI.md#CreateInvitation) | **Post** /organizations/{organizationId}/invitations | Create invitation
[**CreateOrganization**](DefaultAPI.md#CreateOrganization) | **Post** /organizations | Create organization
[**CreatePrivateRegion**](DefaultAPI.md#CreatePrivateRegion) | **Post** /organizations/{organizationId}/regions | Create a private region
[**CreateStack**](DefaultAPI.md#CreateStack) | **Post** /organizations/{organizationId}/stacks | Create stack
[**DeclineInvitation**](DefaultAPI.md#DeclineInvitation) | **Post** /me/invitations/{invitationId}/reject | Decline invitation
[**DeleteInvitation**](DefaultAPI.md#DeleteInvitation) | **Delete** /organizations/{organizationId}/invitations/{invitationId} | Delete invitation
[**DeleteOrganization**](DefaultAPI.md#DeleteOrganization) | **Delete** /organizations/{organizationId} | Delete organization
[**DeleteRegion**](DefaultAPI.md#DeleteRegion) | **Delete** /organizations/{organizationId}/regions/{regionId} | Delete region
[**DeleteStack**](DefaultAPI.md#DeleteStack) | **Delete** /organizations/{organizationId}/stacks/{stackId} | Delete stack
[**DisableStack**](DefaultAPI.md#DisableStack) | **Put** /organizations/{organizationId}/stacks/{stackId}/disable | Disable stack
[**DisableStargate**](DefaultAPI.md#DisableStargate) | **Put** /organizations/{organizationId}/stacks/{stackId}/stargate/disable | Disable stargate on a stack
[**EnableStargate**](DefaultAPI.md#EnableStargate) | **Put** /organizations/{organizationId}/stacks/{stackId}/stargate/enable | Enable stargate on a stack
[**GetRegion**](DefaultAPI.md#GetRegion) | **Get** /organizations/{organizationId}/regions/{regionId} | Get region
[**GetServerInfo**](DefaultAPI.md#GetServerInfo) | **Get** /_info | Get server info
[**ListInvitations**](DefaultAPI.md#ListInvitations) | **Get** /me/invitations | List invitations of the user
[**ListOrganizationInvitations**](DefaultAPI.md#ListOrganizationInvitations) | **Get** /organizations/{organizationId}/invitations | List invitations of the organization
[**ListOrganizations**](DefaultAPI.md#ListOrganizations) | **Get** /organizations | List organizations of the connected user
[**ListOrganizationsExpanded**](DefaultAPI.md#ListOrganizationsExpanded) | **Get** /organizations/expanded | List organizations of the connected user with expanded data
[**ListRegions**](DefaultAPI.md#ListRegions) | **Get** /organizations/{organizationId}/regions | List regions
[**ListStacks**](DefaultAPI.md#ListStacks) | **Get** /organizations/{organizationId}/stacks | List stacks
[**ListUsers**](DefaultAPI.md#ListUsers) | **Get** /organizations/{organizationId}/users | List users
[**ReadConnectedUser**](DefaultAPI.md#ReadConnectedUser) | **Get** /me | Read user
[**ReadOrganization**](DefaultAPI.md#ReadOrganization) | **Get** /organizations/{organizationId} | Read organization
[**ReadStack**](DefaultAPI.md#ReadStack) | **Get** /organizations/{organizationId}/stacks/{stackId} | Read stack
[**ReadUser**](DefaultAPI.md#ReadUser) | **Get** /organizations/{organizationId}/users/{userId} | Read user
[**RestoreStack**](DefaultAPI.md#RestoreStack) | **Put** /organizations/{organizationId}/stacks/{stackId}/restore | Restore stack
[**UnlinkUserFromOrganization**](DefaultAPI.md#UnlinkUserFromOrganization) | **Delete** /organizations/{organizationId}/users/{userId} | Unlink user from organization
[**UpdateConnectedUser**](DefaultAPI.md#UpdateConnectedUser) | **Put** /me | Update user
[**UpdateUser**](DefaultAPI.md#UpdateUser) | **Put** /organizations/{organizationId}/users/{userId} | Update user



## AcceptInvitation

> AcceptInvitation(ctx, invitationId).Execute()

Accept invitation

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    invitationId := "invitationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.AcceptInvitation(context.Background(), invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.AcceptInvitation``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**invitationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAcceptInvitationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BillingPortal

> BillingPortalResponse BillingPortal(ctx, organizationId).Execute()

Access to the billing portal

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.BillingPortal(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.BillingPortal``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BillingPortal`: BillingPortalResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.BillingPortal`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiBillingPortalRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BillingPortalResponse**](BillingPortalResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BillingSetup

> BillingSetupResponse BillingSetup(ctx, organizationId).Execute()

Create a billing setup

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.BillingSetup(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.BillingSetup``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BillingSetup`: BillingSetupResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.BillingSetup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiBillingSetupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BillingSetupResponse**](BillingSetupResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateInvitation

> CreateInvitationResponse CreateInvitation(ctx, organizationId).Email(email).Execute()

Create invitation

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    email := "email_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.CreateInvitation(context.Background(), organizationId).Email(email).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateInvitation``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateInvitation`: CreateInvitationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateInvitation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInvitationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **email** | **string** |  | 

### Return type

[**CreateInvitationResponse**](CreateInvitationResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateOrganization

> CreateOrganizationResponse CreateOrganization(ctx).Body(body).Execute()

Create organization

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    body := OrganizationData(987) // OrganizationData |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.CreateOrganization(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOrganization`: CreateOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateOrganization`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **OrganizationData** |  | 

### Return type

[**CreateOrganizationResponse**](CreateOrganizationResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreatePrivateRegion

> CreatedPrivateRegionResponse CreatePrivateRegion(ctx, organizationId).CreatePrivateRegionRequest(createPrivateRegionRequest).Execute()

Create a private region

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    createPrivateRegionRequest := *openapiclient.NewCreatePrivateRegionRequest("Name_example") // CreatePrivateRegionRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.CreatePrivateRegion(context.Background(), organizationId).CreatePrivateRegionRequest(createPrivateRegionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreatePrivateRegion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreatePrivateRegion`: CreatedPrivateRegionResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreatePrivateRegion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreatePrivateRegionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createPrivateRegionRequest** | [**CreatePrivateRegionRequest**](CreatePrivateRegionRequest.md) |  | 

### Return type

[**CreatedPrivateRegionResponse**](CreatedPrivateRegionResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateStack

> CreateStackResponse CreateStack(ctx, organizationId).CreateStackRequest(createStackRequest).Execute()

Create stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    createStackRequest := *openapiclient.NewCreateStackRequest("Name_example", map[string]string{"key": "Inner_example"}, "RegionID_example") // CreateStackRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.CreateStack(context.Background(), organizationId).CreateStackRequest(createStackRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateStack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createStackRequest** | [**CreateStackRequest**](CreateStackRequest.md) |  | 

### Return type

[**CreateStackResponse**](CreateStackResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeclineInvitation

> DeclineInvitation(ctx, invitationId).Execute()

Decline invitation

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    invitationId := "invitationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DeclineInvitation(context.Background(), invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeclineInvitation``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**invitationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeclineInvitationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInvitation

> DeleteInvitation(ctx, organizationId, invitationId).Execute()

Delete invitation

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    invitationId := "invitationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DeleteInvitation(context.Background(), organizationId, invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteInvitation``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**invitationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInvitationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOrganization

> DeleteOrganization(ctx, organizationId).Execute()

Delete organization

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DeleteOrganization(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteRegion

> DeleteRegion(ctx, organizationId, regionId).Execute()

Delete region

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    regionId := "regionId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DeleteRegion(context.Background(), organizationId, regionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteRegion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**regionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRegionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteStack

> DeleteStack(ctx, organizationId, stackId).Execute()

Delete stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DeleteStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DisableStack

> DisableStack(ctx, organizationId, stackId).Execute()

Disable stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DisableStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DisableStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDisableStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DisableStargate

> DisableStargate(ctx, organizationId, stackId).Execute()

Disable stargate on a stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.DisableStargate(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DisableStargate``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDisableStargateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EnableStargate

> EnableStargate(ctx, organizationId, stackId).Execute()

Enable stargate on a stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.EnableStargate(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.EnableStargate``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiEnableStargateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRegion

> GetRegionResponse GetRegion(ctx, organizationId, regionId).Execute()

Get region

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    regionId := "regionId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.GetRegion(context.Background(), organizationId, regionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetRegion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRegion`: GetRegionResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetRegion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**regionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRegionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetRegionResponse**](GetRegionResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetServerInfo

> ServerInfo GetServerInfo(ctx).Execute()

Get server info

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.GetServerInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetServerInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetServerInfo`: ServerInfo
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetServerInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetServerInfoRequest struct via the builder pattern


### Return type

[**ServerInfo**](ServerInfo.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInvitations

> ListInvitationsResponse ListInvitations(ctx).Status(status).Organization(organization).Execute()

List invitations of the user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    status := "status_example" // string | Status of organizations (optional)
    organization := "organization_example" // string | Status of organizations (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListInvitations(context.Background()).Status(status).Organization(organization).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListInvitations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListInvitations`: ListInvitationsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListInvitations`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInvitationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **status** | **string** | Status of organizations | 
 **organization** | **string** | Status of organizations | 

### Return type

[**ListInvitationsResponse**](ListInvitationsResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrganizationInvitations

> ListInvitationsResponse ListOrganizationInvitations(ctx, organizationId).Status(status).Execute()

List invitations of the organization

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    status := "status_example" // string | Status of organizations (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListOrganizationInvitations(context.Background(), organizationId).Status(status).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListOrganizationInvitations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizationInvitations`: ListInvitationsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListOrganizationInvitations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListOrganizationInvitationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **status** | **string** | Status of organizations | 

### Return type

[**ListInvitationsResponse**](ListInvitationsResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrganizations

> ListOrganizationResponse ListOrganizations(ctx).Execute()

List organizations of the connected user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListOrganizations(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListOrganizations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizations`: ListOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListOrganizations`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOrganizationsRequest struct via the builder pattern


### Return type

[**ListOrganizationResponse**](ListOrganizationResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrganizationsExpanded

> ListOrganizationExpandedResponse ListOrganizationsExpanded(ctx).Execute()

List organizations of the connected user with expanded data

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListOrganizationsExpanded(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListOrganizationsExpanded``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizationsExpanded`: ListOrganizationExpandedResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListOrganizationsExpanded`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOrganizationsExpandedRequest struct via the builder pattern


### Return type

[**ListOrganizationExpandedResponse**](ListOrganizationExpandedResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRegions

> ListRegionsResponse ListRegions(ctx, organizationId).Execute()

List regions

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListRegions(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListRegions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRegions`: ListRegionsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListRegions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListRegionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListRegionsResponse**](ListRegionsResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStacks

> ListStacksResponse ListStacks(ctx, organizationId).Deleted(deleted).Execute()

List stacks

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    deleted := true // bool | Include deleted stacks (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListStacks(context.Background(), organizationId).Deleted(deleted).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListStacks``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListStacks`: ListStacksResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListStacks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListStacksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleted** | **bool** | Include deleted stacks | 

### Return type

[**ListStacksResponse**](ListStacksResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUsers

> ListUsersResponse ListUsers(ctx, organizationId).Execute()

List users

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ListUsers(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListUsers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListUsers`: ListUsersResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListUsers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListUsersResponse**](ListUsersResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadConnectedUser

> ReadUserResponse ReadConnectedUser(ctx).Execute()

Read user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ReadConnectedUser(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReadConnectedUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadConnectedUser`: ReadUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReadConnectedUser`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReadConnectedUserRequest struct via the builder pattern


### Return type

[**ReadUserResponse**](ReadUserResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadOrganization

> CreateOrganizationResponse ReadOrganization(ctx, organizationId).Execute()

Read organization

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ReadOrganization(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReadOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadOrganization`: CreateOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReadOrganization`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreateOrganizationResponse**](CreateOrganizationResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadStack

> CreateStackResponse ReadStack(ctx, organizationId, stackId).Execute()

Read stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ReadStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReadStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReadStack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**CreateStackResponse**](CreateStackResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadUser

> ReadUserResponse ReadUser(ctx, organizationId, userId).Execute()

Read user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    userId := "userId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ReadUser(context.Background(), organizationId, userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReadUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadUser`: ReadUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReadUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ReadUserResponse**](ReadUserResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestoreStack

> CreateStackResponse RestoreStack(ctx, organizationId, stackId).Execute()

Restore stack

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    stackId := "stackId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.RestoreStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RestoreStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RestoreStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RestoreStack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestoreStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**CreateStackResponse**](CreateStackResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UnlinkUserFromOrganization

> UnlinkUserFromOrganization(ctx, organizationId, userId).Execute()

Unlink user from organization



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    userId := "userId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.UnlinkUserFromOrganization(context.Background(), organizationId, userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UnlinkUserFromOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUnlinkUserFromOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateConnectedUser

> UpdateUserResponse UpdateConnectedUser(ctx).Body(body).Execute()

Update user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    body := UpdatableUserData(987) // UpdatableUserData |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.UpdateConnectedUser(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UpdateConnectedUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateConnectedUser`: UpdateUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UpdateConnectedUser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateConnectedUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **UpdatableUserData** |  | 

### Return type

[**UpdateUserResponse**](UpdateUserResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUser

> UpdateUserResponse UpdateUser(ctx, organizationId, userId).Body(body).Execute()

Update user

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/formancehq/fctl/membershipclient"
)

func main() {
    organizationId := "organizationId_example" // string | 
    userId := "userId_example" // string | 
    body := UpdatableUserData(987) // UpdatableUserData |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.UpdateUser(context.Background(), organizationId, userId).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UpdateUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateUser`: UpdateUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UpdateUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **UpdatableUserData** |  | 

### Return type

[**UpdateUserResponse**](UpdateUserResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

