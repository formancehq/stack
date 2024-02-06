package internal

import (
	"context"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"time"
)

func NewDynamicSharedInformerFactory(client *dynamic.DynamicClient) dynamicinformer.DynamicSharedInformerFactory {
	return dynamicinformer.NewDynamicSharedInformerFactory(client, 5*time.Second)
}

func runInformers(lc fx.Lifecycle, informers []cache.SharedIndexInformer) {
	stopCh := make(chan struct{})
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for _, informer := range informers {
					informer.Run(stopCh)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			close(stopCh)
			return nil
		},
	})
}

func createInformer(factory dynamicinformer.DynamicSharedInformerFactory, resource string, handler cache.ResourceEventHandler) (cache.SharedIndexInformer, error) {
	informer := factory.
		ForResource(schema.GroupVersionResource{
			Group:    "formance.com",
			Version:  "v1beta1",
			Resource: resource,
		}).
		Informer()
	_, err := informer.AddEventHandler(handler)
	if err != nil {
		return nil, err
	}
	return informer, nil
}

func CreateVersionsInformer(factory dynamicinformer.DynamicSharedInformerFactory,
	logger logging.Logger, client MembershipClient) (cache.SharedIndexInformer, error) {
	return createInformer(factory, "versions", VersionsEventHandler(logger, client))
}

func CreateStacksInformer(factory dynamicinformer.DynamicSharedInformerFactory,
	logger logging.Logger, client MembershipClient) (cache.SharedIndexInformer, error) {
	return createInformer(factory, "stacks", StacksEventHandler(logger, client))
}

func CreateRestMapper(config *rest.Config) (meta.RESTMapper, error) {
	discovery := discovery.NewDiscoveryClientForConfigOrDie(config)

	groupResources, err := restmapper.GetAPIGroupResources(discovery)
	if err != nil {
		return nil, err
	}

	return restmapper.NewDiscoveryRESTMapper(groupResources), nil
}

func runMembershipClient(lc fx.Lifecycle, client *membershipClient) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.connect(ctx); err != nil {
				return err
			}
			go func() {
				if err := client.Start(context.Background()); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: client.Stop,
	})
}

func runMembershipListener(lc fx.Lifecycle, client *membershipListener) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go client.Start(context.Background())
			return nil
		},
	})
}

func NewModule(serverAddress string, authenticator Authenticator, clientInfo ClientInfo, opts ...grpc.DialOption) fx.Option {
	return fx.Options(
		fx.Supply(clientInfo),
		fx.Provide(rest.RESTClientFor),
		fx.Provide(dynamic.NewForConfig),
		fx.Provide(NewDynamicSharedInformerFactory),
		fx.Provide(CreateRestMapper),
		fx.Provide(func() *membershipClient {
			return NewMembershipClient(authenticator, clientInfo, serverAddress, opts...)
		}),
		fx.Provide(func(membershipClient *membershipClient) MembershipClient {
			return membershipClient
		}),
		fx.Provide(NewMembershipListener),
		fx.Provide(fx.Annotate(CreateVersionsInformer, fx.ResultTags(`group:"informers"`))),
		fx.Provide(fx.Annotate(CreateStacksInformer, fx.ResultTags(`group:"informers"`))),
		fx.Invoke(runMembershipClient),
		fx.Invoke(runMembershipListener),
		fx.Invoke(fx.Annotate(runInformers, fx.ParamTags(``, `group:"informers"`))),
	)
}
