//nolint:nosnakecase
package grpc

import (
	"context"
	"net"
	"net/url"
	"testing"
	"time"

	clientv1beta3 "github.com/formancehq/operator/pkg/client/v1beta3"
	"github.com/pkg/errors"

	"google.golang.org/grpc/metadata"

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

type stacksClient struct {
	stacks map[string]*v1beta3.Stack
}

func (k *stacksClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta3.StackList, error) {
	return nil, errors.New("not implemented")
}

func (k *stacksClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return watch.NewFake(), nil
}

func (k *stacksClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error) {
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

func (k *stacksClient) Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	k.stacks[stack.Name] = stack
	return stack, nil
}

func (k *stacksClient) Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	if k.stacks == nil {
		k.stacks = make(map[string]*v1beta3.Stack)
	}
	k.stacks[stack.Name] = stack
	return nil, nil
}

func (k *stacksClient) Delete(ctx context.Context, name string) error {
	delete(k.stacks, name)
	return nil
}

type versionsClient struct {
	versions map[string]*v1beta3.Versions
}

func (k *versionsClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta3.VersionsList, error) {
	return nil, errors.New("not implemented")
}

func (k *versionsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return watch.NewFake(), nil
}

func (k *versionsClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Versions, error) {
	versions, ok := k.versions[name]
	if !ok {
		return nil, &apierrors.StatusError{
			ErrStatus: metav1.Status{
				Reason: metav1.StatusReasonNotFound,
			},
		}
	}
	return versions, nil
}

func (k *versionsClient) Update(ctx context.Context, versions *v1beta3.Versions) (*v1beta3.Versions, error) {
	k.versions[versions.Name] = versions
	return versions, nil
}

func (k *versionsClient) Create(ctx context.Context, versions *v1beta3.Versions) (*v1beta3.Versions, error) {
	if k.versions == nil {
		k.versions = make(map[string]*v1beta3.Versions)
	}
	k.versions[versions.Name] = versions
	return nil, nil
}

func (k *versionsClient) Delete(ctx context.Context, name string) error {
	delete(k.versions, name)
	return nil
}

type k8sClient struct {
	stacksClient   *stacksClient
	versionsClient *versionsClient
}

func (k k8sClient) Stacks() clientv1beta3.StackInterface {
	return k.stacksClient
}

func (k k8sClient) Versions() clientv1beta3.VersionsInterface {
	return k.versionsClient
}

var _ K8SClient = &k8sClient{}

func newK8SClient() *k8sClient {
	return &k8sClient{
		stacksClient: &stacksClient{
			stacks: map[string]*v1beta3.Stack{},
		},
		versionsClient: &versionsClient{
			versions: map[string]*v1beta3.Versions{},
		},
	}
}

func TestModule(t *testing.T) {
	mockServer := &mockServer{}
	k8sClient := newK8SClient()

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

	clientID := uuid.NewString()
	version := "v1.0.0"
	app := fx.New(
		fx.NopLogger,
		fx.Supply(fx.Annotate(k8sClient, fx.As(new(K8SClient)))),
		NewModule(lis.Addr().String(), TokenAuthenticator(""), ClientInfo{
			ID:         clientID,
			BaseUrl:    baseUrl,
			Production: false,
			Version:    version,
		},
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
		return len(k8sClient.stacksClient.stacks) == 1
	}, time.Second, 100*time.Millisecond)
	require.NotEmpty(t, k8sClient.stacksClient.stacks[createdStack.ClusterName])
	require.Equal(t, createdStack.ClusterName, k8sClient.stacksClient.stacks[createdStack.ClusterName].Name)

	md, ok := metadata.FromIncomingContext(mockServer.connectServer.Context())
	require.True(t, ok)
	require.Equal(t, []string{clientID}, md.Get(metadataID))
	require.Equal(t, []string{baseUrl.String()}, md.Get(metadataBaseUrl))
	require.Equal(t, []string{"false"}, md.Get(metadataProduction))
	require.Equal(t, []string{version}, md.Get(metadataVersion))
}
