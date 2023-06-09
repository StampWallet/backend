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

// checks if the GetUserLocalCardTypesResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetUserLocalCardTypesResponse{}

// GetUserLocalCardTypesResponse struct for GetUserLocalCardTypesResponse
type GetUserLocalCardTypesResponse struct {
	Types []GetUserLocalCardTypesResponseTypesInner `json:"types,omitempty"`
}

// NewGetUserLocalCardTypesResponse instantiates a new GetUserLocalCardTypesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetUserLocalCardTypesResponse() *GetUserLocalCardTypesResponse {
	this := GetUserLocalCardTypesResponse{}
	return &this
}

// NewGetUserLocalCardTypesResponseWithDefaults instantiates a new GetUserLocalCardTypesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetUserLocalCardTypesResponseWithDefaults() *GetUserLocalCardTypesResponse {
	this := GetUserLocalCardTypesResponse{}
	return &this
}

// GetTypes returns the Types field value if set, zero value otherwise.
func (o *GetUserLocalCardTypesResponse) GetTypes() []GetUserLocalCardTypesResponseTypesInner {
	if o == nil || isNil(o.Types) {
		var ret []GetUserLocalCardTypesResponseTypesInner
		return ret
	}
	return o.Types
}

// GetTypesOk returns a tuple with the Types field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetUserLocalCardTypesResponse) GetTypesOk() ([]GetUserLocalCardTypesResponseTypesInner, bool) {
	if o == nil || isNil(o.Types) {
		return nil, false
	}
	return o.Types, true
}

// HasTypes returns a boolean if a field has been set.
func (o *GetUserLocalCardTypesResponse) HasTypes() bool {
	if o != nil && !isNil(o.Types) {
		return true
	}

	return false
}

// SetTypes gets a reference to the given []GetUserLocalCardTypesResponseTypesInner and assigns it to the Types field.
func (o *GetUserLocalCardTypesResponse) SetTypes(v []GetUserLocalCardTypesResponseTypesInner) {
	o.Types = v
}

func (o GetUserLocalCardTypesResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetUserLocalCardTypesResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Types) {
		toSerialize["types"] = o.Types
	}
	return toSerialize, nil
}

type NullableGetUserLocalCardTypesResponse struct {
	value *GetUserLocalCardTypesResponse
	isSet bool
}

func (v NullableGetUserLocalCardTypesResponse) Get() *GetUserLocalCardTypesResponse {
	return v.value
}

func (v *NullableGetUserLocalCardTypesResponse) Set(val *GetUserLocalCardTypesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetUserLocalCardTypesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetUserLocalCardTypesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetUserLocalCardTypesResponse(val *GetUserLocalCardTypesResponse) *NullableGetUserLocalCardTypesResponse {
	return &NullableGetUserLocalCardTypesResponse{value: val, isSet: true}
}

func (v NullableGetUserLocalCardTypesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetUserLocalCardTypesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
