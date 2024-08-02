// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type V2HookBodyParams struct {
	Endpoint string   `json:"endpoint"`
	Events   []string `json:"events"`
	Name     *string  `json:"name,omitempty"`
	Retry    *bool    `json:"retry,omitempty"`
	Secret   *string  `json:"secret,omitempty"`
}

func (o *V2HookBodyParams) GetEndpoint() string {
	if o == nil {
		return ""
	}
	return o.Endpoint
}

func (o *V2HookBodyParams) GetEvents() []string {
	if o == nil {
		return []string{}
	}
	return o.Events
}

func (o *V2HookBodyParams) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *V2HookBodyParams) GetRetry() *bool {
	if o == nil {
		return nil
	}
	return o.Retry
}

func (o *V2HookBodyParams) GetSecret() *string {
	if o == nil {
		return nil
	}
	return o.Secret
}
