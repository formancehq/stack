package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/httpserver"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func StartModule(addr, serviceName string, retriesCron time.Duration, retriesSchedule []time.Duration) fx.Option {
	var options []fx.Option

	options = append(options, otlptraces.CLITracesModule(viper.GetViper()))
	options = append(options, publish.CLIPublisherModule(viper.GetViper(), serviceName))

	options = append(options, fx.Provide(
		func() (string, time.Duration, []time.Duration) {
			return addr, retriesCron, retriesSchedule
		},
		httpserver.NewMuxServer,
		postgres.NewStore,
		NewRetrier,
		newWorkerHandler,
	))
	options = append(options, fx.Invoke(httpserver.RegisterHandler))
	options = append(options, fx.Invoke(httpserver.Run))
	options = append(options, fx.Invoke(run))
	options = append(options, fx.Invoke(func(r *message.Router, subscriber message.Subscriber, store storage.Store, httpClient *http.Client) {
		configureMessageRouter(r, subscriber, viper.GetStringSlice(flag.KafkaTopics), store, httpClient, retriesSchedule)
	}))

	logging.Debugf("starting worker with env:")
	for _, e := range os.Environ() {
		logging.Debugf("%s", e)
	}

	return fx.Module("webhooks worker", options...)
}

func run(lc fx.Lifecycle, w *Retrier) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Debugf("starting worker...")
			go func() {
				if err := w.Run(context.Background()); err != nil {
					logging.FromContext(ctx).Errorf("kafka.Retrier.Run: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Debugf("stopping worker...")
			w.Stop(ctx)

			if err := w.store.Close(ctx); err != nil {
				return fmt.Errorf("storage.Store.Close: %w", err)
			}
			return nil
		},
	})
}

func configureMessageRouter(r *message.Router, subscriber message.Subscriber, topics []string,
	store storage.Store, httpClient *http.Client, retriesSchedule []time.Duration,
) {
	for _, topic := range topics {
		r.AddNoPublisherHandler(fmt.Sprintf("messages-%s", topic), topic, subscriber, processMessages(store, httpClient, retriesSchedule))
	}
}

const (
	EventTypeCommittedTransactions = "committed_transactions"
	EventTypeSavedMetadata         = "saved_metadata"
	EventTypeRevertedTransaction   = "reverted_transaction"
	EventTypeSavedPayments         = "saved_payment"
	EventTypeSavedAccounts         = "saved_account"
	EventTypeConnectorReset        = "connector_reset"
)

func processMessages(store storage.Store, httpClient *http.Client, retriesSchedule []time.Duration) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		var event events.Event
		if err := proto.Unmarshal(msg.Payload, &event); err != nil {
			return errors.Wrap(err, "proto.Unmarshal event message")
		}

		eventApp := strings.ToLower(event.App)
		var eventType string
		switch event.Event.(type) {
		case *events.Event_AccountSaved:
			eventType = EventTypeCommittedTransactions
		case *events.Event_MetadataSaved:
			eventType = EventTypeSavedMetadata
		case *events.Event_PaymentSaved:
			eventType = EventTypeSavedPayments
		case *events.Event_ResetConnector:
			eventType = EventTypeConnectorReset
		case *events.Event_TransactionReverted:
			eventType = EventTypeRevertedTransaction
		case *events.Event_TransactionsCommitted:
			eventType = EventTypeCommittedTransactions
		default:
			return errors.New("unknown event type")
		}

		evType := ""
		if eventApp == "" {
			evType = eventType
		} else {
			evType = strings.Join([]string{eventApp, eventType}, ".")
		}

		filter := map[string]any{
			"event_types": evType,
			"active":      true,
		}
		logging.FromContext(msg.Context()).Debugf("searching configs with filter: %+v", filter)
		cfgs, err := store.FindManyConfigs(msg.Context(), filter)
		if err != nil {
			return errors.Wrap(err, "storage.store.FindManyConfigs")
		}

		response := webhooks.BuildResponseFromMessage(&event)
		for _, cfg := range cfgs {
			logging.FromContext(msg.Context()).Debugf("found one config: %+v", cfg)
			data, err := json.Marshal(response)
			if err != nil {
				return errors.Wrap(err, "json.Marshal event message")
			}

			attempt, err := webhooks.MakeAttempt(msg.Context(), httpClient, retriesSchedule, uuid.NewString(),
				uuid.NewString(), 0, cfg, data, false)
			if err != nil {
				return errors.Wrap(err, "sending webhook")
			}

			if attempt.Status == webhooks.StatusAttemptSuccess {
				logging.FromContext(msg.Context()).Infof(
					"webhook sent with ID %s to %s of type %s",
					attempt.WebhookID, cfg.Endpoint, eventType)
			}

			if err := store.InsertOneAttempt(msg.Context(), attempt); err != nil {
				return errors.Wrap(err, "storage.store.InsertOneAttempt")
			}
		}
		return nil
	}
}
