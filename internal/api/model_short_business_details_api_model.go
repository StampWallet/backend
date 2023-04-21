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

type ShortBusinessDetailsApiModel struct {

	PublicId string `json:"publicId,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	GpsCoordinates string `json:"gpsCoordinates,omitempty"`

	BannerImageId string `json:"bannerImageId,omitempty"`

	IconImageId string `json:"iconImageId,omitempty"`
}
