// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

// MappingResponse - OK
type MappingResponse struct {
	Data *Mapping `json:"data,omitempty"`
}

func (o *MappingResponse) GetData() *Mapping {
	if o == nil {
		return nil
	}
	return o.Data
}
