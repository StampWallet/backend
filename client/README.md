# StampWallet API Server Bash client

## Overview

This is a Bash client script for accessing StampWallet API Server service.

The script uses cURL underneath for making all REST calls.

## Usage

```shell
# Make sure the script has executable rights
$ chmod u+x 

# Print the list of operations available on the service
$ ./ -h

# Print the service description
$ ./ --about

# Print detailed information about specific operation
$ ./ <operationId> -h

# Make GET request
./ --host http://<hostname>:<port> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make GET request using arbitrary curl options (must be passed before <operationId>) to an SSL service using username:password
 -k -sS --tlsv1.2 --host https://<hostname> -u <user>:<password> --accept xml <operationId> <queryParam1>=<value1> <header_key1>:<header_value2>

# Make POST request
$ echo '<body_content>' |  --host <hostname> --content-type json <operationId> -

# Make POST request with simple JSON content, e.g.:
# {
#   "key1": "value1",
#   "key2": "value2",
#   "key3": 23
# }
$ echo '<body_content>' |  --host <hostname> --content-type json <operationId> key1==value1 key2=value2 key3:=23 -

# Make POST request with form data
$  --host <hostname> <operationId> key1:=value1 key2:=value2 key3:=23

# Preview the cURL command without actually executing it
$  --host http://<hostname>:<port> --dry-run <operationid>

```

## Docker image

You can easily create a Docker image containing a preconfigured environment
for using the REST Bash client including working autocompletion and short
welcome message with basic instructions, using the generated Dockerfile:

```shell
docker build -t my-rest-client .
docker run -it my-rest-client
```

By default you will be logged into a Zsh environment which has much more
advanced auto completion, but you can switch to Bash, where basic autocompletion
is also available.

## Shell completion

### Bash

The generated bash-completion script can be either directly loaded to the current Bash session using:

```shell
source .bash-completion
```

Alternatively, the script can be copied to the `/etc/bash-completion.d` (or on OSX with Homebrew to `/usr/local/etc/bash-completion.d`):

```shell
sudo cp .bash-completion /etc/bash-completion.d/
```

#### OS X

On OSX you might need to install bash-completion using Homebrew:

```shell
brew install bash-completion
```

and add the following to the `~/.bashrc`:

```shell
if [ -f $(brew --prefix)/etc/bash_completion ]; then
  . $(brew --prefix)/etc/bash_completion
fi
```

### Zsh

In Zsh, the generated `_` Zsh completion file must be copied to one of the folders under `$FPATH` variable.

## Documentation for API Endpoints

