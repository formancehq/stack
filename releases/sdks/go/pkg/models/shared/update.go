// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type Update struct {
	Account *UpdateAccount `json:"account,omitempty"`
}

func (o *Update) GetAccount() *UpdateAccount {
	if o == nil {
		return nil
	}
	return o.Account
}
