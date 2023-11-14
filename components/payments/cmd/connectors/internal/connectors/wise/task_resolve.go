package wise

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

const (
	taskNameMain                   = "main"
	taskNameFetchTransfers         = "fetch-transfers"
	taskNameFetchProfiles          = "fetch-profiles"
	taskNameFetchRecipientAccounts = "fetch-recipient-accounts"
	taskNameInitiatePayment        = "initiate-payment"
	taskNameUpdatePaymentStatus    = "update-payment-status"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string `json:"name" yaml:"name" bson:"name"`
	Key        string `json:"key" yaml:"key" bson:"key"`
	ProfileID  uint64 `json:"profileID" yaml:"profileID" bson:"profileID"`
	TransferID string `json:"transferID" yaml:"transferID" bson:"transferID"`
	PaymentID  string `json:"paymentID" yaml:"paymentID" bson:"paymentID"`
	Attempt    int    `json:"attempt" yaml:"attempt" bson:"attempt"`
}

func (c *Connector) resolveTasks() func(taskDefinition TaskDescriptor) task.Task {
	client := client.NewClient(c.cfg.APIKey)

	return func(taskDefinition TaskDescriptor) task.Task {
		switch taskDefinition.Key {
		case taskNameMain:
			return taskMain()
		case taskNameFetchProfiles:
			return taskFetchProfiles(client)
		case taskNameFetchRecipientAccounts:
			return taskFetchRecipientAccounts(client, taskDefinition.ProfileID)
		case taskNameFetchTransfers:
			return taskFetchTransfers(client, taskDefinition.ProfileID)
		case taskNameInitiatePayment:
			return taskInitiatePayment(client, taskDefinition.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(client, taskDefinition.TransferID, taskDefinition.PaymentID, taskDefinition.Attempt)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDefinition.Key, ErrMissingTask)
		}
	}
}
