// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type ReadTriggerResponse struct {
	Data Trigger `json:"data"`
}

func (o *ReadTriggerResponse) GetData() Trigger {
	if o == nil {
		return Trigger{}
	}
	return o.Data
}
