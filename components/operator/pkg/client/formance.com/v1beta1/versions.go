package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type VersionsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.VersionsList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Versions, error)
	Create(ctx context.Context, Versions *v1beta1.Versions) (*v1beta1.Versions, error)
	Update(ctx context.Context, Versions *v1beta1.Versions) (*v1beta1.Versions, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type VersionsClient struct {
	restClient rest.Interface
}

func (c *VersionsClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.VersionsList, error) {
	result := v1beta1.VersionsList{}
	err := c.restClient.
		Get().
		Resource("Versions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *VersionsClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Versions, error) {
	result := v1beta1.Versions{}
	err := c.restClient.
		Get().
		Resource("Versions").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *VersionsClient) Create(ctx context.Context, Versions *v1beta1.Versions) (*v1beta1.Versions, error) {
	result := v1beta1.Versions{}
	err := c.restClient.
		Post().
		Resource("Versions").
		Body(Versions).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *VersionsClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Versions").
		Name(name).
		Do(ctx).
		Error()
}

func (c *VersionsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Versions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *VersionsClient) Update(ctx context.Context, o *v1beta1.Versions) (*v1beta1.Versions, error) {
	result := v1beta1.Versions{}
	err := c.restClient.
		Put().
		Resource("Versions").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
