package currencycloud

import "github.com/formancehq/payments/internal/models"

const (
	fetchContactID = "fetch_contact_id"
)

func workflow() models.Tasks {
	return []models.TaskTree{
		{
			TaskType:     models.TASK_FETCH_ACCOUNTS,
			Name:         "fetch_accounts",
			Periodically: true,
			NextTasks: []models.TaskTree{
				{
					TaskType:     models.TASK_FETCH_PAYMENTS,
					Name:         "fetch_payments",
					Periodically: true,
					NextTasks:    []models.TaskTree{},
				},
			},
		},
		{
			TaskType:     models.TASK_FETCH_EXTERNAL_ACCOUNTS,
			Name:         "fetch_beneficiaries",
			Periodically: true,
			NextTasks:    []models.TaskTree{},
		},
		{

			TaskType:     models.TASK_FETCH_BALANCES,
			Name:         "fetch_balances",
			Periodically: true,
			NextTasks:    []models.TaskTree{},
		},
	}
}
