//nolint:nosnakecase
package grpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/components/agent/internal/grpc/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	oidcclient "github.com/zitadel/oidc/v2/pkg/client"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/metadata"
	controllererrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type K8SClient interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error)
	Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Delete(ctx context.Context, name string) error
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type Authenticator interface {
	authenticate(ctx context.Context) (metadata.MD, error)
}
type AuthenticatorFn func(ctx context.Context) (metadata.MD, error)

func (fn AuthenticatorFn) authenticate(ctx context.Context) (metadata.MD, error) {
	return fn(ctx)
}

func TokenAuthenticator(token string) AuthenticatorFn {
	return func(ctx context.Context) (metadata.MD, error) {
		return metadata.New(map[string]string{"token": token}), nil
	}
}

func BearerAuthenticator(issuer, clientID, clientSecret string) AuthenticatorFn {

	return func(ctx context.Context) (metadata.MD, error) {

		discovery, err := oidcclient.Discover(issuer, http.DefaultClient)
		if err != nil {
			return nil, err
		}

		config := clientcredentials.Config{
			ClientID:     "region_" + clientID,
			ClientSecret: clientSecret,
			TokenURL:     discovery.TokenEndpoint,
		}

		token, err := config.Token(ctx)
		if err != nil {
			return nil, err
		}

		return metadata.New(map[string]string{
			"bearer": token.AccessToken,
		}), nil
	}
}

type client struct {
	stopChan       chan chan error
	grpcClient     generated.ServerClient
	k8sClient      K8SClient
	id             string
	connectClient  generated.Server_ConnectClient
	connectContext context.Context
	connectCancel  func()
	authenticator  Authenticator
	baseUrl        *url.URL
	production     bool
}

func (client *client) Connect(ctx context.Context) error {
	sharedlogging.FromContext(ctx).WithFields(map[string]any{
		"id": client.id,
	}).Infof("Establish connection to server")
	client.connectContext, client.connectCancel = context.WithCancel(ctx)

	md, err := client.authenticator.authenticate(ctx)
	if err != nil {
		return errors.Wrap(err, "authenticating client")
	}

	md.Append("id", client.id)
	md.Append("baseUrl", client.baseUrl.String())
	md.Append("production", func() string {
		if client.production {
			return "true"
		}
		return "false"
	}())

	connectContext := metadata.NewOutgoingContext(client.connectContext, md)
	connectClient, err := client.grpcClient.Join(connectContext)
	if err != nil {
		return err
	}
	client.connectClient = connectClient

	return nil
}

func (client *client) createStack(stack *generated.Stack) *v1beta3.Stack {
	return &v1beta3.Stack{
		ObjectMeta: metav1.ObjectMeta{
			Name: stack.ClusterName,
		},
		Spec: v1beta3.StackSpec{
			DevProperties: v1beta3.DevProperties{
				Debug: false,
				Dev:   false,
			},
			Seed: stack.Seed,
			Auth: v1beta3.StackAuthSpec{
				DelegatedOIDCServer: v1beta3.DelegatedOIDCServerConfiguration{
					Issuer:       stack.AuthConfig.Issuer,
					ClientID:     stack.AuthConfig.ClientId,
					ClientSecret: stack.AuthConfig.ClientSecret,
				},
				StaticClients: []v1beta3.StaticClient{{
					ClientConfiguration: v1beta3.ClientConfiguration{
						Public: true,
					},
					ID:      "fctl",
					Secrets: []string{},
				}},
			},
			Host:   fmt.Sprintf("%s.%s", stack.ClusterName, client.baseUrl.Host),
			Scheme: client.baseUrl.Scheme,
			Stargate: func() *v1beta3.StackStargateConfig {
				if stack.StargateConfig == nil || !stack.StargateConfig.Enabled {
					return nil
				}
				return &v1beta3.StackStargateConfig{
					StargateServerURL: stack.StargateConfig.Url,
				}
			}(),
		},
	}
}

func (client *client) stackFusion(currentStack *v1beta3.Stack, createStack *v1beta3.Stack) *v1beta3.Stack {
	// Keep old values updated by client
	keepServices := currentStack.Spec.Services
	keepSeed := currentStack.Spec.Seed
	keepVersion := currentStack.Spec.Versions
	keepDebug := currentStack.Spec.DevProperties.Debug
	keepDev := currentStack.Spec.DevProperties.Dev

	// Fusion old and new values
	currentStack = createStack

	// Apply old values
	currentStack.Spec.Services = keepServices
	currentStack.Spec.Seed = keepSeed
	currentStack.Spec.Versions = keepVersion
	currentStack.Spec.DevProperties.Debug = keepDebug
	currentStack.Spec.DevProperties.Dev = keepDev

	return currentStack
}

