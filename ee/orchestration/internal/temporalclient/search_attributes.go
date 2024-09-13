package temporalclient

import (
	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"
	"go.temporal.io/api/enums/v1"
)

var (
	SearchAttributes = map[string]enums.IndexedValueType{
		workflow.SearchAttributeWorkflowID: enums.INDEXED_VALUE_TYPE_TEXT,
		triggers.SearchAttributeTriggerID:  enums.INDEXED_VALUE_TYPE_TEXT,
	}
)
