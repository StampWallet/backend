# UserApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**searchBusinesses**](UserApi.md#searchBusinesses) | **GET** /user/businesses | Search businesses



## searchBusinesses

Search businesses

This endpoint is used to search businesses that match the provided text query or are close to a specified point.

### Example

```bash
 searchBusinesses  text=value  location=value  proximity=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **text** | **string** | Filter by business name | [optional] [default to null]
 **location** | **string** | Filter by business location | [optional] [default to null]
 **proximity** | **integer** | Filter by distance from location in meters | [optional] [default to null]

### Return type

[**GetUserBusinessesSearchResponse**](GetUserBusinessesSearchResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

