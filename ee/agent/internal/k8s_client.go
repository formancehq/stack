package internal

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
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