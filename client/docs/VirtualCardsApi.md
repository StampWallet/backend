# VirtualCardsApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**buyItem**](VirtualCardsApi.md#buyItem) | **POST** /user/cards/virtual/{businessId}/items/{itemDefinitionId} | Buy an item
[**createVirtualCard**](VirtualCardsApi.md#createVirtualCard) | **POST** /user/cards/virtual/{businessId} | Add a new virtual card
[**deleteItem**](VirtualCardsApi.md#deleteItem) | **DELETE** /user/cards/virtual/{businessId}/items/{itemId} | Delete an item
[**deleteVirtualCard**](VirtualCardsApi.md#deleteVirtualCard) | **DELETE** /user/cards/virtual/{businessId} | Delete a virtual card
[**getVirtualCard**](VirtualCardsApi.md#getVirtualCard) | **GET** /user/cards/virtual/{businessId} | Get info about a virtual card



## buyItem

Buy an item

This endpoint is used to buy an item for points from the virtual card.

### Example

```bash
 buyItem businessId=value itemDefinitionId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested | [default to null]
 **itemDefinitionId** | **string** | Public ID of the item definition requested by the user | [default to null]

### Return type

[**PostUserVirtualCardItemResponse**](PostUserVirtualCardItemResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## createVirtualCard

Add a new virtual card

This endpoint is used to register a new virtual card to the account of the currently logged in user.

### Example

```bash
 createVirtualCard businessId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested to be added by user | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteItem

Delete an item

This endpoint is used to return an item, and get back points that were spent on that item.

### Example

```bash
 deleteItem businessId=value itemId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested | [default to null]
 **itemId** | **string** | Public ID of the item requested to be deleted | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteVirtualCard

Delete a virtual card

This endpoint is used to delete a virtual card from the account of the currently logged in user.

### Example

```bash
 deleteVirtualCard businessId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested to be deleted from the account | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getVirtualCard

Get info about a virtual card

This endpoint is used to retrieve details of a virtual card owned by the currently logged in user.

### Example

```bash
 getVirtualCard businessId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested | [default to null]

### Return type

[**GetUserVirtualCardResponse**](GetUserVirtualCardResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

