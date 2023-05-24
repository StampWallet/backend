# SessionsApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**login**](SessionsApi.md#login) | **POST** /auth/sessions | Login
[**logout**](SessionsApi.md#logout) | **DELETE** /auth/sessions | Logout



## login

Login

This endpoint is used to exchange user credentials for temporary credentials that allow access to the API.

### Example

```bash
 login
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postAccountSessionRequest** | [**PostAccountSessionRequest**](PostAccountSessionRequest.md) |  |

### Return type

[**PostAccountSessionResponse**](PostAccountSessionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## logout

Logout

This endpoint invalidates session token passed with the request.

### Example

```bash
 logout
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

