//+build it

package modules

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"text/template"
	"time"

	"github.com/formancehq/stack/ee/ingester/internal/httpclient"

	"github.com/ThreeDotsLabs/watermill"
	wNats "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func startNatsServer(t *testing.T) *server.Server {
	t.Helper()

	natsDir := t.TempDir()
	natsServer, err := server.NewServer(&server.Options{
		Host:         "0.0.0.0",
		Port:         -1,
		JetStream:    true,
		StoreDir:     natsDir,
		Debug:        testing.Verbose(),
	})
	require.NoError(t, err)

	natsServer.ConfigureLogger()
	natsServer.Start()
	t.Cleanup(natsServer.Shutdown)

	return natsServer
}

func TestModule(t *testing.T) {
	t.Parallel()

	// Setup nats
	natsServer := startNatsServer(t)

	conn, err := nats.Connect(natsServer.ClientURL())
	require.NoError(t, err)

	js, err := conn.JetStream()
	require.NoError(t, err)

	stackName := uuid.NewString()[0:8]

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     stackName,
		Subjects: []string{stackName + ".*"},
		Retention: nats.InterestPolicy,
	})
	require.NoError(t, err)

	consumerConfig := &nats.ConsumerConfig{
		Durable:            stackName,
		Name:               stackName,
		Description:        fmt.Sprintf("test %s", stackName),
		FilterSubject:      fmt.Sprintf("%s.module1", stackName),
		DeliverSubject: "test",
	}
	_, err = js.AddConsumer(stackName, consumerConfig)
	require.NoError(t, err)

	// Create a new subscriber
	subscriber, err := wNats.NewSubscriber(wNats.SubscriberConfig{
		URL: natsServer.ClientURL(),
		QueueGroupPrefix: consumerConfig.DeliverGroup,
		JetStream: wNats.JetStreamConfig{
			DurablePrefix: consumerConfig.Durable,
		},
	}, watermill.NewStdLogger(testing.Verbose(), testing.Verbose()))
	require.NoError(t, err)

	moduleFactory := NewModuleFactory(
		subscriber,
		stackName,
		httpclient.NewStackAuthenticatedClientFromHTTPClient(http.DefaultClient),
		PullConfiguration{
			ModuleURLTpl: template.Must(template.New("").Parse("http://localhost")),
			PullPageSize: 100,
		},
		logging.Testing(),
	)
	module := moduleFactory.Create("module1")

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Subscribe
	subscription, err := module.Subscribe(ctx)
	require.NoError(t, err)

	// Publish a new message
	_, err = js.Publish(stackName+".module1", []byte("hello world"))
	require.NoError(t, err)

	// Check reception
	require.Eventually(t, func() bool {
		select {
		case msg := <-subscription:
			msg.Ack()
			return true
		default:
			return false
		}
	}, time.Second, 20*time.Millisecond)

	// Close subscription
	cancel()
	select {
	case <-subscription:
		// ensure connection is closed
	}

	// Subscribe again
	ctx, cancel = context.WithCancel(context.Background())
	t.Cleanup(cancel)

	subscription, err = module.Subscribe(ctx)
	require.NoError(t, err)

	// Ensure we don't receive already received messages
	select {
	case _, ok := <-subscription:
		require.True(t, ok)
		require.Fail(t, "should not have received any messages")
	case <-time.After(2*time.Second):
	}
}