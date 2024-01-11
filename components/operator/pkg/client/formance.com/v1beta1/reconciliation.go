package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ReconciliationInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.ReconciliationList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Reconciliation, error)
	Create(ctx context.Context, Reconciliation *v1beta1.Reconciliation) (*v1beta1.Reconciliation, error)
	Update(ctx context.Context, Reconciliation *v1beta1.Reconciliation) (*v1beta1.Reconciliation, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type reconciliationClient struct {
	restClient rest.Interface
}

func (c *reconciliationClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.ReconciliationList, error) {
	result := v1beta1.ReconciliationList{}
	err := c.restClient.
		Get().
		Resource("Reconciliations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *reconciliationClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Reconciliation, error) {
	result := v1beta1.Reconciliation{}
	err := c.restClient.
		Get().
		Resource("Reconciliations").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *reconciliationClient) Create(ctx context.Context, Reconciliation *v1beta1.Reconciliation) (*v1beta1.Reconciliation, error) {
	result := v1beta1.Reconciliation{}
	err := c.restClient.
		Post().
		Resource("Reconciliations").
		Body(Reconciliation).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *reconciliationClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Reconciliations").
		Name(name).
		Do(ctx).
		Error()
}

func (c *reconciliationClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Reconciliations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *reconciliationClient) Update(ctx context.Context, o *v1beta1.Reconciliation) (*v1beta1.Reconciliation, error) {
	result := v1beta1.Reconciliation{}
	err := c.restClient.
		Put().
		Resource("Reconciliations").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
