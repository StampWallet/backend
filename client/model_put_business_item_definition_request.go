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
	"time"
)

// checks if the PutBusinessItemDefinitionRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PutBusinessItemDefinitionRequest{}

// PutBusinessItemDefinitionRequest struct for PutBusinessItemDefinitionRequest
type PutBusinessItemDefinitionRequest struct {
	Name        *string       `json:"name,omitempty"`
	Price       NullableInt32 `json:"price,omitempty"`
	Description *string       `json:"description,omitempty"`
	StartDate   NullableTime  `json:"startDate,omitempty"`
	EndDate     NullableTime  `json:"endDate,omitempty"`
	MaxAmount   NullableInt32 `json:"maxAmount,omitempty"`
	Available   *bool         `json:"available,omitempty"`
}

// NewPutBusinessItemDefinitionRequest instantiates a new PutBusinessItemDefinitionRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPutBusinessItemDefinitionRequest() *PutBusinessItemDefinitionRequest {
	this := PutBusinessItemDefinitionRequest{}
	return &this
}

// NewPutBusinessItemDefinitionRequestWithDefaults instantiates a new PutBusinessItemDefinitionRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPutBusinessItemDefinitionRequestWithDefaults() *PutBusinessItemDefinitionRequest {
	this := PutBusinessItemDefinitionRequest{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PutBusinessItemDefinitionRequest) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutBusinessItemDefinitionRequest) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PutBusinessItemDefinitionRequest) SetName(v string) {
	o.Name = &v
}

// GetPrice returns the Price field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PutBusinessItemDefinitionRequest) GetPrice() int32 {
	if o == nil || isNil(o.Price.Get()) {
		var ret int32
		return ret
	}
	return *o.Price.Get()
}

// GetPriceOk returns a tuple with the Price field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PutBusinessItemDefinitionRequest) GetPriceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Price.Get(), o.Price.IsSet()
}

// HasPrice returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasPrice() bool {
	if o != nil && o.Price.IsSet() {
		return true
	}

	return false
}

// SetPrice gets a reference to the given NullableInt32 and assigns it to the Price field.
func (o *PutBusinessItemDefinitionRequest) SetPrice(v int32) {
	o.Price.Set(&v)
}

// SetPriceNil sets the value for Price to be an explicit nil
func (o *PutBusinessItemDefinitionRequest) SetPriceNil() {
	o.Price.Set(nil)
}

// UnsetPrice ensures that no value is present for Price, not even an explicit nil
func (o *PutBusinessItemDefinitionRequest) UnsetPrice() {
	o.Price.Unset()
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *PutBusinessItemDefinitionRequest) GetDescription() string {
	if o == nil || isNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutBusinessItemDefinitionRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || isNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasDescription() bool {
	if o != nil && !isNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *PutBusinessItemDefinitionRequest) SetDescription(v string) {
	o.Description = &v
}

// GetStartDate returns the StartDate field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PutBusinessItemDefinitionRequest) GetStartDate() time.Time {
	if o == nil || isNil(o.StartDate.Get()) {
		var ret time.Time
		return ret
	}
	return *o.StartDate.Get()
}

// GetStartDateOk returns a tuple with the StartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PutBusinessItemDefinitionRequest) GetStartDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.StartDate.Get(), o.StartDate.IsSet()
}

// HasStartDate returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasStartDate() bool {
	if o != nil && o.StartDate.IsSet() {
		return true
	}

	return false
}

// SetStartDate gets a reference to the given NullableTime and assigns it to the StartDate field.
func (o *PutBusinessItemDefinitionRequest) SetStartDate(v time.Time) {
	o.StartDate.Set(&v)
}

// SetStartDateNil sets the value for StartDate to be an explicit nil
func (o *PutBusinessItemDefinitionRequest) SetStartDateNil() {
	o.StartDate.Set(nil)
}

// UnsetStartDate ensures that no value is present for StartDate, not even an explicit nil
func (o *PutBusinessItemDefinitionRequest) UnsetStartDate() {
	o.StartDate.Unset()
}

// GetEndDate returns the EndDate field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PutBusinessItemDefinitionRequest) GetEndDate() time.Time {
	if o == nil || isNil(o.EndDate.Get()) {
		var ret time.Time
		return ret
	}
	return *o.EndDate.Get()
}

// GetEndDateOk returns a tuple with the EndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PutBusinessItemDefinitionRequest) GetEndDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.EndDate.Get(), o.EndDate.IsSet()
}

