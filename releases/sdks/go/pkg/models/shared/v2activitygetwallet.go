// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type V2ActivityGetWallet struct {
	ID string `json:"id"`
}

func (o *V2ActivityGetWallet) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}
