// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type PoolResponse struct {
	Data Pool `json:"data"`
}

func (o *PoolResponse) GetData() Pool {
	if o == nil {
		return Pool{}
	}
	return o.Data
}
