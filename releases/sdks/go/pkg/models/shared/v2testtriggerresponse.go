// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type V2TestTriggerResponse struct {
	Data V2TriggerTest `json:"data"`
}

func (o *V2TestTriggerResponse) GetData() V2TriggerTest {
	if o == nil {
		return V2TriggerTest{}
	}
	return o.Data
}
