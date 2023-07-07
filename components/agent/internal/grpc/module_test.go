//nolint:nosnakecase
package grpc

import (
	"context"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/components/agent/internal/grpc/generated"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type mockServer struct {
	generated.UnimplementedServerServer
	connectServer generated.Server_JoinServer
}

func (m *mockServer) Join(server generated.Server_JoinServer) error {
	m.connectServer = server
	<-server.Context().Done()
	return nil
}

type k8sClient struct {
	stacks map[string]*v1beta3.Stack
}

func (k *k8sClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return watch.NewFake(), nil
}

func (k *k8sClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error) {
	stack, ok := k.stacks[name]
	if !ok {
		return nil, &apierrors.StatusError{
			ErrStatus: metav1.Status{
				Reason: metav1.StatusReasonNotFound,
			},
		}
	}
	return stack, nil
}

func (k *k8sClient) Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	k.stacks[stack.Name] = stack
	return stack, nil
}

func (k *k8sClient) Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	if k.stacks == nil {
		k.stacks = make(map[string]*v1beta3.Stack)
	}
	k.stacks[stack.Name] = stack
	return nil, nil
}

func (k *k8sClient) Delete(ctx context.Context, name string) error {
	delete(k.stacks, name)
	return nil
}

var _ K8SClient = &k8sClient{}

func TestModule(t *testing.T) {
	mockServer := &mockServer{}
	k8sClient := &k8sClient{}

	grpcServer := grpc.NewServer()
	generated.RegisterServerServer(grpcServer, mockServer)

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		require.NoError(t, grpcServer.Serve(lis))
	}()
	defer func() {
		grpcServer.Stop()
	}()

	baseUrl, err := url.Parse("http://example.net")
	require.NoError(t, err)

	app := fx.New(
		fx.NopLogger,
		fx.Supply(fx.Annotate(k8sClient, fx.As(new(K8SClient)))),
		NewModule(uuid.NewString(), lis.Addr().String(), baseUrl, false, TokenAuthenticator(""),
			grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	require.NoError(t, app.Start(context.Background()))
	defer func() {
		require.NoError(t, app.Stop(context.Background()))
	}()

	require.Eventually(t, func() bool {
		return mockServer.connectServer != nil
	}, time.Second, 100*time.Millisecond)

	createdStack := generated.Stack{
		ClusterName: uuid.NewString(),
		Seed:        uuid.NewString(),
		AuthConfig: &generated.AuthConfig{
			ClientId:     uuid.NewString(),
			ClientSecret: uuid.NewString(),
			Issuer:       uuid.NewString(),
		},
		StaticClients: []*generated.AuthClient{{
			Public: true,
			Id:     uuid.NewString(),
		}},
	}
	require.NoError(t, mockServer.connectServer.Send(&generated.Order{
		Message: &generated.Order_ExistingStack{
			ExistingStack: &createdStack,
		},
	}))
	require.Eventually(t, func() bool {
		return len(k8sClient.stacks) == 1
	}, time.Second, 100*time.Millisecond)
	require.NotEmpty(t, k8sClient.stacks[createdStack.ClusterName])
	require.Equal(t, createdStack.ClusterName, k8sClient.stacks[createdStack.ClusterName].Name)
}
