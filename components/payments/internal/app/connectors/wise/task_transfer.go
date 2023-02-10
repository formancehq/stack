package wise

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskTransfer(logger logging.Logger, client *client, transfer Transfer) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
	) error {
		profiles, err := client.getProfiles()
		if err != nil {
			return err
		}

		var profileID uint64

		for _, profile := range profiles {
			if fmt.Sprint(profile.ID) == transfer.Source {
				profileID = profile.ID
			}
		}

		quote, err := client.createQuote(profileID, transfer.Currency, transfer.Amount)
		if err != nil {
			return err
		}

		destinationAccount, err := strconv.ParseUint(transfer.Destination, 10, 64)
		if err != nil {
			return err
		}

		transactionID := uuid.New().String()

		err = client.createTransfer(quote, destinationAccount, transactionID)
		if err != nil {
			return err
		}

		logger.Infof("transfer created: %s", transactionID)

		return nil
	}
}
