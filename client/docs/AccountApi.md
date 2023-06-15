# AccountApi

All URIs are relative to **

Method | HTTP request | Description
------------- | ------------- | -------------
[**changeEmail**](AccountApi.md#changeEmail) | **POST** /auth/account/email | Change email
[**changePassword**](AccountApi.md#changePassword) | **POST** /auth/account/password | Change password
[**confirmEmail**](AccountApi.md#confirmEmail) | **POST** /auth/account/emailConfirmation | Confirm email
[**createAccount**](AccountApi.md#createAccount) | **POST** /auth/account | Create a new account



## changeEmail

Change email

This endpoint can be used to change email address of currently logged in user. Changing email address requires email confirmation

### Example

```bash
 changeEmail
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postAccountEmailRequest** | [**PostAccountEmailRequest**](PostAccountEmailRequest.md) |  |

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## changePassword

Change password

This endpoint can be used to change password of currently logged in user. Requires the user to provide their old password

### Example

```bash
 changePassword
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postAccountPasswordRequest** | [**PostAccountPasswordRequest**](PostAccountPasswordRequest.md) |  |

### Return type

[**DefaultResponse**](DefaultResponse.md)

### Authorization

[sessionToken](../README.md#sessionToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## confirmEmail

Confirm email

When an account is created, user receives an email with a link to a static website. That website simply posts it's parameters (unique to each email) to this endpoint. The parameters will be unique and hard to guess, allowing to verify that user really has access to the email address.

### Example

```bash
 confirmEmail
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postAccountEmailConfirmationRequest** | [**PostAccountEmailConfirmationRequest**](PostAccountEmailConfirmationRequest.md) |  |

### Return type

[**PostAccountResponse**](PostAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## createAccount

Create a new account

Create a new account with specified password and email, send a confirmation email

### Example

```bash
 createAccount
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **postAccountRequest** | [**PostAccountRequest**](PostAccountRequest.md) |  |

### Return type

[**PostAccountResponse**](PostAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

