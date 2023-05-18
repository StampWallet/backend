# ItemDefinitionsApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**addItemDefinition**](ItemDefinitionsApi.md#addItemDefinition) | **POST** /business/itemDefinitions | Add a new item definition
[**deleteItemDefinition**](ItemDefinitionsApi.md#deleteItemDefinition) | **DELETE** /business/itemDefinitions/{definitionId} | Delete an exiting item definition
[**getItemDefinitions**](ItemDefinitionsApi.md#getItemDefinitions) | **GET** /business/itemDefinitions | Get list of item definitions
[**updateItemDefinition**](ItemDefinitionsApi.md#updateItemDefinition) | **PATCH** /business/itemDefinitions/{definitionId} | Update an exiting item definition



## addItemDefinition

Add a new item definition

This endpoint is used to add new item definitions (benefits).

### Example

```bash
 addItemDefinition
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postBusinessItemDefinitionRequest** | [**PostBusinessItemDefinitionRequest**](PostBusinessItemDefinitionRequest.md) |  |

### Return type

[**PostBusinessItemDefinitionResponse**](PostBusinessItemDefinitionResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteItemDefinition

Delete an exiting item definition

This endpoint is used to delete existing item definitions (benefits).

### Example

```bash
 deleteItemDefinition definitionId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **definitionId** | **string** | Public id of the definition to update | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getItemDefinitions

Get list of item definitions

This endpoint is used to retrieve data about existing item definitions (benefits).

### Example

```bash
 getItemDefinitions
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetBusinessItemDefinitionsResponse**](GetBusinessItemDefinitionsResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## updateItemDefinition

Update an exiting item definition

This endpoint is used to change details of existing item definitions (benefits).

### Example

```bash
 updateItemDefinition definitionId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **definitionId** | **string** | Public id of the definition to update | [default to null]
 **patchBusinessItemDefinitionRequest** | [**PatchBusinessItemDefinitionRequest**](PatchBusinessItemDefinitionRequest.md) |  |

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

