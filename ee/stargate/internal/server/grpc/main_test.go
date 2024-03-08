package grpc_test

import (
	"context"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/generated"
	stargateserver "github.com/formancehq/stack/components/stargate/internal/server/grpc"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/metrics"
	"github.com/formancehq/stack/libs/go-libs/logging"
	natsServer "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
	nc  *nats.Conn
)

func TestMain(m *testing.M) {
	server, err := natsServer.NewServer(&natsServer.Options{
		Host:      "0.0.0.0",
		Port:      4322,
		JetStream: true,
	})
	if err != nil {
		panic(err)
	}

	server.Start()

	nc, err = nats.Connect("nats://127.0.0.1:4322")
	if err != nil {
		panic(err)
	}

	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	generated.RegisterStargateServiceServer(srv, stargateserver.NewServer(logging.Testing(), nc, metrics.NewNoOpMetricsRegistry()))

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	code := m.Run()

	server.Shutdown()

	os.Exit(code)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type Client struct {
	conn           *grpc.ClientConn
	stargateClient generated.StargateServiceClient
}

func NewClient() *Client {
	conn, err := grpc.Dial(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		conn:           conn,
		stargateClient: generated.NewStargateServiceClient(conn),
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) RunStream(t *testing.T, ctx context.Context, organizationID, stackID string, responseChan chan *generated.StargateClientMessage) chan *generated.StargateServerMessage {
	ctx = metadata.AppendToOutgoingContext(ctx, "organization-id", organizationID, "stack-id", stackID)
	stream, err := c.stargateClient.Stargate(ctx)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	incomingMessageChan := make(chan *generated.StargateServerMessage)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-responseChan:
				err := stream.Send(ev)
				require.NoError(t, err)
			}
		}
	}()

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				if err == io.EOF || status.Code(err) == codes.Canceled {
					return
				}

				require.NoError(t, err)
			}

			select {
			case <-ctx.Done():
				return
			case incomingMessageChan <- in:
			}
		}
	}()

	return incomingMessageChan
}
