// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
	"time"
)

type OrchestrationV2Transaction struct {
	Metadata  map[string]string `json:"metadata"`
	Postings  []V2Posting       `json:"postings"`
	Reference *string           `json:"reference,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Txid      *big.Int          `json:"txid"`
}

func (o OrchestrationV2Transaction) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(o, "", false)
}

func (o *OrchestrationV2Transaction) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &o, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *OrchestrationV2Transaction) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *OrchestrationV2Transaction) GetPostings() []V2Posting {
	if o == nil {
		return []V2Posting{}
	}
	return o.Postings
}

func (o *OrchestrationV2Transaction) GetReference() *string {
	if o == nil {
		return nil
	}
	return o.Reference
}

func (o *OrchestrationV2Transaction) GetTimestamp() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.Timestamp
}

func (o *OrchestrationV2Transaction) GetTxid() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Txid
}
