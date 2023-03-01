/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230301
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
	"fmt"
)

// TasksCursorCursorAllOfDataInner - struct for TasksCursorCursorAllOfDataInner
type TasksCursorCursorAllOfDataInner struct {
	TaskBankingCircle *TaskBankingCircle
	TaskCurrencyCloud *TaskCurrencyCloud
	TaskDummyPay *TaskDummyPay
	TaskModulr *TaskModulr
	TaskStripe *TaskStripe
	TaskWise *TaskWise
}

// TaskBankingCircleAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskBankingCircle wrapped in TasksCursorCursorAllOfDataInner
func TaskBankingCircleAsTasksCursorCursorAllOfDataInner(v *TaskBankingCircle) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskBankingCircle: v,
	}
}

// TaskCurrencyCloudAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskCurrencyCloud wrapped in TasksCursorCursorAllOfDataInner
func TaskCurrencyCloudAsTasksCursorCursorAllOfDataInner(v *TaskCurrencyCloud) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskCurrencyCloud: v,
	}
}

// TaskDummyPayAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskDummyPay wrapped in TasksCursorCursorAllOfDataInner
func TaskDummyPayAsTasksCursorCursorAllOfDataInner(v *TaskDummyPay) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskDummyPay: v,
	}
}

// TaskModulrAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskModulr wrapped in TasksCursorCursorAllOfDataInner
func TaskModulrAsTasksCursorCursorAllOfDataInner(v *TaskModulr) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskModulr: v,
	}
}

// TaskStripeAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskStripe wrapped in TasksCursorCursorAllOfDataInner
func TaskStripeAsTasksCursorCursorAllOfDataInner(v *TaskStripe) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskStripe: v,
	}
}

// TaskWiseAsTasksCursorCursorAllOfDataInner is a convenience function that returns TaskWise wrapped in TasksCursorCursorAllOfDataInner
func TaskWiseAsTasksCursorCursorAllOfDataInner(v *TaskWise) TasksCursorCursorAllOfDataInner {
	return TasksCursorCursorAllOfDataInner{
		TaskWise: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *TasksCursorCursorAllOfDataInner) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into TaskBankingCircle
	err = newStrictDecoder(data).Decode(&dst.TaskBankingCircle)
	if err == nil {
		jsonTaskBankingCircle, _ := json.Marshal(dst.TaskBankingCircle)
		if string(jsonTaskBankingCircle) == "{}" { // empty struct
			dst.TaskBankingCircle = nil
		} else {
			match++
		}
	} else {
		dst.TaskBankingCircle = nil
	}

	// try to unmarshal data into TaskCurrencyCloud
	err = newStrictDecoder(data).Decode(&dst.TaskCurrencyCloud)
	if err == nil {
		jsonTaskCurrencyCloud, _ := json.Marshal(dst.TaskCurrencyCloud)
		if string(jsonTaskCurrencyCloud) == "{}" { // empty struct
			dst.TaskCurrencyCloud = nil
		} else {
			match++
		}
	} else {
		dst.TaskCurrencyCloud = nil
	}

	// try to unmarshal data into TaskDummyPay
	err = newStrictDecoder(data).Decode(&dst.TaskDummyPay)
	if err == nil {
		jsonTaskDummyPay, _ := json.Marshal(dst.TaskDummyPay)
		if string(jsonTaskDummyPay) == "{}" { // empty struct
			dst.TaskDummyPay = nil
		} else {
			match++
		}
	} else {
		dst.TaskDummyPay = nil
	}

	// try to unmarshal data into TaskModulr
	err = newStrictDecoder(data).Decode(&dst.TaskModulr)
	if err == nil {
		jsonTaskModulr, _ := json.Marshal(dst.TaskModulr)
		if string(jsonTaskModulr) == "{}" { // empty struct
			dst.TaskModulr = nil
		} else {
			match++
		}
	} else {
		dst.TaskModulr = nil
	}

	// try to unmarshal data into TaskStripe
	err = newStrictDecoder(data).Decode(&dst.TaskStripe)
	if err == nil {
		jsonTaskStripe, _ := json.Marshal(dst.TaskStripe)
		if string(jsonTaskStripe) == "{}" { // empty struct
			dst.TaskStripe = nil
		} else {
			match++
		}
	} else {
		dst.TaskStripe = nil
	}

	// try to unmarshal data into TaskWise
	err = newStrictDecoder(data).Decode(&dst.TaskWise)
	if err == nil {
		jsonTaskWise, _ := json.Marshal(dst.TaskWise)
		if string(jsonTaskWise) == "{}" { // empty struct
			dst.TaskWise = nil
		} else {
			match++
		}
	} else {
		dst.TaskWise = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.TaskBankingCircle = nil
		dst.TaskCurrencyCloud = nil
		dst.TaskDummyPay = nil
		dst.TaskModulr = nil
		dst.TaskStripe = nil
		dst.TaskWise = nil

		return fmt.Errorf("data matches more than one schema in oneOf(TasksCursorCursorAllOfDataInner)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(TasksCursorCursorAllOfDataInner)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src TasksCursorCursorAllOfDataInner) MarshalJSON() ([]byte, error) {
	if src.TaskBankingCircle != nil {
		return json.Marshal(&src.TaskBankingCircle)
	}

	if src.TaskCurrencyCloud != nil {
		return json.Marshal(&src.TaskCurrencyCloud)
	}

	if src.TaskDummyPay != nil {
		return json.Marshal(&src.TaskDummyPay)
	}

	if src.TaskModulr != nil {
		return json.Marshal(&src.TaskModulr)
	}

	if src.TaskStripe != nil {
		return json.Marshal(&src.TaskStripe)
	}

	if src.TaskWise != nil {
		return json.Marshal(&src.TaskWise)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *TasksCursorCursorAllOfDataInner) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.TaskBankingCircle != nil {
		return obj.TaskBankingCircle
	}

	if obj.TaskCurrencyCloud != nil {
		return obj.TaskCurrencyCloud
	}

	if obj.TaskDummyPay != nil {
		return obj.TaskDummyPay
	}

	if obj.TaskModulr != nil {
		return obj.TaskModulr
	}

	if obj.TaskStripe != nil {
		return obj.TaskStripe
	}

	if obj.TaskWise != nil {
		return obj.TaskWise
	}

	// all schemas are nil
	return nil
}

type NullableTasksCursorCursorAllOfDataInner struct {
	value *TasksCursorCursorAllOfDataInner
	isSet bool
}

func (v NullableTasksCursorCursorAllOfDataInner) Get() *TasksCursorCursorAllOfDataInner {
	return v.value
}

func (v *NullableTasksCursorCursorAllOfDataInner) Set(val *TasksCursorCursorAllOfDataInner) {
	v.value = val
	v.isSet = true
}

func (v NullableTasksCursorCursorAllOfDataInner) IsSet() bool {
	return v.isSet
}

func (v *NullableTasksCursorCursorAllOfDataInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTasksCursorCursorAllOfDataInner(val *TasksCursorCursorAllOfDataInner) *NullableTasksCursorCursorAllOfDataInner {
	return &NullableTasksCursorCursorAllOfDataInner{value: val, isSet: true}
}

func (v NullableTasksCursorCursorAllOfDataInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTasksCursorCursorAllOfDataInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