All URIs are relative to **

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AccountApi* | [**changeEmail**](docs/AccountApi.md#changeemail) | **POST** /auth/account/email | Change email
*AccountApi* | [**changePassword**](docs/AccountApi.md#changepassword) | **POST** /auth/account/password | Change password
*AccountApi* | [**confirmEmail**](docs/AccountApi.md#confirmemail) | **POST** /auth/account/emailConfirmation | Confirm email
*AccountApi* | [**createAccount**](docs/AccountApi.md#createaccount) | **POST** /auth/account | Create a new account
*BusinessApi* | [**addMenuImage**](docs/BusinessApi.md#addmenuimage) | **POST** /business/menuImages/ | Add menu image to business
*BusinessApi* | [**createBusinessAccount**](docs/BusinessApi.md#createbusinessaccount) | **POST** /business/account | Create a business account
*BusinessApi* | [**deleteMenuImage**](docs/BusinessApi.md#deletemenuimage) | **DELETE** /business/menuImages/{menuImageId} | Delete menu image from business
*BusinessApi* | [**getBusinessAccountInfo**](docs/BusinessApi.md#getbusinessaccountinfo) | **GET** /business/info | Get business info
*BusinessApi* | [**updateBusinessAccount**](docs/BusinessApi.md#updatebusinessaccount) | **PATCH** /business/info | Update business account
*CardsApi* | [**getUserCards**](docs/CardsApi.md#getusercards) | **GET** /user/cards | Get list of user&#39;s cards
*FileApi* | [**deleteFile**](docs/FileApi.md#deletefile) | **DELETE** /file/{fileId} | Delete a file
*FileApi* | [**getFile**](docs/FileApi.md#getfile) | **GET** /file/{fileId} | Get file
*FileApi* | [**uploadFile**](docs/FileApi.md#uploadfile) | **POST** /file/{fileId} | Upload file
*ItemDefinitionsApi* | [**addItemDefinition**](docs/ItemDefinitionsApi.md#additemdefinition) | **POST** /business/itemDefinitions | Add a new item definition
*ItemDefinitionsApi* | [**deleteItemDefinition**](docs/ItemDefinitionsApi.md#deleteitemdefinition) | **DELETE** /business/itemDefinitions/{definitionId} | Delete an exiting item definition
*ItemDefinitionsApi* | [**getItemDefinitions**](docs/ItemDefinitionsApi.md#getitemdefinitions) | **GET** /business/itemDefinitions | Get list of item definitions
*ItemDefinitionsApi* | [**updateItemDefinition**](docs/ItemDefinitionsApi.md#updateitemdefinition) | **PATCH** /business/itemDefinitions/{definitionId} | Update an exiting item definition
*LocalCardsApi* | [**createLocalCard**](docs/LocalCardsApi.md#createlocalcard) | **POST** /user/cards/local | Add a new local card
*LocalCardsApi* | [**deleteLocalCard**](docs/LocalCardsApi.md#deletelocalcard) | **DELETE** /user/cards/local/{cardId} | Delete a local card
*LocalCardsApi* | [**getLocalCardTypes**](docs/LocalCardsApi.md#getlocalcardtypes) | **GET** /user/cards/local/types | Get list of local card types
*SessionsApi* | [**login**](docs/SessionsApi.md#login) | **POST** /auth/sessions | Login
*SessionsApi* | [**logout**](docs/SessionsApi.md#logout) | **DELETE** /auth/sessions | Logout
*TransactionApi* | [**getTransactionStatus**](docs/TransactionApi.md#gettransactionstatus) | **GET** /user/cards/virtual/{businessId}/transactions/{transactionCode} | Get info about a transaction
*TransactionApi* | [**startTransaction**](docs/TransactionApi.md#starttransaction) | **POST** /user/cards/virtual/{businessId}/transactions | Start a transaction
*TransactionsApi* | [**finishTransaction**](docs/TransactionsApi.md#finishtransaction) | **POST** /business/transactions/{transactionCode} | Finish a transaction
*TransactionsApi* | [**getTransactionDetails**](docs/TransactionsApi.md#gettransactiondetails) | **GET** /business/transactions/{transactionCode} | Get info about a started transaction
*UserApi* | [**getBusiness**](docs/UserApi.md#getbusiness) | **GET** /user/businesses/{businessId} | Get business info
*UserApi* | [**searchBusinesses**](docs/UserApi.md#searchbusinesses) | **GET** /user/businesses | Search businesses
*VirtualCardsApi* | [**buyItem**](docs/VirtualCardsApi.md#buyitem) | **POST** /user/cards/virtual/{businessId}/itemsDefinitions/{itemDefinitionId} | Buy an item
*VirtualCardsApi* | [**createVirtualCard**](docs/VirtualCardsApi.md#createvirtualcard) | **POST** /user/cards/virtual/{businessId} | Add a new virtual card
*VirtualCardsApi* | [**deleteItem**](docs/VirtualCardsApi.md#deleteitem) | **DELETE** /user/cards/virtual/{businessId}/items/{itemId} | Delete an item
*VirtualCardsApi* | [**deleteVirtualCard**](docs/VirtualCardsApi.md#deletevirtualcard) | **DELETE** /user/cards/virtual/{businessId} | Delete a virtual card
*VirtualCardsApi* | [**getVirtualCard**](docs/VirtualCardsApi.md#getvirtualcard) | **GET** /user/cards/virtual/{businessId} | Get info about a virtual card


## Documentation For Models

 - [DefaultResponse](docs/DefaultResponse.md)
 - [DefaultResponseStatusEnum](docs/DefaultResponseStatusEnum.md)
 - [GetBusinessAccountResponse](docs/GetBusinessAccountResponse.md)
 - [GetBusinessAccountResponseAllOf](docs/GetBusinessAccountResponseAllOf.md)
 - [GetBusinessItemDefinitionsResponse](docs/GetBusinessItemDefinitionsResponse.md)
 - [GetBusinessTransactionResponse](docs/GetBusinessTransactionResponse.md)
 - [GetUserBusinessesSearchResponse](docs/GetUserBusinessesSearchResponse.md)
 - [GetUserCardsResponse](docs/GetUserCardsResponse.md)
 - [GetUserLocalCardTypesResponse](docs/GetUserLocalCardTypesResponse.md)
 - [GetUserLocalCardTypesResponseTypesInner](docs/GetUserLocalCardTypesResponseTypesInner.md)
 - [GetUserVirtualCardResponse](docs/GetUserVirtualCardResponse.md)
 - [GetUserVirtualCardTransactionResponse](docs/GetUserVirtualCardTransactionResponse.md)
 - [ItemActionAPIModel](docs/ItemActionAPIModel.md)
 - [ItemActionTypeEnum](docs/ItemActionTypeEnum.md)
 - [ItemDefinitionAPIModel](docs/ItemDefinitionAPIModel.md)
 - [LocalCardAPIModel](docs/LocalCardAPIModel.md)
 - [OwnedItemAPIModel](docs/OwnedItemAPIModel.md)
 - [PatchBusinessAccountRequest](docs/PatchBusinessAccountRequest.md)
 - [PatchBusinessItemDefinitionRequest](docs/PatchBusinessItemDefinitionRequest.md)
 - [PostAccountEmailConfirmationRequest](docs/PostAccountEmailConfirmationRequest.md)
 - [PostAccountEmailRequest](docs/PostAccountEmailRequest.md)
 - [PostAccountPasswordRequest](docs/PostAccountPasswordRequest.md)
 - [PostAccountRequest](docs/PostAccountRequest.md)
 - [PostAccountResponse](docs/PostAccountResponse.md)
 - [PostAccountSessionRequest](docs/PostAccountSessionRequest.md)
 - [PostAccountSessionResponse](docs/PostAccountSessionResponse.md)
 - [PostBusinessAccountMenuImageResponse](docs/PostBusinessAccountMenuImageResponse.md)
 - [PostBusinessAccountRequest](docs/PostBusinessAccountRequest.md)
 - [PostBusinessAccountResponse](docs/PostBusinessAccountResponse.md)
 - [PostBusinessItemDefinitionRequest](docs/PostBusinessItemDefinitionRequest.md)
 - [PostBusinessItemDefinitionResponse](docs/PostBusinessItemDefinitionResponse.md)
 - [PostBusinessTransactionRequest](docs/PostBusinessTransactionRequest.md)
 - [PostUserLocalCardsRequest](docs/PostUserLocalCardsRequest.md)
 - [PostUserLocalCardsResponse](docs/PostUserLocalCardsResponse.md)
 - [PostUserVirtualCardItemResponse](docs/PostUserVirtualCardItemResponse.md)
 - [PostUserVirtualCardTransactionRequest](docs/PostUserVirtualCardTransactionRequest.md)
 - [PostUserVirtualCardTransactionResponse](docs/PostUserVirtualCardTransactionResponse.md)
 - [PublicBusinessDetailsAPIModel](docs/PublicBusinessDetailsAPIModel.md)
 - [ShortBusinessDetailsAPIModel](docs/ShortBusinessDetailsAPIModel.md)
 - [ShortVirtualCardAPIModel](docs/ShortVirtualCardAPIModel.md)
 - [TransactionItemDetailAPIModel](docs/TransactionItemDetailAPIModel.md)
 - [TransactionStateEnum](docs/TransactionStateEnum.md)


## Documentation For Authorization


## sessionToken

- **Type**: HTTP basic authentication

