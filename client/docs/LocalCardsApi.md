# LocalCardsApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**createLocalCard**](LocalCardsApi.md#createLocalCard) | **POST** /user/cards/local | Add a new local card
[**deleteLocalCard**](LocalCardsApi.md#deleteLocalCard) | **DELETE** /user/cards/local/{cardId} | Delete a local card
[**getLocalCardTypes**](LocalCardsApi.md#getLocalCardTypes) | **GET** /user/cards/local/types | Get list of local card types



## createLocalCard

Add a new local card

This endpoint is used to add a new local card to user's account.

### Example

```bash
 createLocalCard
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postUserLocalCardsRequest** | [**PostUserLocalCardsRequest**](PostUserLocalCardsRequest.md) |  |

### Return type

[**PostUserLocalCardsResponse**](PostUserLocalCardsResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## deleteLocalCard

Delete a local card

This endpoint is used to delete a local card from account of the currently logged in user.

### Example

```bash
 deleteLocalCard cardId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cardId** | **string** | Public id of the card to delete | [default to null]

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getLocalCardTypes

Get list of local card types

This endpoint is used to get a list of supported local card types.

### Example

```bash
 getLocalCardTypes
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetUserLocalCardTypesResponse**](GetUserLocalCardTypesResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

