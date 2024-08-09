package fawry

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/fawry/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
)

type Webhook struct {
}

func (w *Webhook) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := task.ConnectorContextFromContext(r.Context())

		b, err := io.ReadAll(r.Body)

		n := client.Notification{}
		if err := json.Unmarshal(b, &n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, _ := contextutil.DetachedWithTimeout(r.Context(), 30*time.Second)
		td, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:    TaskIngestPayment,
			Key:     TaskIngestPayment,
			Payload: string(b),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn.Scheduler().Schedule(ctx, td, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
	}
}

func NewWebhook() *Webhook {
	return &Webhook{}
}