// HasEndDate returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasEndDate() bool {
	if o != nil && o.EndDate.IsSet() {
		return true
	}

	return false
}

// SetEndDate gets a reference to the given NullableTime and assigns it to the EndDate field.
func (o *PutBusinessItemDefinitionRequest) SetEndDate(v time.Time) {
	o.EndDate.Set(&v)
}

// SetEndDateNil sets the value for EndDate to be an explicit nil
func (o *PutBusinessItemDefinitionRequest) SetEndDateNil() {
	o.EndDate.Set(nil)
}

// UnsetEndDate ensures that no value is present for EndDate, not even an explicit nil
func (o *PutBusinessItemDefinitionRequest) UnsetEndDate() {
	o.EndDate.Unset()
}

// GetMaxAmount returns the MaxAmount field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PutBusinessItemDefinitionRequest) GetMaxAmount() int32 {
	if o == nil || isNil(o.MaxAmount.Get()) {
		var ret int32
		return ret
	}
	return *o.MaxAmount.Get()
}

// GetMaxAmountOk returns a tuple with the MaxAmount field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PutBusinessItemDefinitionRequest) GetMaxAmountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.MaxAmount.Get(), o.MaxAmount.IsSet()
}

// HasMaxAmount returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasMaxAmount() bool {
	if o != nil && o.MaxAmount.IsSet() {
		return true
	}

	return false
}

// SetMaxAmount gets a reference to the given NullableInt32 and assigns it to the MaxAmount field.
func (o *PutBusinessItemDefinitionRequest) SetMaxAmount(v int32) {
	o.MaxAmount.Set(&v)
}

// SetMaxAmountNil sets the value for MaxAmount to be an explicit nil
func (o *PutBusinessItemDefinitionRequest) SetMaxAmountNil() {
	o.MaxAmount.Set(nil)
}

// UnsetMaxAmount ensures that no value is present for MaxAmount, not even an explicit nil
func (o *PutBusinessItemDefinitionRequest) UnsetMaxAmount() {
	o.MaxAmount.Unset()
}

// GetAvailable returns the Available field value if set, zero value otherwise.
func (o *PutBusinessItemDefinitionRequest) GetAvailable() bool {
	if o == nil || isNil(o.Available) {
		var ret bool
		return ret
	}
	return *o.Available
}

// GetAvailableOk returns a tuple with the Available field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutBusinessItemDefinitionRequest) GetAvailableOk() (*bool, bool) {
	if o == nil || isNil(o.Available) {
		return nil, false
	}
	return o.Available, true
}

// HasAvailable returns a boolean if a field has been set.
func (o *PutBusinessItemDefinitionRequest) HasAvailable() bool {
	if o != nil && !isNil(o.Available) {
		return true
	}

	return false
}

// SetAvailable gets a reference to the given bool and assigns it to the Available field.
func (o *PutBusinessItemDefinitionRequest) SetAvailable(v bool) {
	o.Available = &v
}

func (o PutBusinessItemDefinitionRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PutBusinessItemDefinitionRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if o.Price.IsSet() {
		toSerialize["price"] = o.Price.Get()
	}
	if !isNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if o.StartDate.IsSet() {
		toSerialize["startDate"] = o.StartDate.Get()
	}
	if o.EndDate.IsSet() {
		toSerialize["endDate"] = o.EndDate.Get()
	}
	if o.MaxAmount.IsSet() {
		toSerialize["maxAmount"] = o.MaxAmount.Get()
	}
	if !isNil(o.Available) {
		toSerialize["available"] = o.Available
	}
	return toSerialize, nil
}

type NullablePutBusinessItemDefinitionRequest struct {
	value *PutBusinessItemDefinitionRequest
	isSet bool
}

func (v NullablePutBusinessItemDefinitionRequest) Get() *PutBusinessItemDefinitionRequest {
	return v.value
}

func (v *NullablePutBusinessItemDefinitionRequest) Set(val *PutBusinessItemDefinitionRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePutBusinessItemDefinitionRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePutBusinessItemDefinitionRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePutBusinessItemDefinitionRequest(val *PutBusinessItemDefinitionRequest) *NullablePutBusinessItemDefinitionRequest {
	return &NullablePutBusinessItemDefinitionRequest{value: val, isSet: true}
}

func (v NullablePutBusinessItemDefinitionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePutBusinessItemDefinitionRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}