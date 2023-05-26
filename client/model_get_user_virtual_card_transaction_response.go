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

// checks if the GetUserVirtualCardTransactionResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetUserVirtualCardTransactionResponse{}

// GetUserVirtualCardTransactionResponse struct for GetUserVirtualCardTransactionResponse
type GetUserVirtualCardTransactionResponse struct {
	ItemId      *string               `json:"itemId,omitempty"`
	State       *TransactionStateEnum `json:"state,omitempty"`
	AddedPoints *int32                `json:"addedPoints,omitempty"`
	ItemActions []ItemActionAPIModel  `json:"itemActions,omitempty"`
}

// NewGetUserVirtualCardTransactionResponse instantiates a new GetUserVirtualCardTransactionResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetUserVirtualCardTransactionResponse() *GetUserVirtualCardTransactionResponse {
	this := GetUserVirtualCardTransactionResponse{}
	return &this
}

// NewGetUserVirtualCardTransactionResponseWithDefaults instantiates a new GetUserVirtualCardTransactionResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetUserVirtualCardTransactionResponseWithDefaults() *GetUserVirtualCardTransactionResponse {
	this := GetUserVirtualCardTransactionResponse{}
	return &this
}

// GetItemId returns the ItemId field value if set, zero value otherwise.
func (o *GetUserVirtualCardTransactionResponse) GetItemId() string {
	if o == nil || isNil(o.ItemId) {
		var ret string
		return ret
	}
	return *o.ItemId
}

// GetItemIdOk returns a tuple with the ItemId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetUserVirtualCardTransactionResponse) GetItemIdOk() (*string, bool) {
	if o == nil || isNil(o.ItemId) {
		return nil, false
	}
	return o.ItemId, true
}

// HasItemId returns a boolean if a field has been set.
func (o *GetUserVirtualCardTransactionResponse) HasItemId() bool {
	if o != nil && !isNil(o.ItemId) {
		return true
	}

	return false
}

// SetItemId gets a reference to the given string and assigns it to the ItemId field.
func (o *GetUserVirtualCardTransactionResponse) SetItemId(v string) {
	o.ItemId = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *GetUserVirtualCardTransactionResponse) GetState() TransactionStateEnum {
	if o == nil || isNil(o.State) {
		var ret TransactionStateEnum
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetUserVirtualCardTransactionResponse) GetStateOk() (*TransactionStateEnum, bool) {
	if o == nil || isNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *GetUserVirtualCardTransactionResponse) HasState() bool {
	if o != nil && !isNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given TransactionStateEnum and assigns it to the State field.
func (o *GetUserVirtualCardTransactionResponse) SetState(v TransactionStateEnum) {
	o.State = &v
}

// GetAddedPoints returns the AddedPoints field value if set, zero value otherwise.
func (o *GetUserVirtualCardTransactionResponse) GetAddedPoints() int32 {
	if o == nil || isNil(o.AddedPoints) {
		var ret int32
		return ret
	}
	return *o.AddedPoints
}

// GetAddedPointsOk returns a tuple with the AddedPoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetUserVirtualCardTransactionResponse) GetAddedPointsOk() (*int32, bool) {
	if o == nil || isNil(o.AddedPoints) {
		return nil, false
	}
	return o.AddedPoints, true
}

// HasAddedPoints returns a boolean if a field has been set.
func (o *GetUserVirtualCardTransactionResponse) HasAddedPoints() bool {
	if o != nil && !isNil(o.AddedPoints) {
		return true
	}

	return false
}

// SetAddedPoints gets a reference to the given int32 and assigns it to the AddedPoints field.
func (o *GetUserVirtualCardTransactionResponse) SetAddedPoints(v int32) {
	o.AddedPoints = &v
}

// GetItemActions returns the ItemActions field value if set, zero value otherwise.
func (o *GetUserVirtualCardTransactionResponse) GetItemActions() []ItemActionAPIModel {
	if o == nil || isNil(o.ItemActions) {
		var ret []ItemActionAPIModel
		return ret
	}
	return o.ItemActions
}

// GetItemActionsOk returns a tuple with the ItemActions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetUserVirtualCardTransactionResponse) GetItemActionsOk() ([]ItemActionAPIModel, bool) {
	if o == nil || isNil(o.ItemActions) {
		return nil, false
	}
	return o.ItemActions, true
}

// HasItemActions returns a boolean if a field has been set.
func (o *GetUserVirtualCardTransactionResponse) HasItemActions() bool {
	if o != nil && !isNil(o.ItemActions) {
		return true
	}

	return false
}

// SetItemActions gets a reference to the given []ItemActionAPIModel and assigns it to the ItemActions field.
func (o *GetUserVirtualCardTransactionResponse) SetItemActions(v []ItemActionAPIModel) {
	o.ItemActions = v
}

func (o GetUserVirtualCardTransactionResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetUserVirtualCardTransactionResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.ItemId) {
		toSerialize["itemId"] = o.ItemId
	}
	if !isNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !isNil(o.AddedPoints) {
		toSerialize["addedPoints"] = o.AddedPoints
	}
	if !isNil(o.ItemActions) {
		toSerialize["itemActions"] = o.ItemActions
	}
	return toSerialize, nil
}

type NullableGetUserVirtualCardTransactionResponse struct {
	value *GetUserVirtualCardTransactionResponse
	isSet bool
}

func (v NullableGetUserVirtualCardTransactionResponse) Get() *GetUserVirtualCardTransactionResponse {
	return v.value
}

func (v *NullableGetUserVirtualCardTransactionResponse) Set(val *GetUserVirtualCardTransactionResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetUserVirtualCardTransactionResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetUserVirtualCardTransactionResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetUserVirtualCardTransactionResponse(val *GetUserVirtualCardTransactionResponse) *NullableGetUserVirtualCardTransactionResponse {
	return &NullableGetUserVirtualCardTransactionResponse{value: val, isSet: true}
}

func (v NullableGetUserVirtualCardTransactionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetUserVirtualCardTransactionResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}