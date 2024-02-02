package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type GatewayInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.GatewayList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Gateway, error)
	Create(ctx context.Context, Gateway *v1beta1.Gateway) (*v1beta1.Gateway, error)
	Update(ctx context.Context, Gateway *v1beta1.Gateway) (*v1beta1.Gateway, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type gatewayClient struct {
	restClient rest.Interface
}

func (c *gatewayClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.GatewayList, error) {
	result := v1beta1.GatewayList{}
	err := c.restClient.
		Get().
		Resource("Gateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *gatewayClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Gateway, error) {
	result := v1beta1.Gateway{}
	err := c.restClient.
		Get().
		Resource("Gateways").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *gatewayClient) Create(ctx context.Context, Gateway *v1beta1.Gateway) (*v1beta1.Gateway, error) {
	result := v1beta1.Gateway{}
	err := c.restClient.
		Post().
		Resource("Gateways").
		Body(Gateway).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *gatewayClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Gateways").
		Name(name).
		Do(ctx).
		Error()
}

func (c *gatewayClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Gateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *gatewayClient) Update(ctx context.Context, o *v1beta1.Gateway) (*v1beta1.Gateway, error) {
	result := v1beta1.Gateway{}
	err := c.restClient.
		Put().
		Resource("Gateways").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
