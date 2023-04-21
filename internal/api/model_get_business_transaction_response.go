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

type GetBusinessTransactionResponse struct {

	PublicId string `json:"publicId,omitempty"`

	VirtualCardId int32 `json:"virtualCardId,omitempty"`

	State TransactionStateEnum `json:"state,omitempty"`

	Items []TransactionItemDetailApiModel `json:"items,omitempty"`
}
