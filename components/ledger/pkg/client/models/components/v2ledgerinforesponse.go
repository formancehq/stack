// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package components

type V2LedgerInfoResponse struct {
	Data *V2LedgerInfo `json:"data,omitempty"`
}

func (o *V2LedgerInfoResponse) GetData() *V2LedgerInfo {
	if o == nil {
		return nil
	}
	return o.Data
}
