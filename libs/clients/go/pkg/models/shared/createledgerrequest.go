// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type CreateLedgerRequest struct {
	Bucket *string `json:"bucket,omitempty"`
}

func (o *CreateLedgerRequest) GetBucket() *string {
	if o == nil {
		return nil
	}
	return o.Bucket
}