func (client *client) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	watcher, err := client.k8sClient.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	var (
		closed = false
		errCh  = make(chan error, 1)
		msgs   = make(chan *generated.Order)
	)
	go func() {
		for {
			msg := &generated.Order{}
			if err := client.connectClient.RecvMsg(msg); err != nil {
				if err == io.EOF {
					if !closed {
						errCh <- err
					}
					return
				}
				errCh <- err
				return
			}
			select {
			case msgs <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch := <-client.stopChan:
			closed = true
			if err := client.connectClient.CloseSend(); err != nil {
				ch <- err
				//nolint:nilerr
				return nil
			}
			client.connectCancel()
			for {
				msg := &generated.Order{}
				if err := client.connectClient.RecvMsg(msg); err != nil { // Drain messages
					break
				}
			}

			ch <- nil
			return nil
		case err := <-errCh:
			sharedlogging.FromContext(ctx).Errorf("Stream closed with error: %s", err)
			return err
		case k8sUpdate := <-watcher.ResultChan():
			stack := k8sUpdate.Object.(*v1beta3.Stack)
			sharedlogging.FromContext(ctx).Infof("Got update for stack '%s'", stack.Name)
			if err := client.connectClient.SendMsg(&generated.Message{
				Message: &generated.Message_StatusChanged{
					StatusChanged: &generated.StatusChanged{
						StackId: stack.Name,
						Status: func() generated.StackStatus {
							if stack.IsReady() {
								return generated.StackStatus_Ready
							}
							return generated.StackStatus_Progressing
						}(),
					},
				},
			}); err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to send stack status to server: %s", err)
			}
		case msg := <-msgs:
			switch msg := msg.Message.(type) {
			// TODO: Implement UpdateOrCreate
			case *generated.Order_ExistingStack:
				createStack := client.createStack(msg.ExistingStack)
				existingStack, err := client.k8sClient.Get(ctx, createStack.Name, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						if _, err := client.k8sClient.Create(ctx, createStack); err != nil {
							sharedlogging.FromContext(ctx).Errorf("Creating stack cluster side: %s", err)
							continue
						}
						sharedlogging.FromContext(ctx).Infof("Stack %s created", msg.ExistingStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}

				newStack := client.stackFusion(existingStack, createStack)
				if _, err := client.k8sClient.Update(ctx, newStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Updating stack cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s updated", newStack.Name)

			case *generated.Order_DeletedStack:
				if err := client.k8sClient.Delete(ctx, msg.DeletedStack.ClusterName); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Creating deleting cluster side: %s", err)
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s deleted", msg.DeletedStack.ClusterName)
			case *generated.Order_DisabledStack:
				existingStack, err := client.k8sClient.Get(ctx, msg.DisabledStack.ClusterName, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						sharedlogging.FromContext(ctx).Infof("Cannot disable not existing stack: %s", msg.DisabledStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}
				existingStack.Spec.Disabled = true
				if _, err := client.k8sClient.Update(ctx, existingStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Updating stack cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s disabled", msg.DisabledStack.ClusterName)
			case *generated.Order_EnabledStack:
				existingStack, err := client.k8sClient.Get(ctx, msg.EnabledStack.ClusterName, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						sharedlogging.FromContext(ctx).Infof("Cannot enable not existing stack: %s", msg.EnabledStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}
				existingStack.Spec.Disabled = false
				if _, err := client.k8sClient.Update(ctx, existingStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Updating stack cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s enabled", msg.EnabledStack.ClusterName)
			case *generated.Order_Ping:
				sharedlogging.FromContext(ctx).Debugf("Receive ping")
				if err := client.connectClient.SendMsg(&generated.Message{
					Message: &generated.Message_Pong{
						Pong: &generated.Pong{},
					},
				}); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Unable to send pong to server: %s", err)
				}
			}
		}
	}
}

func (client *client) Stop(ctx context.Context) error {
	ch := make(chan error)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case client.stopChan <- ch:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-ch:
			return err
		}
	}
}

func newClient(id string, grpcClient generated.ServerClient, k8sClient K8SClient,
	baseUrl *url.URL, authenticator Authenticator, production bool) *client {
	return &client{
		stopChan:      make(chan chan error),
		grpcClient:    grpcClient,
		id:            id,
		k8sClient:     k8sClient,
		authenticator: authenticator,
		baseUrl:       baseUrl,
		production:    production,
	}
}
