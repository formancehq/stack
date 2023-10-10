package nats

import (
	"context"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

func NewClient(natsConfig *v1beta3.NatsConfig, clientId string) (*nats.Conn, error) {
	options := []nats.Option{
		nats.Name(clientId),
	}

	conn, err := nats.Connect(natsConfig.URL, options...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to nats-core")
	}

	return conn, nil
}

func ExistSubject(ctx nats.JetStreamContext, subject string, logger logr.Logger, execContext context.Context) (bool, error) {
	_, err := ctx.StreamNameBySubject(subject, nats.Context(execContext))
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}
