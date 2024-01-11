package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type WalletsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.WalletsList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Wallets, error)
	Create(ctx context.Context, Wallets *v1beta1.Wallets) (*v1beta1.Wallets, error)
	Update(ctx context.Context, Wallets *v1beta1.Wallets) (*v1beta1.Wallets, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type walletsClient struct {
	restClient rest.Interface
}

func (c *walletsClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.WalletsList, error) {
	result := v1beta1.WalletsList{}
	err := c.restClient.
		Get().
		Resource("Wallets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *walletsClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Wallets, error) {
	result := v1beta1.Wallets{}
	err := c.restClient.
		Get().
		Resource("Wallets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *walletsClient) Create(ctx context.Context, Wallets *v1beta1.Wallets) (*v1beta1.Wallets, error) {
	result := v1beta1.Wallets{}
	err := c.restClient.
		Post().
		Resource("Wallets").
		Body(Wallets).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *walletsClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Wallets").
		Name(name).
		Do(ctx).
		Error()
}

func (c *walletsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Wallets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *walletsClient) Update(ctx context.Context, o *v1beta1.Wallets) (*v1beta1.Wallets, error) {
	result := v1beta1.Wallets{}
	err := c.restClient.
		Put().
		Resource("Wallets").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
