// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type StageSendSourcePayment struct {
	ID string `json:"id"`
}

func (o *StageSendSourcePayment) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}
