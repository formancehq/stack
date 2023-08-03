package nats

import "github.com/nats-io/nats.go"

func DeleteSubject(conn *nats.Conn, subject string) error {

	js, _ := conn.JetStream()

	return js.DeleteStream(
		subject,
	)

}
