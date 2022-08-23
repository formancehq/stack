package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/pkg/storage"
	"github.com/numary/webhooks/pkg/svix"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/spf13/viper"
)

type Worker struct {
	Reader Reader
	Store  storage.Store

	svixApp svix.App

	stopChan chan chan struct{}
}

func NewWorker(store storage.Store, svixApp svix.App) (*Worker, error) {
	cfg, err := NewKafkaReaderConfig()
	if err != nil {
		return nil, fmt.Errorf("kafka.NewKafkaReaderConfig: %w", err)
	}

	return &Worker{
		Reader:   kafkago.NewReader(cfg),
		Store:    store,
		svixApp:  svixApp,
		stopChan: make(chan chan struct{}),
	}, nil
}

var ErrMechanism = errors.New("unrecognized SASL mechanism")

func NewKafkaReaderConfig() (kafkago.ReaderConfig, error) {
	dialer := kafkago.DefaultDialer
	if viper.GetBool(constants.KafkaTLSEnabledFlag) {
		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if viper.GetBool(constants.KafkaSASLEnabledFlag) {
		var alg scram.Algorithm
		switch mechanism := viper.GetString(constants.KafkaSASLMechanismFlag); mechanism {
		case "SCRAM-SHA-512":
			alg = scram.SHA512
		case "SCRAM-SHA-256":
			alg = scram.SHA256
		default:
			return kafkago.ReaderConfig{}, ErrMechanism
		}
		mechanism, err := scram.Mechanism(alg,
			viper.GetString(constants.KafkaUsernameFlag), viper.GetString(constants.KafkaPasswordFlag))
		if err != nil {
			return kafkago.ReaderConfig{}, fmt.Errorf("scram.Mechanism: %w", err)
		}
		dialer.SASLMechanism = mechanism
	}

	brokers := viper.GetStringSlice(constants.KafkaBrokersFlag)
	if len(brokers) == 0 {
		brokers = []string{constants.DefaultKafkaBroker}
	}

	topics := viper.GetStringSlice(constants.KafkaTopicsFlag)
	if len(topics) == 0 {
		topics = []string{constants.DefaultKafkaTopic}
	}

	groupID := viper.GetString(constants.KafkaGroupIDFlag)
	if groupID == "" {
		groupID = constants.DefaultKafkaGroupID
	}

	return kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    1,
		MaxBytes:    10e5,
		Dialer:      dialer,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	msgChan := make(chan kafkago.Message)
	errChan := make(chan error)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go fetchMessages(ctxWithCancel, w.Reader, msgChan, errChan)

	for {
		select {
		case ch := <-w.stopChan:
			sharedlogging.GetLogger(ctx).Debug("worker: received from stopChan")
			close(ch)
			return nil
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("worker: context done: %s", ctx.Err())
			return nil
		case err := <-errChan:
			return fmt.Errorf("kafka.Worker.fetchMessages: %w", err)
		case msg := <-msgChan:
			ctx = sharedlogging.ContextWithLogger(ctx,
				sharedlogging.GetLogger(ctx).WithFields(map[string]any{
					"offset": msg.Offset,
				}))
			sharedlogging.GetLogger(ctx).WithFields(map[string]any{
				"time":      msg.Time.UTC().Format(time.RFC3339),
				"partition": msg.Partition,
				"data":      string(msg.Value),
				"headers":   msg.Headers,
			}).Debug("worker: new kafka message fetched")

			if err := processEventMessage(ctx, msg.Value, w.svixApp); err != nil {
				return fmt.Errorf("processEventMessage: %w", err)
			}

			if err := w.Reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("kafka.Reader.CommitMessages: %w", err)
			}
		}
	}
}

func fetchMessages(ctx context.Context, reader Reader, msgChan chan kafkago.Message, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if !errors.Is(err, io.EOF) && ctx.Err() == nil {
					select {
					case errChan <- fmt.Errorf("kafka.Reader.FetchMessage: %w", err):
					case <-ctx.Done():
						return
					}
				}
				continue
			}

			select {
			case msgChan <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (w *Worker) Stop(ctx context.Context) {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		sharedlogging.GetLogger(ctx).Debugf("worker stopped: context done: %s", ctx.Err())
		return
	case w.stopChan <- ch:
		select {
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("worker stopped via stopChan: context done: %s", ctx.Err())
			return
		case <-ch:
			sharedlogging.GetLogger(ctx).Debug("worker stopped via stopChan")
		}
	}
}
