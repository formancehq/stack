package wise

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchProfiles(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
	) error {
		profiles, err := client.GetProfiles()
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			logger.Infof(fmt.Sprintf("scheduling fetch-transfers: %d", profile.ID))

			descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transfers from client by profile",
				Key:       taskNameFetchTransfers,
				ProfileID: profile.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, descriptor, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				Restart:        true,
			})
			if err != nil {
				return err
			}
		}

		return nil
	}
}
