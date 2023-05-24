# TransactionApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**getTransactionStatus**](TransactionApi.md#getTransactionStatus) | **GET** /user/cards/virtual/{businessId}/transactions/{transactionCode} | Get info about a transaction
[**startTransaction**](TransactionApi.md#startTransaction) | **POST** /user/cards/virtual/{businessId}/transactions | Start a transaction



## getTransactionStatus

Get info about a transaction

This endpoint is used in the last step of transaction processing, it's used to check the status of the transaction.

### Example

```bash
 getTransactionStatus businessId=value transactionCode=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested | [default to null]
 **transactionCode** | **string** | Transaction code | [default to null]

### Return type

[**GetUserVirtualCardTransactionResponse**](GetUserVirtualCardTransactionResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## startTransaction

Start a transaction

This endpoint is used in the first step of transaction processing, the app should use it to start a transaction optionally providing items to be exchanged.

### Example

```bash
 startTransaction businessId=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **businessId** | **string** | Public ID of the business which card was requested | [default to null]
 **postUserVirtualCardTransactionRequest** | [**PostUserVirtualCardTransactionRequest**](PostUserVirtualCardTransactionRequest.md) |  |

### Return type

[**PostUserVirtualCardTransactionResponse**](PostUserVirtualCardTransactionResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

