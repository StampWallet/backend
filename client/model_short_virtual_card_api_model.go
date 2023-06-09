/*
StampWallet API Server

StampWallet API Server REST Specification

API version: 0.1.0
Contact: fbstachura@gmail.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the ShortVirtualCardAPIModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ShortVirtualCardAPIModel{}

// ShortVirtualCardAPIModel struct for ShortVirtualCardAPIModel
type ShortVirtualCardAPIModel struct {
	BusinessDetails *ShortBusinessDetailsAPIModel `json:"businessDetails,omitempty"`
	Points          int32                         `json:"points"`
}

// NewShortVirtualCardAPIModel instantiates a new ShortVirtualCardAPIModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShortVirtualCardAPIModel(points int32) *ShortVirtualCardAPIModel {
	this := ShortVirtualCardAPIModel{}
	this.Points = points
	return &this
}

// NewShortVirtualCardAPIModelWithDefaults instantiates a new ShortVirtualCardAPIModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShortVirtualCardAPIModelWithDefaults() *ShortVirtualCardAPIModel {
	this := ShortVirtualCardAPIModel{}
	return &this
}

// GetBusinessDetails returns the BusinessDetails field value if set, zero value otherwise.
func (o *ShortVirtualCardAPIModel) GetBusinessDetails() ShortBusinessDetailsAPIModel {
	if o == nil || isNil(o.BusinessDetails) {
		var ret ShortBusinessDetailsAPIModel
		return ret
	}
	return *o.BusinessDetails
}

// GetBusinessDetailsOk returns a tuple with the BusinessDetails field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ShortVirtualCardAPIModel) GetBusinessDetailsOk() (*ShortBusinessDetailsAPIModel, bool) {
	if o == nil || isNil(o.BusinessDetails) {
		return nil, false
	}
	return o.BusinessDetails, true
}

// HasBusinessDetails returns a boolean if a field has been set.
func (o *ShortVirtualCardAPIModel) HasBusinessDetails() bool {
	if o != nil && !isNil(o.BusinessDetails) {
		return true
	}

	return false
}

// SetBusinessDetails gets a reference to the given ShortBusinessDetailsAPIModel and assigns it to the BusinessDetails field.
func (o *ShortVirtualCardAPIModel) SetBusinessDetails(v ShortBusinessDetailsAPIModel) {
	o.BusinessDetails = &v
}

// GetPoints returns the Points field value
func (o *ShortVirtualCardAPIModel) GetPoints() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Points
}

// GetPointsOk returns a tuple with the Points field value
// and a boolean to check if the value has been set.
func (o *ShortVirtualCardAPIModel) GetPointsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Points, true
}

// SetPoints sets field value
func (o *ShortVirtualCardAPIModel) SetPoints(v int32) {
	o.Points = v
}

func (o ShortVirtualCardAPIModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ShortVirtualCardAPIModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.BusinessDetails) {
		toSerialize["businessDetails"] = o.BusinessDetails
	}
	toSerialize["points"] = o.Points
	return toSerialize, nil
}

type NullableShortVirtualCardAPIModel struct {
	value *ShortVirtualCardAPIModel
	isSet bool
}

func (v NullableShortVirtualCardAPIModel) Get() *ShortVirtualCardAPIModel {
	return v.value
}

func (v *NullableShortVirtualCardAPIModel) Set(val *ShortVirtualCardAPIModel) {
	v.value = val
	v.isSet = true
}

func (v NullableShortVirtualCardAPIModel) IsSet() bool {
	return v.isSet
}

func (v *NullableShortVirtualCardAPIModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableShortVirtualCardAPIModel(val *ShortVirtualCardAPIModel) *NullableShortVirtualCardAPIModel {
	return &NullableShortVirtualCardAPIModel{value: val, isSet: true}
}

func (v NullableShortVirtualCardAPIModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableShortVirtualCardAPIModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
