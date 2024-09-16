package adyen

import "github.com/formancehq/payments/internal/models"

func workflow() models.Tasks {
	return []models.TaskTree{
		{
			TaskType:     models.TASK_FETCH_ACCOUNTS,
			Name:         "fetch_accounts",
			Periodically: true,
			NextTasks: []models.TaskTree{
				{
					TaskType:     models.TASK_CREATE_WEBHOOKS,
					Name:         "create_webhooks",
					Periodically: false,
					NextTasks:    []models.TaskTree{},
				},
			},
		},
	}
}
