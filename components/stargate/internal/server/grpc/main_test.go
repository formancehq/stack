package grpc_test

import (
	"context"
	"io"
	"net"
	"os"
	"testing"

	"github.com/formancehq/stack/components/stargate/internal/api"
	stargateserver "github.com/formancehq/stack/components/stargate/internal/server/grpc"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/opentelemetry"
	natsserver "github.com/nats-io/nats-server/test"
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
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	api.RegisterStargateServiceServer(srv, stargateserver.NewServer(nc, opentelemetry.NewNoOpMetricsRegistry()))

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	code := m.Run()

	os.Exit(code)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type Client struct {
	conn           *grpc.ClientConn
	stargateClient api.StargateServiceClient
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
		stargateClient: api.NewStargateServiceClient(conn),
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) RunStream(t *testing.T, ctx context.Context, organizationID, stackID string, responseChan chan *api.StargateClientMessage) chan *api.StargateServerMessage {
	ctx = metadata.AppendToOutgoingContext(ctx, "organization-id", organizationID, "stack-id", stackID)
	stream, err := c.stargateClient.Stargate(ctx)
	require.NoError(t, err)

	incomingMessageChan := make(chan *api.StargateServerMessage)
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
