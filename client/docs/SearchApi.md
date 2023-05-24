# \SearchApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SearchBusinesses**](SearchApi.md#SearchBusinesses) | **Get** /user/businesses | Search businesses



## SearchBusinesses

> GetUserBusinessesSearchResponse SearchBusinesses(ctx).Text(text).Location(location).Proximity(proximity).Execute()

Search businesses



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    text := "text_example" // string | Filter by business name (optional)
    location := "location_example" // string | Filter by business location (optional)
    proximity := int32(56) // int32 | Filter by distance from location in meters (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.SearchApi.SearchBusinesses(context.Background()).Text(text).Location(location).Proximity(proximity).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SearchApi.SearchBusinesses``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SearchBusinesses`: GetUserBusinessesSearchResponse
    fmt.Fprintf(os.Stdout, "Response from `SearchApi.SearchBusinesses`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSearchBusinessesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **text** | **string** | Filter by business name | 
 **location** | **string** | Filter by business location | 
 **proximity** | **int32** | Filter by distance from location in meters | 

### Return type

[**GetUserBusinessesSearchResponse**](GetUserBusinessesSearchResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

