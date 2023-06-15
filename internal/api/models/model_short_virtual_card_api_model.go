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

type ShortVirtualCardApiModel struct {
	BusinessDetails ShortBusinessDetailsApiModel `json:"businessDetails,omitempty"`

	Points int32 `json:"points"`
}
