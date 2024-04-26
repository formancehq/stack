package internal

import (
	"context"
	"reflect"
	"time"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
)

func NewDynamicSharedInformerFactory(client *dynamic.DynamicClient) dynamicinformer.DynamicSharedInformerFactory {
	return dynamicinformer.NewDynamicSharedInformerFactory(client, 2500*time.Millisecond)
}

func runInformers(lc fx.Lifecycle, factory dynamicinformer.DynamicSharedInformerFactory) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			stopCh := make(chan struct{})
			factory.Start(stopCh)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			factory.Shutdown()
			return nil
		},
	})
}

func createInformer(factory dynamicinformer.DynamicSharedInformerFactory, resource string, handler cache.ResourceEventHandler) error {
	informer := factory.
		ForResource(schema.GroupVersionResource{
			Group:    "formance.com",
			Version:  "v1beta1",
			Resource: resource,
		}).
		Informer()

	_, err := informer.AddEventHandler(handler)
	if err != nil {
		return errors.Wrap(err, "unable to add event handler")
	}
	return nil
}

func CreateVersionsInformer(factory dynamicinformer.DynamicSharedInformerFactory,
	logger logging.Logger, client MembershipClient) error {
	logger = logger.WithFields(map[string]any{
		"component": "versions",
	})
	logger.Info("Creating informer")
	return createInformer(factory, "versions", VersionsEventHandler(logger, client))
}

func CreateStacksInformer(factory dynamicinformer.DynamicSharedInformerFactory,
	logger logging.Logger, client MembershipClient, stacks InMemoryStacksModules) error {
	logger = logger.WithFields(map[string]any{
		"component": "stacks",
	})
	logger.Info("Creating informer")
	return createInformer(factory, "stacks", NewStackEventHandler(logger, client, stacks))
}

func CreateModulesInformers(factory dynamicinformer.DynamicSharedInformerFactory,
	restMapper meta.RESTMapper,
	logger logging.Logger, client MembershipClient) error {

	for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
		object := reflect.New(rtype).Interface()
		_, ok := object.(v1beta1.Module)
		if !ok {
			continue
		}

		restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		logger = logger.WithFields(map[string]any{
			"component": restMapping.Resource.Resource,
		})

		logger.Info("Creating informer")
		if err := createInformer(factory, restMapping.Resource.Resource, NewModuleEventHandler(logger, client)); err != nil {
			return err
		}
	}
	return nil
}

func CreateRestMapper(config *rest.Config) (meta.RESTMapper, error) {
	discovery := discovery.NewDiscoveryClientForConfigOrDie(config)

	groupResources, err := restmapper.GetAPIGroupResources(discovery)
	if err != nil {
		return nil, err
	}

	return restmapper.NewDiscoveryRESTMapper(groupResources), nil
}

func runMembershipClient(lc fx.Lifecycle, client *membershipClient, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.connect(logging.ContextWithLogger(ctx, logger)); err != nil {
				return err
			}
			go func() {
				if err := client.Start(logging.ContextWithLogger(context.Background(), logger)); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: client.Stop,
	})
}

func runMembershipListener(lc fx.Lifecycle, client *membershipListener, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go client.Start(logging.ContextWithLogger(context.Background(), logger))
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
		fx.Provide(func() InMemoryStacksModules {
			return map[string][]string{}
		}),
		fx.Provide(NewMembershipListener),
		fx.Invoke(CreateVersionsInformer),
		fx.Invoke(CreateStacksInformer),
		fx.Invoke(CreateModulesInformers),
		fx.Invoke(runMembershipClient),
		fx.Invoke(runMembershipListener),
		fx.Invoke(runInformers),
	)
}
