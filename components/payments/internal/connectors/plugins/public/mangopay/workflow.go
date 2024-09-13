package mangopay

import "github.com/formancehq/payments/internal/models"

const (
	fetchUsersName = "fetch_users"
)

func workflow() models.Tasks {
	return []models.TaskTree{
		{
			TaskType:     models.TASK_FETCH_OTHERS,
			Name:         fetchUsersName,
			Periodically: true,
			NextTasks: []models.TaskTree{
				{
					TaskType:     models.TASK_FETCH_ACCOUNTS,
					Name:         "fetch_accounts",
					Periodically: true,
					NextTasks: []models.TaskTree{
						{
							TaskType:     models.TASK_FETCH_PAYMENTS,
							Name:         "fetch_payments",
							Periodically: false, // We will be using webhooks after polling the history
							NextTasks:    []models.TaskTree{},
						},
						{
							TaskType:     models.TASK_FETCH_BALANCES,
							Name:         "fetch_balances",
							Periodically: true,
							NextTasks:    []models.TaskTree{},
						},
					},
				},
				{
					TaskType:     models.TASK_FETCH_EXTERNAL_ACCOUNTS,
					Name:         "fetch_external_accounts",
					Periodically: true,
					NextTasks:    []models.TaskTree{},
				},
			},
		},
		{
			TaskType:     models.TASK_CREATE_WEBHOOKS,
			Name:         "create_webhooks",
			Periodically: false,
			NextTasks:    []models.TaskTree{},
		},
	}
}
