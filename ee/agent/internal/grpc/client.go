//nolint:nosnakecase
package grpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	clientv1beta3 "github.com/formancehq/operator/pkg/client/stack.formance.com/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
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
	Stacks() clientv1beta3.StackInterface
	Versions() clientv1beta3.VersionsInterface
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

type ClientInfo struct {
	ID         string
	BaseUrl    *url.URL
	Production bool
	Version    string
}

const (
	metadataID         = "id"
	metadataBaseUrl    = "baseUrl"
	metadataProduction = "production"
	metadataVersion    = "version"
)

type client struct {
	clientInfo     ClientInfo
	stopChan       chan chan error
	grpcClient     generated.ServerClient
	k8sClient      K8SClient
	connectClient  generated.Server_ConnectClient
	connectContext context.Context
	connectCancel  func()
	authenticator  Authenticator
}

func (client *client) Connect(ctx context.Context) error {
	sharedlogging.FromContext(ctx).WithFields(map[string]any{
		"id": client.clientInfo.ID,
	}).Infof("Establish connection to server")
	client.connectContext, client.connectCancel = context.WithCancel(ctx)

	md, err := client.authenticator.authenticate(ctx)
	if err != nil {
		return errors.Wrap(err, "authenticating client")
	}

	md.Append("id", client.clientInfo.ID)
	md.Append("baseUrl", client.clientInfo.BaseUrl.String())
	md.Append("production", func() string {
		if client.clientInfo.Production {
			return "true"
		}
		return "false"
	}())
	md.Append("version", client.clientInfo.Version)

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
					ID: "fctl",
				}},
			},
			Host:   fmt.Sprintf("%s.%s", stack.ClusterName, client.clientInfo.BaseUrl.Host),
			Scheme: client.clientInfo.BaseUrl.Scheme,
			Stargate: func() *v1beta3.StackStargateConfig {
				if stack.StargateConfig == nil || !stack.StargateConfig.Enabled {
					return nil
				}
				return &v1beta3.StackStargateConfig{
					StargateServerURL: stack.StargateConfig.Url,
				}
			}(),
			Disabled: stack.Disabled,
			Versions: stack.Versions,
		},
	}
}

func mergeStack(currentStack *v1beta3.Stack, into *v1beta3.Stack) *v1beta3.Stack {
	into.SetResourceVersion(currentStack.GetResourceVersion())
	into.Spec.Services = currentStack.Spec.Services
	into.Spec.Seed = currentStack.Spec.Seed
	into.Spec.DevProperties.Debug = currentStack.Spec.DevProperties.Debug
	into.Spec.DevProperties.Dev = currentStack.Spec.DevProperties.Dev
	if into.Spec.Versions == "" {
		into.Spec.Versions = currentStack.Spec.Versions
	}
	return into
}

