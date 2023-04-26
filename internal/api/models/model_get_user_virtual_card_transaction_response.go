/*
 * StampWallet API Server
 *
 * StampWallet API Server REST Specification
 *
 * API version: 0.1.0
 * Contact: fbstachura@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

type GetUserVirtualCardTransactionResponse struct {

	ItemId string `json:"itemId,omitempty"`

	State TransactionStateEnum `json:"state,omitempty"`

	AddedPoints int32 `json:"addedPoints,omitempty"`

	ItemActions []ItemActionApiModel `json:"itemActions,omitempty"`
}