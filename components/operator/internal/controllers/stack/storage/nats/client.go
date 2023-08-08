package nats

import (
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

func ExistSubject(conn *nats.Conn, subject string, logger logr.Logger) (bool, error) {

	js, _ := conn.JetStream()

	_, err := js.StreamNameBySubject(subject)
	if err != nil {
		if err.Error() == "nats: no stream matches subject" {
			return false, nil
		}
		return false, errors.Wrap(err, "cannot get stream info")
	}

	return true, nil

}

func DeleteSubject(conn *nats.Conn, subject string) error {

	js, _ := conn.JetStream()

	return js.DeleteStream(
		subject,
	)

}
