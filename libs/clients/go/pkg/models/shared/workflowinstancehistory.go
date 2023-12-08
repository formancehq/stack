// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"time"
)

type WorkflowInstanceHistory struct {
	Error        *string    `json:"error,omitempty"`
	Input        Stage      `json:"input"`
	Name         string     `json:"name"`
	StartedAt    time.Time  `json:"startedAt"`
	Terminated   bool       `json:"terminated"`
	TerminatedAt *time.Time `json:"terminatedAt,omitempty"`
}

func (w WorkflowInstanceHistory) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(w, "", false)
}

func (w *WorkflowInstanceHistory) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &w, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *WorkflowInstanceHistory) GetError() *string {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *WorkflowInstanceHistory) GetInput() Stage {
	if o == nil {
		return Stage{}
	}
	return o.Input
}

func (o *WorkflowInstanceHistory) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *WorkflowInstanceHistory) GetStartedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.StartedAt
}

func (o *WorkflowInstanceHistory) GetTerminated() bool {
	if o == nil {
		return false
	}
	return o.Terminated
}

func (o *WorkflowInstanceHistory) GetTerminatedAt() *time.Time {
	if o == nil {
		return nil
	}
	return o.TerminatedAt
}
