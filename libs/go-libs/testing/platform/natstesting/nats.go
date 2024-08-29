package natstesting

import (
	"fmt"
	"os"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

type NatsT interface {
	require.TestingT
	Helper()
	TempDir() string
	Cleanup(func())
}

type NatsServer struct {
	URL string

	*nats.Conn
}

func (s *NatsServer) ClientURL() string {
	return s.URL
}

func (s *NatsServer) initClient() error {
	if s.Conn != nil {
		return nil
	}
	var err error
	s.Conn, err = nats.Connect(s.URL)

	return err
}

func (s *NatsServer) Client(t NatsT) *nats.Conn {
	require.NoError(t, s.initClient())

	return s.Conn
}

func (s *NatsServer) WithJetStream(t NatsT, fn func(js nats.JetStreamContext)) {
	client := s.Client(t)

	js, err := client.JetStream()
	require.NoError(t, err)

	fn(js)
}

func (s *NatsServer) CreateConsumer(t NatsT, stream string, config *nats.ConsumerConfig) (consumer *nats.ConsumerInfo) {
	s.WithJetStream(t, func(js nats.JetStreamContext) {
		var err error
		consumer, err = js.AddConsumer(stream, config)
		require.NoError(t, err)
	})

	return consumer
}

func (s *NatsServer) CreateStream(t NatsT, name string) (stream *nats.StreamInfo) {
	s.WithJetStream(t, func(js nats.JetStreamContext) {
		var err error
		stream, err = js.AddStream(&nats.StreamConfig{
			Name:      name,
			Subjects:  []string{name + ".*"},
			Retention: nats.InterestPolicy,
		})
		require.NoError(t, err)
	})

	return stream
}

func (s *NatsServer) Publish(t NatsT, stack, module string, bytes []byte) {
	s.WithJetStream(t, func(js nats.JetStreamContext) {
		_, err := js.Publish(fmt.Sprintf("%s.%s", stack, module), bytes)
		require.NoError(t, err)
	})
}

func (s *NatsServer) ConsumerInfo(t NatsT, stream, consumerName string) (ret *nats.ConsumerInfo) {
	s.WithJetStream(t, func(js nats.JetStreamContext) {
		var err error
		ret, err = js.ConsumerInfo(stream, consumerName)
		require.NoError(t, err)
	})
	return ret
}

func CreateServer(t NatsT, debug bool, logger logging.Logger) *NatsServer {
	t.Helper()

	// Create a Nats server
	natsDir := t.TempDir()

	var err error
	natsServer, err := natsserver.NewServer(&natsserver.Options{
		Host:      "0.0.0.0",
		Port:      -1,
		JetStream: true,
		StoreDir:  natsDir,
		Debug:     debug,
	})
	require.NoError(t, err)

	if debug {
		r, w, err := os.Pipe()
		require.NoError(t, err)
		os.Stderr = w

		natsServer.ConfigureLogger()

		go logging.StreamReader(logger.WithField("service", "nats"), r, logging.Logger.Debug)
	}

	natsServer.Start()
	t.Cleanup(natsServer.Shutdown)
	require.Eventually(t, natsServer.Running, 5*time.Second, 50*time.Millisecond)

	return &NatsServer{
		URL: natsServer.ClientURL(),
	}
}
