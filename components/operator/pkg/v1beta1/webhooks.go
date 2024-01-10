package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type WebhooksInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.WebhooksList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Webhooks, error)
	Create(ctx context.Context, Webhooks *v1beta1.Webhooks) (*v1beta1.Webhooks, error)
	Update(ctx context.Context, Webhooks *v1beta1.Webhooks) (*v1beta1.Webhooks, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type webhooksClient struct {
	restClient rest.Interface
}

func (c *webhooksClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.WebhooksList, error) {
	result := v1beta1.WebhooksList{}
	err := c.restClient.
		Get().
		Resource("Webhooks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *webhooksClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Webhooks, error) {
	result := v1beta1.Webhooks{}
	err := c.restClient.
		Get().
		Resource("Webhooks").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *webhooksClient) Create(ctx context.Context, Webhooks *v1beta1.Webhooks) (*v1beta1.Webhooks, error) {
	result := v1beta1.Webhooks{}
	err := c.restClient.
		Post().
		Resource("Webhooks").
		Body(Webhooks).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *webhooksClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Webhooks").
		Name(name).
		Do(ctx).
		Error()
}

func (c *webhooksClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Webhooks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *webhooksClient) Update(ctx context.Context, o *v1beta1.Webhooks) (*v1beta1.Webhooks, error) {
	result := v1beta1.Webhooks{}
	err := c.restClient.
		Put().
		Resource("Webhooks").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
