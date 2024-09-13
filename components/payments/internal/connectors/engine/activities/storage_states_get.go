package activities

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStatesGet(ctx context.Context, id models.StateID) (*models.State, error) {
	resp, err := a.storage.StatesGet(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &models.State{
				ID:          id,
				ConnectorID: id.ConnectorID,
				State:       nil,
			}, nil
		}
	}
	return &resp, nil
}

var StorageStatesGetActivity = Activities{}.StorageStatesGet

func StorageStatesGet(ctx workflow.Context, id models.StateID) (*models.State, error) {
	ret := models.State{}
	if err := executeActivity(ctx, StorageStatesGetActivity, &ret, id); err != nil {
		return nil, err
	}
	return &ret, nil
}
