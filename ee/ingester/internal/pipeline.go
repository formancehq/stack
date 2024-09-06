package ingester

import (
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/google/uuid"
)

type PipelineConfiguration struct {
	Module      string `json:"module"`
	ConnectorID string `json:"connectorID"`
}

func (p PipelineConfiguration) String() string {
	return fmt.Sprintf("%s/%s", p.Module, p.ConnectorID)
}

func NewPipelineConfiguration(module, connectorID string) PipelineConfiguration {
	return PipelineConfiguration{
		Module:      module,
		ConnectorID: connectorID,
	}
}

type Pipeline struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	State     State     `json:"state"`
	PipelineConfiguration
}

func (p Pipeline) String() string {
	return fmt.Sprintf("%s (%s): %s", p.ID, p.PipelineConfiguration, p.State)
}

func NewPipeline(pipelineConfiguration PipelineConfiguration, state State) Pipeline {
	return Pipeline{
		ID:                    uuid.NewString(),
		PipelineConfiguration: pipelineConfiguration,
		State:                 state,
		CreatedAt:             time.Now(),
	}
}
