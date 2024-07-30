// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type CreateWorkflowRequest struct {
	Name   *string          `json:"name,omitempty"`
	Stages []map[string]any `json:"stages"`
}

func (o *CreateWorkflowRequest) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *CreateWorkflowRequest) GetStages() []map[string]any {
	if o == nil {
		return []map[string]any{}
	}
	return o.Stages
}
