# BusinessApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMenuImage**](BusinessApi.md#addMenuImage) | **POST** /business/menuImages/ | Add menu image to business
[**createBusinessAccount**](BusinessApi.md#createBusinessAccount) | **POST** /business/account | Create a business account
[**deleteMenuImage**](BusinessApi.md#deleteMenuImage) | **DELETE** /business/menuImages/{menuImageId} | Delete menu image from business
[**getBusinessAccountInfo**](BusinessApi.md#getBusinessAccountInfo) | **GET** /business/info | Get business info
[**updateBusinessAccount**](BusinessApi.md#updateBusinessAccount) | **PATCH** /business/info | Update business account



## addMenuImage

Add menu image to business

This endpoint is used to add a new menu image to business details. Returns a new fileId to be used with '/file/' endpoints.

### Example

```bash
 addMenuImage
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**PostBusinessAccountMenuImageResponse**](PostBusinessAccountMenuImageResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## createBusinessAccount

Create a business account

This endpoint is used to attach a new business account to an existing, logged in user account. Busies details are provided in the request. Responds with business id and ids of banner and icon image slots.

### Example

```bash
 createBusinessAccount
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postBusinessAccountRequest** | [**PostBusinessAccountRequest**](PostBusinessAccountRequest.md) |  |

### Return type

[**PostBusinessAccountResponse**](PostBusinessAccountResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteMenuImage

Delete menu image from business

This endpoint is used to delete a menu image from business details.

### Example

```bash
 deleteMenuImage menuImageId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **menuImageId** | **string** | Public id of the menu image to be deleted | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getBusinessAccountInfo

Get business info

Responds with information about business owned by the logged in user.

### Example

```bash
 getBusinessAccountInfo
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetBusinessAccountResponse**](GetBusinessAccountResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## updateBusinessAccount

Update business account

This endpoint is used to update business account data

### Example

```bash
 updateBusinessAccount
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **patchBusinessAccountRequest** | [**PatchBusinessAccountRequest**](PatchBusinessAccountRequest.md) |  |

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

