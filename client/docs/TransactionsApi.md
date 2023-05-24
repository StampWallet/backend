# TransactionsApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**finishTransaction**](TransactionsApi.md#finishTransaction) | **POST** /business/transactions/{transactionCode} | Finish a transaction
[**getTransactionDetails**](TransactionsApi.md#getTransactionDetails) | **GET** /business/transactions/{transactionCode} | Get info about a started transaction



## finishTransaction

Finish a transaction

This endpoint is used in the third step of transaction processing, the app should use it to update transaction details with data about points added to user's account and actions that were taken on items included in the transaction.

### Example

```bash
 finishTransaction transactionCode=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **transactionCode** | **string** | Transaction code (scanned or typed in) | [default to null]
 **postBusinessTransactionRequest** | [**PostBusinessTransactionRequest**](PostBusinessTransactionRequest.md) |  |

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## getTransactionDetails

Get info about a started transaction

This endpoint is used in the second step of transaction processing, the app should use it to retrieve details about a transaction started by a user, after scanning user's transaction code.

### Example

```bash
 getTransactionDetails transactionCode=value
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **transactionCode** | **string** | Transaction code (scanned or typed in) | [default to null]

### Return type

[**GetBusinessTransactionResponse**](GetBusinessTransactionResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: Not Applicable
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

