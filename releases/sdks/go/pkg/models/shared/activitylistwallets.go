// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type ActivityListWallets struct {
	Name *string `json:"name,omitempty"`
}

func (o *ActivityListWallets) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}
