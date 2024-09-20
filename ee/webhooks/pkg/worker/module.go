package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/formancehq/go-libs/contextutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alitto/pond"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/publish"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

var Tracer = otel.Tracer("listener")

func StartModule(cmd *cobra.Command, retriesCron time.Duration, retryPolicy webhooks.BackoffPolicy, debug bool, topics []string) fx.Option {
	var options []fx.Option

	options = append(options, fx.Invoke(func(r *message.Router, subscriber message.Subscriber, store storage.Store, httpClient *http.Client) {
		configureMessageRouter(r, subscriber, topics, store, httpClient, retryPolicy, pond.New(50, 50))
	}))
	options = append(options, publish.FXModuleFromFlags(cmd, debug))
	options = append(options, fx.Provide(
		func() (time.Duration, webhooks.BackoffPolicy) {
			return retriesCron, retryPolicy
		},
		NewRetrier,
	))
	options = append(options, fx.Invoke(run))

	logging.Debugf("starting worker with env:")
	for _, e := range os.Environ() {
		logging.Debugf("%s", e)
	}

	return fx.Options(options...)
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
	store storage.Store, httpClient *http.Client, retryPolicy webhooks.BackoffPolicy, pool *pond.WorkerPool,
) {
	for _, topic := range topics {
		r.AddNoPublisherHandler(fmt.Sprintf("messages-%s", topic), topic, subscriber, processMessages(store, httpClient, retryPolicy, pool))
	}
}

func processMessages(store storage.Store, httpClient *http.Client, retryPolicy webhooks.BackoffPolicy, pool *pond.WorkerPool) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		pool.Submit(func() {

			var ev *publish.EventMessage
			span, ev, err := publish.UnmarshalMessage(msg)
			if err != nil {
				logging.FromContext(msg.Context()).Error(err.Error())
				return
			}

			ctx, span := Tracer.Start(msg.Context(), "HandleEvent",
				trace.WithLinks(trace.Link{
					SpanContext: span.SpanContext(),
				}),
				trace.WithAttributes(
					attribute.String("event-id", msg.UUID),
					attribute.Bool("duplicate", false),
					attribute.String("event-type", ev.Type),
					attribute.String("event-payload", string(msg.Payload)),
				),
			)
			defer span.End()
			defer func() {
				if err != nil {
					span.RecordError(err)
				}
			}()
			ctx, _ = contextutil.Detached(ctx)

			eventApp := strings.ToLower(ev.App)
			eventType := strings.ToLower(ev.Type)

			if eventApp == "" {
				ev.Type = eventType
			} else {
				ev.Type = strings.Join([]string{eventApp, eventType}, ".")
			}

			filter := map[string]any{
				"event_types": ev.Type,
				"active":      true,
			}
			logging.FromContext(ctx).Debugf("searching configs with event types: %+v", ev.Type)
			cfgs, err := store.FindManyConfigs(ctx, filter)
			if err != nil {
				logging.FromContext(ctx).Error(err)
				return
			}

			for _, cfg := range cfgs {
				logging.FromContext(ctx).Debugf("found one config: %+v", cfg)
				data, err := json.Marshal(ev)
				if err != nil {
					logging.FromContext(ctx).Error(err)
					return
				}

				attempt, err := webhooks.MakeAttempt(ctx, httpClient, retryPolicy, uuid.NewString(),
					uuid.NewString(), 0, cfg, data, false)
				if err != nil {
					logging.FromContext(ctx).Error(err)
					return
				}

				if attempt.Status == webhooks.StatusAttemptSuccess {
					logging.FromContext(ctx).Debugf(
						"webhook sent with ID %s to %s of type %s",
						attempt.WebhookID, cfg.Endpoint, ev.Type)
				}

				if err := store.InsertOneAttempt(ctx, attempt); err != nil {
					logging.FromContext(ctx).Error(err)
					return
				}
			}
		})
		return nil
	}
}
