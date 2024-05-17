package internal

import (
	"context"
	"reflect"
	"sort"
	"time"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func NewDynamicSharedInformerFactory(client *dynamic.DynamicClient, resyncPeriod time.Duration) dynamicinformer.DynamicSharedInformerFactory {
	return dynamicinformer.NewDynamicSharedInformerFactory(client, resyncPeriod)
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

func NewK8SConfig(kubeConfigPath string) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		logging.Info("Does not seems to be in cluster, trying to load k8s client from kube config file")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	}

	config.GroupVersion = &v1beta1.GroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.APIPath = "/apis"

	return config, nil
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
	restMapper meta.RESTMapper, logger logging.Logger, client MembershipClient) error {

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

func retrieveModuleList(ctx context.Context, config *rest.Config) (modules []string, eeModules []string, err error) {
	config = rest.CopyConfig(config)
	config.GroupVersion = &apiextensions.SchemeGroupVersion

	apiextensionsClient, err := apiextensionsv1client.NewForConfig(config)
	if err != nil {
		return
	}

	crds, err := apiextensionsClient.CustomResourceDefinitions().List(ctx, metav1.ListOptions{
		LabelSelector: "formance.com/kind=module",
	})
	if err != nil {
		return
	}

	modules = collectionutils.Map(crds.Items, func(item v1.CustomResourceDefinition) string {
		return item.Status.AcceptedNames.Singular
	})

	eeModules = collectionutils.Reduce(crds.Items, func(acc []string, item v1.CustomResourceDefinition) []string {
		if item.Labels["formance.com/is-ee"] == "true" {
			return append(acc, item.Status.AcceptedNames.Singular)
		}
		return acc
	}, []string{})

	sort.Strings(modules)
	sort.Strings(eeModules)

	return
}

func runMembershipClient(lc fx.Lifecycle, membershipClient *membershipClient, logger logging.Logger, config *rest.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			modules, eeModules, err := retrieveModuleList(ctx, config)
			if err != nil {
				return err
			}

			if err := membershipClient.connect(logging.ContextWithLogger(ctx, logger), modules, eeModules); err != nil {
				return err
			}
			go func() {
				if err := membershipClient.Start(logging.ContextWithLogger(context.Background(), logger)); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: membershipClient.Stop,
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

func NewModule(serverAddress string, authenticator Authenticator, clientInfo ClientInfo, resyncPeriod time.Duration, opts ...grpc.DialOption) fx.Option {
	return fx.Options(
		fx.Supply(clientInfo),
		fx.Provide(rest.RESTClientFor),
		fx.Provide(dynamic.NewForConfig),
		fx.Provide(func(client *dynamic.DynamicClient) dynamicinformer.DynamicSharedInformerFactory {
			return NewDynamicSharedInformerFactory(client, resyncPeriod)
		}),
		fx.Provide(func(restClient *rest.RESTClient, informerFactory dynamicinformer.DynamicSharedInformerFactory) K8SClient {
			return NewCachedK8SClient(NewDefaultK8SClient(restClient), informerFactory)
		}),
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
