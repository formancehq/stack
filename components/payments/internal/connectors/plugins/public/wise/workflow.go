package wise

import "github.com/formancehq/payments/internal/models"

const (
	fetchProfileName = "fetch_profiles"
)

func workflow() models.Tasks {
	return []models.TaskTree{
		{
			TaskType:     models.TASK_FETCH_OTHERS,
			Name:         fetchProfileName,
			Periodically: true,
			NextTasks: []models.TaskTree{
				{
					TaskType:     models.TASK_FETCH_ACCOUNTS,
					Name:         "fetch_accounts",
					Periodically: true,
					NextTasks: []models.TaskTree{
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
					Name:         "fetch_recipient_accounts",
					Periodically: true,
					NextTasks:    []models.TaskTree{},
				},
				{
					TaskType:     models.TASK_CREATE_WEBHOOKS,
					Name:         "create_webhooks",
					Periodically: false,
					NextTasks:    []models.TaskTree{},
				},
				{
					TaskType:     models.TASK_FETCH_PAYMENTS,
					Name:         "fetch_payments",
					Periodically: false,
					NextTasks:    []models.TaskTree{},
				},
			},
		},
	}
}
