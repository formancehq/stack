// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
)

type V2BulkElementResultErrorSchemas struct {
	ErrorCode        string  `json:"errorCode"`
	ErrorDescription string  `json:"errorDescription"`
	ErrorDetails     *string `json:"errorDetails,omitempty"`
	ResponseType     string  `json:"responseType"`
}

func (o *V2BulkElementResultErrorSchemas) GetErrorCode() string {
	if o == nil {
		return ""
	}
	return o.ErrorCode
}

func (o *V2BulkElementResultErrorSchemas) GetErrorDescription() string {
	if o == nil {
		return ""
	}
	return o.ErrorDescription
}

func (o *V2BulkElementResultErrorSchemas) GetErrorDetails() *string {
	if o == nil {
		return nil
	}
	return o.ErrorDetails
}

func (o *V2BulkElementResultErrorSchemas) GetResponseType() string {
	if o == nil {
		return ""
	}
	return o.ResponseType
}

type V2BulkElementResultDeleteMetadataSchemas struct {
	ResponseType string `json:"responseType"`
}

func (o *V2BulkElementResultDeleteMetadataSchemas) GetResponseType() string {
	if o == nil {
		return ""
	}
	return o.ResponseType
}

type V2BulkElementResultRevertTransactionSchemas struct {
	Data         V2Transaction `json:"data"`
	ResponseType string        `json:"responseType"`
}

func (o *V2BulkElementResultRevertTransactionSchemas) GetData() V2Transaction {
	if o == nil {
		return V2Transaction{}
	}
	return o.Data
}

func (o *V2BulkElementResultRevertTransactionSchemas) GetResponseType() string {
	if o == nil {
		return ""
	}
	return o.ResponseType
}

type Schemas struct {
	ResponseType string `json:"responseType"`
}

func (o *Schemas) GetResponseType() string {
	if o == nil {
		return ""
	}
	return o.ResponseType
}

type V2BulkElementResultCreateTransactionSchemas struct {
	Data         V2Transaction `json:"data"`
	ResponseType string        `json:"responseType"`
}

func (o *V2BulkElementResultCreateTransactionSchemas) GetData() V2Transaction {
	if o == nil {
		return V2Transaction{}
	}
	return o.Data
}

func (o *V2BulkElementResultCreateTransactionSchemas) GetResponseType() string {
	if o == nil {
		return ""
	}
	return o.ResponseType
}

type V2BulkElementResultType string

const (
	V2BulkElementResultTypeAddMetadata       V2BulkElementResultType = "ADD_METADATA"
	V2BulkElementResultTypeCreateTransaction V2BulkElementResultType = "CREATE_TRANSACTION"
	V2BulkElementResultTypeDeleteMetadata    V2BulkElementResultType = "DELETE_METADATA"
	V2BulkElementResultTypeError             V2BulkElementResultType = "ERROR"
	V2BulkElementResultTypeRevertTransaction V2BulkElementResultType = "REVERT_TRANSACTION"
)

type V2BulkElementResult struct {
	V2BulkElementResultCreateTransactionSchemas *V2BulkElementResultCreateTransactionSchemas
	Schemas                                     *Schemas
	V2BulkElementResultRevertTransactionSchemas *V2BulkElementResultRevertTransactionSchemas
	V2BulkElementResultDeleteMetadataSchemas    *V2BulkElementResultDeleteMetadataSchemas
	V2BulkElementResultErrorSchemas             *V2BulkElementResultErrorSchemas

	Type V2BulkElementResultType
}

func CreateV2BulkElementResultAddMetadata(addMetadata Schemas) V2BulkElementResult {
	typ := V2BulkElementResultTypeAddMetadata

	typStr := string(typ)
	addMetadata.ResponseType = typStr

	return V2BulkElementResult{
		Schemas: &addMetadata,
		Type:    typ,
	}
}

func CreateV2BulkElementResultCreateTransaction(createTransaction V2BulkElementResultCreateTransactionSchemas) V2BulkElementResult {
	typ := V2BulkElementResultTypeCreateTransaction

	typStr := string(typ)
	createTransaction.ResponseType = typStr

	return V2BulkElementResult{
		V2BulkElementResultCreateTransactionSchemas: &createTransaction,
		Type: typ,
	}
}

