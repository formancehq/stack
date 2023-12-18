# \DefaultApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptInvitation**](DefaultApi.md#AcceptInvitation) | **Post** /me/invitations/{invitationId}/accept | Accept invitation
[**CreateInvitation**](DefaultApi.md#CreateInvitation) | **Post** /organizations/{organizationId}/invitations | Create invitation
[**CreateOrganization**](DefaultApi.md#CreateOrganization) | **Post** /organizations | Create organization
[**CreatePrivateRegion**](DefaultApi.md#CreatePrivateRegion) | **Post** /organizations/{organizationId}/regions | Create a private region
[**CreateStack**](DefaultApi.md#CreateStack) | **Post** /organizations/{organizationId}/stacks | Create stack
[**DeclineInvitation**](DefaultApi.md#DeclineInvitation) | **Post** /me/invitations/{invitationId}/reject | Decline invitation
[**DeleteInvitation**](DefaultApi.md#DeleteInvitation) | **Delete** /organizations/{organizationId}/invitations/{invitationId} | Delete invitation
[**DeleteOrganization**](DefaultApi.md#DeleteOrganization) | **Delete** /organizations/{organizationId} | Delete organization
[**DeleteRegion**](DefaultApi.md#DeleteRegion) | **Delete** /organizations/{organizationId}/regions/{regionId} | Delete region
[**DeleteStack**](DefaultApi.md#DeleteStack) | **Delete** /organizations/{organizationId}/stacks/{stackId} | Delete stack
[**DeleteStackUserAccess**](DefaultApi.md#DeleteStackUserAccess) | **Delete** /organizations/{organizationId}/stacks/{stackId}/users/{userId} | Delete stack user access role within an organization
[**DeleteUserFromOrganization**](DefaultApi.md#DeleteUserFromOrganization) | **Delete** /organizations/{organizationId}/users/{userId} | delete user from organization
[**DisableStack**](DefaultApi.md#DisableStack) | **Put** /organizations/{organizationId}/stacks/{stackId}/disable | Disable stack
[**DisableStargate**](DefaultApi.md#DisableStargate) | **Put** /organizations/{organizationId}/stacks/{stackId}/stargate/disable | Disable stargate on a stack
[**EnableStack**](DefaultApi.md#EnableStack) | **Put** /organizations/{organizationId}/stacks/{stackId}/enable | Enable stack
[**EnableStargate**](DefaultApi.md#EnableStargate) | **Put** /organizations/{organizationId}/stacks/{stackId}/stargate/enable | Enable stargate on a stack
[**GetRegion**](DefaultApi.md#GetRegion) | **Get** /organizations/{organizationId}/regions/{regionId} | Get region
[**GetRegionVersions**](DefaultApi.md#GetRegionVersions) | **Get** /organizations/{organizationId}/regions/{regionId}/versions | Get region versions
[**GetServerInfo**](DefaultApi.md#GetServerInfo) | **Get** /_info | Get server info
[**GetStack**](DefaultApi.md#GetStack) | **Get** /organizations/{organizationId}/stacks/{stackId} | Find stack
[**ListInvitations**](DefaultApi.md#ListInvitations) | **Get** /me/invitations | List invitations of the user
[**ListOrganizationInvitations**](DefaultApi.md#ListOrganizationInvitations) | **Get** /organizations/{organizationId}/invitations | List invitations of the organization
[**ListOrganizations**](DefaultApi.md#ListOrganizations) | **Get** /organizations | List organizations of the connected user
[**ListOrganizationsExpanded**](DefaultApi.md#ListOrganizationsExpanded) | **Get** /organizations/expanded | List organizations of the connected user with expanded data
[**ListRegions**](DefaultApi.md#ListRegions) | **Get** /organizations/{organizationId}/regions | List regions
[**ListStackUsersAccesses**](DefaultApi.md#ListStackUsersAccesses) | **Get** /organizations/{organizationId}/stacks/{stackId}/users | List stack users accesses within an organization
[**ListStacks**](DefaultApi.md#ListStacks) | **Get** /organizations/{organizationId}/stacks | List stacks
[**ListUsersOfOrganization**](DefaultApi.md#ListUsersOfOrganization) | **Get** /organizations/{organizationId}/users | List users of organization
[**ReadConnectedUser**](DefaultApi.md#ReadConnectedUser) | **Get** /me | Read user
[**ReadOrganization**](DefaultApi.md#ReadOrganization) | **Get** /organizations/{organizationId} | Read organization
[**ReadUserOfOrganization**](DefaultApi.md#ReadUserOfOrganization) | **Get** /organizations/{organizationId}/users/{userId} | Read user of organization
[**RestoreStack**](DefaultApi.md#RestoreStack) | **Put** /organizations/{organizationId}/stacks/{stackId}/restore | Restore stack
[**UpsertOrganizationUser**](DefaultApi.md#UpsertOrganizationUser) | **Put** /organizations/{organizationId}/users/{userId} | Update user role within an organization
[**UpsertStackUserAccess**](DefaultApi.md#UpsertStackUserAccess) | **Put** /organizations/{organizationId}/stacks/{stackId}/users/{userId} | Update stack user access role within an organization



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
    r, err := apiClient.DefaultApi.AcceptInvitation(context.Background(), invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.AcceptInvitation``: %v\n", err)
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
    resp, r, err := apiClient.DefaultApi.CreateInvitation(context.Background(), organizationId).Email(email).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateInvitation``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateInvitation`: CreateInvitationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateInvitation`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.CreateOrganization(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOrganization`: CreateOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateOrganization`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.CreatePrivateRegion(context.Background(), organizationId).CreatePrivateRegionRequest(createPrivateRegionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreatePrivateRegion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreatePrivateRegion`: CreatedPrivateRegionResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreatePrivateRegion`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.CreateStack(context.Background(), organizationId).CreateStackRequest(createStackRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateStack`: %v\n", resp)
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
    r, err := apiClient.DefaultApi.DeclineInvitation(context.Background(), invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeclineInvitation``: %v\n", err)
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
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DeleteInvitation(context.Background(), organizationId, invitationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteInvitation``: %v\n", err)
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
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DeleteOrganization(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteOrganization``: %v\n", err)
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
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DeleteRegion(context.Background(), organizationId, regionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteRegion``: %v\n", err)
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
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DeleteStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteStack``: %v\n", err)
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
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteStackUserAccess

> DeleteStackUserAccess(ctx, organizationId, stackId, userId).Execute()

Delete stack user access role within an organization

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
    userId := "userId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultApi.DeleteStackUserAccess(context.Background(), organizationId, stackId, userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteStackUserAccess``: %v\n", err)
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
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteStackUserAccessRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUserFromOrganization

> DeleteUserFromOrganization(ctx, organizationId, userId).Execute()

delete user from organization



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
    r, err := apiClient.DefaultApi.DeleteUserFromOrganization(context.Background(), organizationId, userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeleteUserFromOrganization``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteUserFromOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DisableStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DisableStack``: %v\n", err)
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
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.DisableStargate(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DisableStargate``: %v\n", err)
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
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EnableStack

> EnableStack(ctx, organizationId, stackId).Execute()

Enable stack

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
    r, err := apiClient.DefaultApi.EnableStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.EnableStack``: %v\n", err)
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

Other parameters are passed through a pointer to a apiEnableStackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

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
    r, err := apiClient.DefaultApi.EnableStargate(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.EnableStargate``: %v\n", err)
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
- **Accept**: application/json

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
    resp, r, err := apiClient.DefaultApi.GetRegion(context.Background(), organizationId, regionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetRegion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRegion`: GetRegionResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetRegion`: %v\n", resp)
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


## GetRegionVersions

> GetRegionVersionsResponse GetRegionVersions(ctx, organizationId, regionId).Execute()

Get region versions

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
    resp, r, err := apiClient.DefaultApi.GetRegionVersions(context.Background(), organizationId, regionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetRegionVersions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRegionVersions`: GetRegionVersionsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetRegionVersions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**regionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRegionVersionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetRegionVersionsResponse**](GetRegionVersionsResponse.md)

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
    resp, r, err := apiClient.DefaultApi.GetServerInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetServerInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetServerInfo`: ServerInfo
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetServerInfo`: %v\n", resp)
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


## GetStack

> CreateStackResponse GetStack(ctx, organizationId, stackId).Execute()

Find stack

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
    resp, r, err := apiClient.DefaultApi.GetStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetStack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStackRequest struct via the builder pattern


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
    resp, r, err := apiClient.DefaultApi.ListInvitations(context.Background()).Status(status).Organization(organization).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListInvitations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListInvitations`: ListInvitationsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListInvitations`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.ListOrganizationInvitations(context.Background(), organizationId).Status(status).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListOrganizationInvitations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizationInvitations`: ListInvitationsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListOrganizationInvitations`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.ListOrganizations(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListOrganizations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizations`: ListOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListOrganizations`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.ListOrganizationsExpanded(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListOrganizationsExpanded``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrganizationsExpanded`: ListOrganizationExpandedResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListOrganizationsExpanded`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.ListRegions(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListRegions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRegions`: ListRegionsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListRegions`: %v\n", resp)
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


## ListStackUsersAccesses

> StackUserAccessResponse ListStackUsersAccesses(ctx, organizationId, stackId).Execute()

List stack users accesses within an organization

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
    resp, r, err := apiClient.DefaultApi.ListStackUsersAccesses(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListStackUsersAccesses``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListStackUsersAccesses`: StackUserAccessResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListStackUsersAccesses`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**stackId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListStackUsersAccessesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**StackUserAccessResponse**](StackUserAccessResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStacks

> ListStacksResponse ListStacks(ctx, organizationId).All(all).Deleted(deleted).Execute()

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
    all := true // bool | Include deleted and disabled stacks (optional)
    deleted := true // bool | Include deleted stacks (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.ListStacks(context.Background(), organizationId).All(all).Deleted(deleted).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListStacks``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListStacks`: ListStacksResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListStacks`: %v\n", resp)
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

 **all** | **bool** | Include deleted and disabled stacks | 
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


## ListUsersOfOrganization

> ListUsersResponse ListUsersOfOrganization(ctx, organizationId).Execute()

List users of organization

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
    resp, r, err := apiClient.DefaultApi.ListUsersOfOrganization(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ListUsersOfOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListUsersOfOrganization`: ListUsersResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ListUsersOfOrganization`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListUsersOfOrganizationRequest struct via the builder pattern


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
    resp, r, err := apiClient.DefaultApi.ReadConnectedUser(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ReadConnectedUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadConnectedUser`: ReadUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ReadConnectedUser`: %v\n", resp)
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
    resp, r, err := apiClient.DefaultApi.ReadOrganization(context.Background(), organizationId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ReadOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadOrganization`: CreateOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ReadOrganization`: %v\n", resp)
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


## ReadUserOfOrganization

> ReadUserResponse ReadUserOfOrganization(ctx, organizationId, userId).Execute()

Read user of organization

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
    resp, r, err := apiClient.DefaultApi.ReadUserOfOrganization(context.Background(), organizationId, userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.ReadUserOfOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadUserOfOrganization`: ReadUserResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.ReadUserOfOrganization`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadUserOfOrganizationRequest struct via the builder pattern


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
    resp, r, err := apiClient.DefaultApi.RestoreStack(context.Background(), organizationId, stackId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.RestoreStack``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RestoreStack`: CreateStackResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.RestoreStack`: %v\n", resp)
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


## UpsertOrganizationUser

> UpsertOrganizationUser(ctx, organizationId, userId).RequestBody(requestBody).Execute()

Update user role within an organization

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
    requestBody := []string{"Property_example"} // []string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultApi.UpsertOrganizationUser(context.Background(), organizationId, userId).RequestBody(requestBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UpsertOrganizationUser``: %v\n", err)
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

Other parameters are passed through a pointer to a apiUpsertOrganizationUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **requestBody** | **[]string** |  | 

### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpsertStackUserAccess

> UpsertStackUserAccess(ctx, organizationId, stackId, userId).RequestBody(requestBody).Execute()

Update stack user access role within an organization

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
    userId := "userId_example" // string | 
    requestBody := []string{"Property_example"} // []string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultApi.UpsertStackUserAccess(context.Background(), organizationId, stackId, userId).RequestBody(requestBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UpsertStackUserAccess``: %v\n", err)
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
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpsertStackUserAccessRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **requestBody** | **[]string** |  | 

### Return type

 (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

