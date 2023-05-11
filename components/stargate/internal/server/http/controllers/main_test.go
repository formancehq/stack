package controllers

import (
	"os"
	"testing"

	natsserver "github.com/nats-io/nats-server/test"
)

func TestMain(m *testing.M) {
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	code := m.Run()
	os.Exit(code)
}
