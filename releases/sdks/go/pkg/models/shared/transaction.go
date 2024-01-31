// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
	"time"
)

type Transaction struct {
	Metadata          map[string]interface{}       `json:"metadata,omitempty"`
	PostCommitVolumes map[string]map[string]Volume `json:"postCommitVolumes,omitempty"`
	Postings          []Posting                    `json:"postings"`
	PreCommitVolumes  map[string]map[string]Volume `json:"preCommitVolumes,omitempty"`
	Reference         *string                      `json:"reference,omitempty"`
	Timestamp         time.Time                    `json:"timestamp"`
	Txid              *big.Int                     `json:"txid"`
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(t, "", false)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &t, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *Transaction) GetMetadata() map[string]interface{} {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *Transaction) GetPostCommitVolumes() map[string]map[string]Volume {
	if o == nil {
		return nil
	}
	return o.PostCommitVolumes
}

func (o *Transaction) GetPostings() []Posting {
	if o == nil {
		return []Posting{}
	}
	return o.Postings
}

func (o *Transaction) GetPreCommitVolumes() map[string]map[string]Volume {
	if o == nil {
		return nil
	}
	return o.PreCommitVolumes
}

func (o *Transaction) GetReference() *string {
	if o == nil {
		return nil
	}
	return o.Reference
}

func (o *Transaction) GetTimestamp() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.Timestamp
}

func (o *Transaction) GetTxid() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Txid
}