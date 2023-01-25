/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
	"fmt"
)

// WorkflowInstanceHistoryStageOutput struct for WorkflowInstanceHistoryStageOutput
type WorkflowInstanceHistoryStageOutput struct {
	AccountResponse *AccountResponse
	DebitWalletResponse *DebitWalletResponse
	GetWalletResponse *GetWalletResponse
	PaymentResponse *PaymentResponse
	TransactionResponse *TransactionResponse
	TransactionsResponse *TransactionsResponse
}

// Unmarshal JSON data into any of the pointers in the struct
func (dst *WorkflowInstanceHistoryStageOutput) UnmarshalJSON(data []byte) error {
	var err error
	// try to unmarshal JSON data into AccountResponse
	err = json.Unmarshal(data, &dst.AccountResponse);
	if err == nil {
		jsonAccountResponse, _ := json.Marshal(dst.AccountResponse)
		if string(jsonAccountResponse) == "{}" { // empty struct
			dst.AccountResponse = nil
		} else {
			return nil // data stored in dst.AccountResponse, return on the first match
		}
	} else {
		dst.AccountResponse = nil
	}

	// try to unmarshal JSON data into DebitWalletResponse
	err = json.Unmarshal(data, &dst.DebitWalletResponse);
	if err == nil {
		jsonDebitWalletResponse, _ := json.Marshal(dst.DebitWalletResponse)
		if string(jsonDebitWalletResponse) == "{}" { // empty struct
			dst.DebitWalletResponse = nil
		} else {
			return nil // data stored in dst.DebitWalletResponse, return on the first match
		}
	} else {
		dst.DebitWalletResponse = nil
	}

	// try to unmarshal JSON data into GetWalletResponse
	err = json.Unmarshal(data, &dst.GetWalletResponse);
	if err == nil {
		jsonGetWalletResponse, _ := json.Marshal(dst.GetWalletResponse)
		if string(jsonGetWalletResponse) == "{}" { // empty struct
			dst.GetWalletResponse = nil
		} else {
			return nil // data stored in dst.GetWalletResponse, return on the first match
		}
	} else {
		dst.GetWalletResponse = nil
	}

	// try to unmarshal JSON data into PaymentResponse
	err = json.Unmarshal(data, &dst.PaymentResponse);
	if err == nil {
		jsonPaymentResponse, _ := json.Marshal(dst.PaymentResponse)
		if string(jsonPaymentResponse) == "{}" { // empty struct
			dst.PaymentResponse = nil
		} else {
			return nil // data stored in dst.PaymentResponse, return on the first match
		}
	} else {
		dst.PaymentResponse = nil
	}

	// try to unmarshal JSON data into TransactionResponse
	err = json.Unmarshal(data, &dst.TransactionResponse);
	if err == nil {
		jsonTransactionResponse, _ := json.Marshal(dst.TransactionResponse)
		if string(jsonTransactionResponse) == "{}" { // empty struct
			dst.TransactionResponse = nil
		} else {
			return nil // data stored in dst.TransactionResponse, return on the first match
		}
	} else {
		dst.TransactionResponse = nil
	}

	// try to unmarshal JSON data into TransactionsResponse
	err = json.Unmarshal(data, &dst.TransactionsResponse);
	if err == nil {
		jsonTransactionsResponse, _ := json.Marshal(dst.TransactionsResponse)
		if string(jsonTransactionsResponse) == "{}" { // empty struct
			dst.TransactionsResponse = nil
		} else {
			return nil // data stored in dst.TransactionsResponse, return on the first match
		}
	} else {
		dst.TransactionsResponse = nil
	}

	return fmt.Errorf("data failed to match schemas in anyOf(WorkflowInstanceHistoryStageOutput)")
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src *WorkflowInstanceHistoryStageOutput) MarshalJSON() ([]byte, error) {
	if src.AccountResponse != nil {
		return json.Marshal(&src.AccountResponse)
	}

	if src.DebitWalletResponse != nil {
		return json.Marshal(&src.DebitWalletResponse)
	}

	if src.GetWalletResponse != nil {
		return json.Marshal(&src.GetWalletResponse)
	}

	if src.PaymentResponse != nil {
		return json.Marshal(&src.PaymentResponse)
	}

	if src.TransactionResponse != nil {
		return json.Marshal(&src.TransactionResponse)
	}

	if src.TransactionsResponse != nil {
		return json.Marshal(&src.TransactionsResponse)
	}

	return nil, nil // no data in anyOf schemas
}

type NullableWorkflowInstanceHistoryStageOutput struct {
	value *WorkflowInstanceHistoryStageOutput
	isSet bool
}

func (v NullableWorkflowInstanceHistoryStageOutput) Get() *WorkflowInstanceHistoryStageOutput {
	return v.value
}

func (v *NullableWorkflowInstanceHistoryStageOutput) Set(val *WorkflowInstanceHistoryStageOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowInstanceHistoryStageOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowInstanceHistoryStageOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowInstanceHistoryStageOutput(val *WorkflowInstanceHistoryStageOutput) *NullableWorkflowInstanceHistoryStageOutput {
	return &NullableWorkflowInstanceHistoryStageOutput{value: val, isSet: true}
}

func (v NullableWorkflowInstanceHistoryStageOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowInstanceHistoryStageOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


