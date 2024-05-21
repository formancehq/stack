// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type Trigger struct {
	CreatedAt  time.Time      `json:"createdAt"`
	Event      string         `json:"event"`
	Filter     *string        `json:"filter,omitempty"`
	ID         string         `json:"id"`
	Name       *string        `json:"name,omitempty"`
	Vars       map[string]any `json:"vars,omitempty"`
	WorkflowID string         `json:"workflowID"`
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(t, "", false)
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &t, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *Trigger) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *Trigger) GetEvent() string {
	if o == nil {
		return ""
	}
	return o.Event
}

func (o *Trigger) GetFilter() *string {
	if o == nil {
		return nil
	}
	return o.Filter
}

func (o *Trigger) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Trigger) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *Trigger) GetVars() map[string]any {
	if o == nil {
		return nil
	}
	return o.Vars
}

func (o *Trigger) GetWorkflowID() string {
	if o == nil {
		return ""
	}
	return o.WorkflowID
}
