package internal

import (
	"context"

	"github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/go-libs/logging"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8SClient interface {
	Get(ctx context.Context, resource string, name string) (*unstructured.Unstructured, error)
	Create(ctx context.Context, resource string, o *unstructured.Unstructured) error
	Patch(ctx context.Context, resource, name string, body []byte) error
	EnsureNotExists(ctx context.Context, resource, name string) error
	EnsureNotExistsBySelector(ctx context.Context, resource string, selector labels.Selector) error
	List(ctx context.Context, resource string, selector labels.Selector) ([]unstructured.Unstructured, error)
}

type defaultK8SClient struct {
	restClient *rest.RESTClient
}

func (c defaultK8SClient) Get(ctx context.Context, resource string, name string) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	if err := c.restClient.Get().
		Resource(resource).
		Name(name).
		Do(ctx).
		Into(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (c defaultK8SClient) Create(ctx context.Context, resource string, o *unstructured.Unstructured) error {
	return c.restClient.
		Post().
		Resource(resource).
		Body(o).
		Do(ctx).
		Into(o)
}

func (c defaultK8SClient) Patch(ctx context.Context, resource, name string, body []byte) error {
	return c.restClient.Patch(types.MergePatchType).
		Name(name).
		Body(body).
		Resource(resource).
		Do(ctx).
		Error()
}

func (c defaultK8SClient) EnsureNotExists(ctx context.Context, resource, name string) error {
	return client.IgnoreNotFound(
		c.restClient.Delete().
			Resource(resource).
			Name(name).
			Do(ctx).
			Error(),
	)
}

func (c defaultK8SClient) EnsureNotExistsBySelector(ctx context.Context, resource string, selector labels.Selector) error {
	logging.FromContext(ctx).Debugf("Deleting resources of type %s with selector %s", resource, selector.String())
	return client.IgnoreNotFound(
		c.restClient.Delete().
			Resource(resource).
			VersionedParams(
				&metav1.ListOptions{
					LabelSelector: selector.String(),
				}, scheme.ParameterCodec).
			Do(ctx).
			Error(),
	)
}

func (c defaultK8SClient) List(ctx context.Context, resource string, add labels.Selector) ([]unstructured.Unstructured, error) {
	list := &unstructured.UnstructuredList{}
	if err := c.restClient.Get().
		Resource(resource).
		VersionedParams(&metav1.ListOptions{
			LabelSelector: add.String(),
		}, scheme.ParameterCodec).
		Do(ctx).
		Into(list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

var _ K8SClient = (*defaultK8SClient)(nil)

func NewDefaultK8SClient(restClient *rest.RESTClient) K8SClient {
	return defaultK8SClient{
		restClient: restClient,
	}
}

type cachedK8SClient struct {
	K8SClient
	informerFactory dynamicinformer.DynamicSharedInformerFactory
}

func (c cachedK8SClient) Get(ctx context.Context, resource string, name string) (*unstructured.Unstructured, error) {
	ret, err := c.informerFactory.ForResource(schema.GroupVersionResource{
		Group:    "formance.com",
		Version:  "v1beta1",
		Resource: resource,
	}).
		Lister().
		Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.K8SClient.Get(ctx, resource, name)
		}
		return nil, err
	}
	logging.FromContext(ctx).Debugf("Cache hit for resource %s/%s", resource, name)
	return ret.(*unstructured.Unstructured), nil
}

func (c cachedK8SClient) List(_ context.Context, resource string, selector labels.Selector) ([]unstructured.Unstructured, error) {
	ret, err := c.informerFactory.ForResource(schema.GroupVersionResource{
		Group:    "formance.com",
		Version:  "v1beta1",
		Resource: resource,
	}).
		Lister().
		List(selector)
	if err != nil {
		return nil, err
	}

	return collectionutils.Map(ret, func(from runtime.Object) unstructured.Unstructured {
		return *from.(*unstructured.Unstructured)
	}), nil
}

var _ K8SClient = (*defaultK8SClient)(nil)

func NewCachedK8SClient(k8sClient K8SClient, informerFactory dynamicinformer.DynamicSharedInformerFactory) K8SClient {
	return cachedK8SClient{
		K8SClient:       k8sClient,
		informerFactory: informerFactory,
	}
}
