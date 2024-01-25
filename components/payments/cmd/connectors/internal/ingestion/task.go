package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
)

func (i *DefaultIngester) UpdateTaskState(ctx context.Context, state any) error {
	taskState, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("error marshaling task state: %w", err)
	}

	if err = i.store.UpdateTaskState(ctx, i.connectorID, i.descriptor, taskState); err != nil {
		return fmt.Errorf("error updating task state: %w", err)
	}

	return nil
}
