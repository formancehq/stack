// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type WorkflowConfig struct {
	Name   *string          `json:"name,omitempty"`
	Stages []map[string]any `json:"stages"`
}

func (o *WorkflowConfig) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *WorkflowConfig) GetStages() []map[string]any {
	if o == nil {
		return []map[string]any{}
	}
	return o.Stages
}