func CreateV2BulkElementResultDeleteMetadata(deleteMetadata V2BulkElementResultDeleteMetadataSchemas) V2BulkElementResult {
	typ := V2BulkElementResultTypeDeleteMetadata

	typStr := string(typ)
	deleteMetadata.ResponseType = typStr

	return V2BulkElementResult{
		V2BulkElementResultDeleteMetadataSchemas: &deleteMetadata,
		Type:                                     typ,
	}
}

func CreateV2BulkElementResultError(error V2BulkElementResultErrorSchemas) V2BulkElementResult {
	typ := V2BulkElementResultTypeError

	typStr := string(typ)
	error.ResponseType = typStr

	return V2BulkElementResult{
		V2BulkElementResultErrorSchemas: &error,
		Type:                            typ,
	}
}

func CreateV2BulkElementResultRevertTransaction(revertTransaction V2BulkElementResultRevertTransactionSchemas) V2BulkElementResult {
	typ := V2BulkElementResultTypeRevertTransaction

	typStr := string(typ)
	revertTransaction.ResponseType = typStr

	return V2BulkElementResult{
		V2BulkElementResultRevertTransactionSchemas: &revertTransaction,
		Type: typ,
	}
}

func (u *V2BulkElementResult) UnmarshalJSON(data []byte) error {

	type discriminator struct {
		ResponseType string `json:"responseType"`
	}

	dis := new(discriminator)
	if err := json.Unmarshal(data, &dis); err != nil {
		return fmt.Errorf("could not unmarshal discriminator: %w", err)
	}

	switch dis.ResponseType {
	case "ADD_METADATA":
		schemas := new(Schemas)
		if err := utils.UnmarshalJSON(data, &schemas, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.Schemas = schemas
		u.Type = V2BulkElementResultTypeAddMetadata
		return nil
	case "CREATE_TRANSACTION":
		v2BulkElementResultCreateTransactionSchemas := new(V2BulkElementResultCreateTransactionSchemas)
		if err := utils.UnmarshalJSON(data, &v2BulkElementResultCreateTransactionSchemas, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.V2BulkElementResultCreateTransactionSchemas = v2BulkElementResultCreateTransactionSchemas
		u.Type = V2BulkElementResultTypeCreateTransaction
		return nil
	case "DELETE_METADATA":
		v2BulkElementResultDeleteMetadataSchemas := new(V2BulkElementResultDeleteMetadataSchemas)
		if err := utils.UnmarshalJSON(data, &v2BulkElementResultDeleteMetadataSchemas, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.V2BulkElementResultDeleteMetadataSchemas = v2BulkElementResultDeleteMetadataSchemas
		u.Type = V2BulkElementResultTypeDeleteMetadata
		return nil
	case "ERROR":
		v2BulkElementResultErrorSchemas := new(V2BulkElementResultErrorSchemas)
		if err := utils.UnmarshalJSON(data, &v2BulkElementResultErrorSchemas, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.V2BulkElementResultErrorSchemas = v2BulkElementResultErrorSchemas
		u.Type = V2BulkElementResultTypeError
		return nil
	case "REVERT_TRANSACTION":
		v2BulkElementResultRevertTransactionSchemas := new(V2BulkElementResultRevertTransactionSchemas)
		if err := utils.UnmarshalJSON(data, &v2BulkElementResultRevertTransactionSchemas, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.V2BulkElementResultRevertTransactionSchemas = v2BulkElementResultRevertTransactionSchemas
		u.Type = V2BulkElementResultTypeRevertTransaction
		return nil
	}

	return errors.New("could not unmarshal into supported union types")
}

func (u V2BulkElementResult) MarshalJSON() ([]byte, error) {
	if u.V2BulkElementResultCreateTransactionSchemas != nil {
		return utils.MarshalJSON(u.V2BulkElementResultCreateTransactionSchemas, "", true)
	}

	if u.Schemas != nil {
		return utils.MarshalJSON(u.Schemas, "", true)
	}

	if u.V2BulkElementResultRevertTransactionSchemas != nil {
		return utils.MarshalJSON(u.V2BulkElementResultRevertTransactionSchemas, "", true)
	}

	if u.V2BulkElementResultDeleteMetadataSchemas != nil {
		return utils.MarshalJSON(u.V2BulkElementResultDeleteMetadataSchemas, "", true)
	}

	if u.V2BulkElementResultErrorSchemas != nil {
		return utils.MarshalJSON(u.V2BulkElementResultErrorSchemas, "", true)
	}

	return nil, errors.New("could not marshal union type: all fields are null")
}
