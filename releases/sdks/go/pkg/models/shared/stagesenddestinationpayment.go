// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type StageSendDestinationPayment struct {
	Psp string `json:"psp"`
}

func (o *StageSendDestinationPayment) GetPsp() string {
	if o == nil {
		return ""
	}
	return o.Psp
}
