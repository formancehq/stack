package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type LedgerInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.LedgerList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Ledger, error)
	Create(ctx context.Context, Ledger *v1beta1.Ledger) (*v1beta1.Ledger, error)
	Update(ctx context.Context, Ledger *v1beta1.Ledger) (*v1beta1.Ledger, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type LedgerClient struct {
	restClient rest.Interface
}

func (c *LedgerClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.LedgerList, error) {
	result := v1beta1.LedgerList{}
	err := c.restClient.
		Get().
		Resource("Ledgers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *LedgerClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Ledger, error) {
	result := v1beta1.Ledger{}
	err := c.restClient.
		Get().
		Resource("Ledgers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *LedgerClient) Create(ctx context.Context, Ledger *v1beta1.Ledger) (*v1beta1.Ledger, error) {
	result := v1beta1.Ledger{}
	err := c.restClient.
		Post().
		Resource("Ledgers").
		Body(Ledger).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *LedgerClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Ledgers").
		Name(name).
		Do(ctx).
		Error()
}

func (c *LedgerClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Ledgers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *LedgerClient) Update(ctx context.Context, o *v1beta1.Ledger) (*v1beta1.Ledger, error) {
	result := v1beta1.Ledger{}
	err := c.restClient.
		Put().
		Resource("Ledgers").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
