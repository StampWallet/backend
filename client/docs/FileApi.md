# FileApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteFile**](FileApi.md#deleteFile) | **DELETE** /file/{fileId} | Delete a file
[**getFile**](FileApi.md#getFile) | **GET** /file/{fileId} | Get file
[**uploadFile**](FileApi.md#uploadFile) | **POST** /file/{fileId} | Upload file



## deleteFile

Delete a file

This endpoint is used to delete files.

### Example

```bash
 deleteFile fileId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fileId** | **string** | ID of file to delete | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getFile

Get file

This endpoint is used to download files by ID.

### Example

```bash
 getFile fileId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fileId** | **string** | ID of file to download | [default to null]

### Return type

**binary**

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: image/png, image/jpeg, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## uploadFile

Upload file

This endpoint is used to upload files.

### Example

```bash
 uploadFile fileId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fileId** | **string** | ID of file to upload/replace | [default to null]
 **file** | **binary** |  | [optional] [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

