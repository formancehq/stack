package controllers

import (
	"os"
	"testing"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

func TestMain(m *testing.M) {
	server, err := natsserver.NewServer(&natsserver.Options{
		Host:      "0.0.0.0",
		Port:      4322,
		JetStream: true,
	})
	if err != nil {
		panic(err)
	}

	server.Start()
	defer server.Shutdown()

	code := m.Run()
	os.Exit(code)
}