func (client *client) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stacksWatcher, err := client.k8sClient.Stacks().Watch(ctx, metav1.ListOptions{
		Watch: true,
	})
	if err != nil {
		return errors.Wrap(err, "trying to watch stacks")
	}
	defer stacksWatcher.Stop()

	versionsWatcher, err := client.k8sClient.Versions().Watch(ctx, metav1.ListOptions{
		Watch: true,
	})
	if err != nil {
		return errors.Wrap(err, "trying to watch versions")
	}
	defer versionsWatcher.Stop()

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
		case k8sUpdate := <-stacksWatcher.ResultChan():
			stack := k8sUpdate.Object.(*v1beta3.Stack)
			var status generated.StackStatus
			switch k8sUpdate.Type {
			case watch.Deleted:
				status = generated.StackStatus_Deleted
			default:
				switch {
				case stack.Spec.Disabled:
					status = generated.StackStatus_Disabled
				case stack.Status.Ready:
					status = generated.StackStatus_Ready
				default:
					status = generated.StackStatus_Progressing
				}
			}

			sharedlogging.FromContext(ctx).Infof("Got update for stack '%s': %s", stack.Name, status)

			if err := client.connectClient.SendMsg(&generated.Message{
				Message: &generated.Message_StatusChanged{
					StatusChanged: &generated.StatusChanged{
						ClusterName: stack.Name,
						Status:      status,
					},
				},
			}); err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to send stack status to server: %s", err)
			}
		case k8sUpdate := <-versionsWatcher.ResultChan():
			version := k8sUpdate.Object.(*v1beta3.Versions)

			var msg *generated.Message
			switch k8sUpdate.Type {
			default:
				sharedlogging.FromContext(ctx).Error("Watch type '%s' not handled for versions", k8sUpdate.Type)
			case watch.Deleted:
				sharedlogging.FromContext(ctx).Infof("Detect versions '%s' as deleted", version.Name)
				msg = &generated.Message{
					Message: &generated.Message_DeletedVersion{
						DeletedVersion: &generated.DeletedVersion{
							Name: version.Name,
						},
					},
				}
			case watch.Added, watch.Modified:
				vOfVersion := reflect.ValueOf(version.Spec)
				tOfVersion := reflect.TypeOf(version.Spec)
				versionsMap := make(map[string]string)
				for i := 0; i < tOfVersion.NumField(); i++ {
					jsonTag := tOfVersion.Field(i).Tag.Get("json")
					name := strings.Split(jsonTag, ",")[0]
					versionsMap[name] = vOfVersion.Field(i).String()
				}

				switch k8sUpdate.Type {
				case watch.Added:
					sharedlogging.FromContext(ctx).
						WithFields(collectionutils.ConvertMap(versionsMap, collectionutils.ToAny[string])).
						Infof("Detect versions '%s' added", version.Name)
					msg = &generated.Message{
						Message: &generated.Message_AddedVersion{
							AddedVersion: &generated.AddedVersion{
								Name:     version.Name,
								Versions: versionsMap,
							},
						},
					}
				case watch.Modified:
					sharedlogging.FromContext(ctx).
						WithFields(collectionutils.ConvertMap(versionsMap, collectionutils.ToAny[string])).
						Infof("Detect versions '%s' modified", version.Name)
					msg = &generated.Message{
						Message: &generated.Message_UpdatedVersion{
							UpdatedVersion: &generated.UpdatedVersion{
								Name:     version.Name,
								Versions: versionsMap,
							},
						},
					}
				}
			}

			if err := client.connectClient.SendMsg(msg); err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to send version update: %s", err)
			}

		case msg := <-msgs:
			switch msg := msg.Message.(type) {
			// TODO: Implement UpdateOrCreate
			case *generated.Order_ExistingStack:
				createStack := client.createStack(msg.ExistingStack)
				existingStack, err := client.k8sClient.Stacks().Get(ctx, createStack.Name, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						if _, err := client.k8sClient.Stacks().Create(ctx, createStack); err != nil {
							sharedlogging.FromContext(ctx).Errorf("Creating stack cluster side: %s", err)
							continue
						}
						sharedlogging.FromContext(ctx).Infof("Stack %s created", msg.ExistingStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}

				newStack := mergeStack(existingStack, createStack)
				if _, err := client.k8sClient.Stacks().Update(ctx, newStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Updating stack cluster side: %s", err)
					continue
				}

				sharedlogging.FromContext(ctx).Infof("Stack %s updated cluster side", newStack.Name)

			case *generated.Order_DeletedStack:
				if err := client.k8sClient.Stacks().Delete(ctx, msg.DeletedStack.ClusterName); err != nil {
					if controllererrors.IsNotFound(err) {
						sharedlogging.FromContext(ctx).Infof("Cannot delete not existing stack: %s", msg.DeletedStack.ClusterName)

						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Deleting cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s deleted", msg.DeletedStack.ClusterName)
			case *generated.Order_DisabledStack:
				sharedlogging.FromContext(ctx).Infof("Incomming DISABLING: %s", msg.DisabledStack.ClusterName)

				existingStack, err := client.k8sClient.Stacks().Get(ctx, msg.DisabledStack.ClusterName, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						sharedlogging.FromContext(ctx).Infof("Cannot disable not existing stack: %s", msg.DisabledStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}
				existingStack.Spec.Disabled = true
				if _, err := client.k8sClient.Stacks().Update(ctx, existingStack); err != nil {
					sharedlogging.FromContext(ctx).Errorf("Updating stack cluster side: %s", err)
					continue
				}
				sharedlogging.FromContext(ctx).Infof("Stack %s disabled", msg.DisabledStack.ClusterName)
			case *generated.Order_EnabledStack:
				sharedlogging.FromContext(ctx).Infof("Incomming ENABLING: %s", msg.EnabledStack.ClusterName)
				existingStack, err := client.k8sClient.Stacks().Get(ctx, msg.EnabledStack.ClusterName, metav1.GetOptions{})
				if err != nil {
					if controllererrors.IsNotFound(err) {
						sharedlogging.FromContext(ctx).Infof("Cannot enable not existing stack: %s", msg.EnabledStack.ClusterName)
						continue
					}
					sharedlogging.FromContext(ctx).Errorf("Reading stack cluster side: %s", err)
					continue
				}
				existingStack.Spec.Disabled = false
				if _, err := client.k8sClient.Stacks().Update(ctx, existingStack); err != nil {
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

func newClient(grpcClient generated.ServerClient, k8sClient K8SClient, authenticator Authenticator, clientInfo ClientInfo) *client {
	return &client{
		stopChan:      make(chan chan error),
		grpcClient:    grpcClient,
		k8sClient:     k8sClient,
		authenticator: authenticator,
		clientInfo:    clientInfo,
	}
}
