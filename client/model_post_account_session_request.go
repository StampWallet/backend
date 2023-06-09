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

// checks if the PostAccountSessionRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PostAccountSessionRequest{}

// PostAccountSessionRequest struct for PostAccountSessionRequest
type PostAccountSessionRequest struct {
	Email    *string `json:"email,omitempty" binding:"required"`
	Password *string `json:"password,omitempty" binding:"required"`
}

// NewPostAccountSessionRequest instantiates a new PostAccountSessionRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPostAccountSessionRequest() *PostAccountSessionRequest {
	this := PostAccountSessionRequest{}
	return &this
}

// NewPostAccountSessionRequestWithDefaults instantiates a new PostAccountSessionRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPostAccountSessionRequestWithDefaults() *PostAccountSessionRequest {
	this := PostAccountSessionRequest{}
	return &this
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *PostAccountSessionRequest) GetEmail() string {
	if o == nil || isNil(o.Email) {
		var ret string
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostAccountSessionRequest) GetEmailOk() (*string, bool) {
	if o == nil || isNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *PostAccountSessionRequest) HasEmail() bool {
	if o != nil && !isNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given string and assigns it to the Email field.
func (o *PostAccountSessionRequest) SetEmail(v string) {
	o.Email = &v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *PostAccountSessionRequest) GetPassword() string {
	if o == nil || isNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostAccountSessionRequest) GetPasswordOk() (*string, bool) {
	if o == nil || isNil(o.Password) {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *PostAccountSessionRequest) HasPassword() bool {
	if o != nil && !isNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *PostAccountSessionRequest) SetPassword(v string) {
	o.Password = &v
}

func (o PostAccountSessionRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PostAccountSessionRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Email) {
		toSerialize["email"] = o.Email
	}
	if !isNil(o.Password) {
		toSerialize["password"] = o.Password
	}
	return toSerialize, nil
}

type NullablePostAccountSessionRequest struct {
	value *PostAccountSessionRequest
	isSet bool
}

func (v NullablePostAccountSessionRequest) Get() *PostAccountSessionRequest {
	return v.value
}

func (v *NullablePostAccountSessionRequest) Set(val *PostAccountSessionRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePostAccountSessionRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePostAccountSessionRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePostAccountSessionRequest(val *PostAccountSessionRequest) *NullablePostAccountSessionRequest {
	return &NullablePostAccountSessionRequest{value: val, isSet: true}
}

func (v NullablePostAccountSessionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePostAccountSessionRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
