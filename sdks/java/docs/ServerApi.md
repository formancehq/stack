# ServerApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getInfo**](ServerApi.md#getInfo) | **GET** api/ledger/_info | Show server information |



## getInfo

> ConfigInfo getInfo()

Show server information

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.models.*;
import com.formance.formance.api.ServerApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        ServerApi apiInstance = new ServerApi(defaultClient);
        try {
            ConfigInfo result = apiInstance.getInfo();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling ServerApi#getInfo");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ConfigInfo**](ConfigInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

