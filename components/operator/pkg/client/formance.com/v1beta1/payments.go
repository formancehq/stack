package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type PaymentsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.PaymentsList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Payments, error)
	Create(ctx context.Context, Payments *v1beta1.Payments) (*v1beta1.Payments, error)
	Update(ctx context.Context, Payments *v1beta1.Payments) (*v1beta1.Payments, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type paymentsClient struct {
	restClient rest.Interface
}

func (c *paymentsClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.PaymentsList, error) {
	result := v1beta1.PaymentsList{}
	err := c.restClient.
		Get().
		Resource("Payments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *paymentsClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Payments, error) {
	result := v1beta1.Payments{}
	err := c.restClient.
		Get().
		Resource("Payments").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *paymentsClient) Create(ctx context.Context, Payments *v1beta1.Payments) (*v1beta1.Payments, error) {
	result := v1beta1.Payments{}
	err := c.restClient.
		Post().
		Resource("Payments").
		Body(Payments).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *paymentsClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Payments").
		Name(name).
		Do(ctx).
		Error()
}

func (c *paymentsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Payments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *paymentsClient) Update(ctx context.Context, o *v1beta1.Payments) (*v1beta1.Payments, error) {
	result := v1beta1.Payments{}
	err := c.restClient.
		Put().
		Resource("Payments").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
