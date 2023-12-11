// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"time"
)

type V2WorkflowInstanceHistoryStage struct {
	Attempt       int64                                 `json:"attempt"`
	Error         *string                               `json:"error,omitempty"`
	Input         V2WorkflowInstanceHistoryStageInput   `json:"input"`
	LastFailure   *string                               `json:"lastFailure,omitempty"`
	Name          string                                `json:"name"`
	NextExecution *time.Time                            `json:"nextExecution,omitempty"`
	Output        *V2WorkflowInstanceHistoryStageOutput `json:"output,omitempty"`
	StartedAt     time.Time                             `json:"startedAt"`
	Terminated    bool                                  `json:"terminated"`
	TerminatedAt  *time.Time                            `json:"terminatedAt,omitempty"`
}

func (v V2WorkflowInstanceHistoryStage) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2WorkflowInstanceHistoryStage) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2WorkflowInstanceHistoryStage) GetAttempt() int64 {
	if o == nil {
		return 0
	}
	return o.Attempt
}

func (o *V2WorkflowInstanceHistoryStage) GetError() *string {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *V2WorkflowInstanceHistoryStage) GetInput() V2WorkflowInstanceHistoryStageInput {
	if o == nil {
		return V2WorkflowInstanceHistoryStageInput{}
	}
	return o.Input
}

func (o *V2WorkflowInstanceHistoryStage) GetLastFailure() *string {
	if o == nil {
		return nil
	}
	return o.LastFailure
}

func (o *V2WorkflowInstanceHistoryStage) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *V2WorkflowInstanceHistoryStage) GetNextExecution() *time.Time {
	if o == nil {
		return nil
	}
	return o.NextExecution
}

func (o *V2WorkflowInstanceHistoryStage) GetOutput() *V2WorkflowInstanceHistoryStageOutput {
	if o == nil {
		return nil
	}
	return o.Output
}

func (o *V2WorkflowInstanceHistoryStage) GetStartedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.StartedAt
}

func (o *V2WorkflowInstanceHistoryStage) GetTerminated() bool {
	if o == nil {
		return false
	}
	return o.Terminated
}

func (o *V2WorkflowInstanceHistoryStage) GetTerminatedAt() *time.Time {
	if o == nil {
		return nil
	}
	return o.TerminatedAt
}
