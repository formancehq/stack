//nolint:nosnakecase
package grpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/components/agent/internal/grpc/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/usagerecord"
	oidcclient "github.com/zitadel/oidc/pkg/client"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/metadata"
	controllererrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8SClient interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error)
	Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Delete(ctx context.Context, name string) error
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

	connectContext := metadata.NewOutgoingContext(client.connectContext, md)
	connectClient, err := client.grpcClient.Connect(connectContext, &generated.ConnectRequest{
		Id:         client.id,
		BaseUrl:    client.baseUrl.String(),
		Production: client.production,
	})
	if err != nil {
		return err
	}
	client.connectClient = connectClient
	return nil
}

func (client *client) k8sStackFromProtobuf(stack *generated.Stack) *v1beta3.Stack {
	return &v1beta3.Stack{
		ObjectMeta: metav1.ObjectMeta{
			Name: stack.ClusterName,
		},
		Spec: v1beta3.StackSpec{
			DevProperties: v1beta3.DevProperties{
				Debug: true,
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

func (client *client) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		closed = false
		errCh  = make(chan error, 1)
		msgs   = make(chan *generated.ConnectResponse)
	)
	go func() {
		for {
			msg := &generated.ConnectResponse{}
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
				msg := &generated.ConnectResponse{}
				if err := client.connectClient.RecvMsg(msg); err != nil { // Drain messages
					break
				}
			}

			ch <- nil
			return nil
		case err := <-errCh:
			sharedlogging.FromContext(ctx).Errorf("Stream closed with error: %s", err)
			return err
		case msg := <-msgs:
			switch msg := msg.Message.(type) {
			case *generated.ConnectResponse_ExistingStack:
				crd := client.k8sStackFromProtobuf(msg.ExistingStack)
				existingStack, err := client.k8sClient.Get(ctx, crd.Name, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						if _, err := client.k8sClient.Create(ctx, crd); err != nil {
							sharedlogging.FromContext(ctx).Errorf("creating stack cluster side: %s", err)
						}
						sharedlogging.FromContext(ctx).Infof("stack %s created", crd.Name)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("reading stack cluster side: %s", err)
					continue
				}
				existingStack.Spec = crd.Spec
				if _, err := client.k8sClient.Update(ctx, existingStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("updating stack cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("stack %s updated", crd.Name)

			case *generated.ConnectResponse_DeletedStack:
				if err := client.k8sClient.Delete(ctx, msg.DeletedStack.ClusterName); err != nil {
					sharedlogging.FromContext(ctx).Errorf("creating deleting cluster side: %s", err)
				}
			case *generated.ConnectResponse_UpdateUsageReport:
				total, err := CountDocument(msg.UpdateUsageReport.ClusterName)
				if err != nil {
					sharedlogging.FromContext(ctx).Errorf("counting documents: %s", err)
					continue
				}

				stripe.Key = msg.UpdateUsageReport.StripeKey

				params := &stripe.UsageRecordParams{
					Action:           stripe.String(stripe.UsageRecordActionSet),
					SubscriptionItem: stripe.String(msg.UpdateUsageReport.StripeSubscriptionId),
					Quantity:         stripe.Int64(total),
					Timestamp:        stripe.Int64(time.Now().Unix()),
				}
				_, err = usagerecord.New(params)

				if err != nil {
					sharedlogging.FromContext(ctx).Errorf("creating usage record: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("usage record: %s", total)
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
