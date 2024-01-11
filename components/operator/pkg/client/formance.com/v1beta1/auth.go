package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type AuthInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.AuthList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Auth, error)
	Create(ctx context.Context, Auth *v1beta1.Auth) (*v1beta1.Auth, error)
	Update(ctx context.Context, Auth *v1beta1.Auth) (*v1beta1.Auth, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type AuthClient struct {
	restClient rest.Interface
}

func (c *AuthClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.AuthList, error) {
	result := v1beta1.AuthList{}
	err := c.restClient.
		Get().
		Resource("Auths").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *AuthClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Auth, error) {
	result := v1beta1.Auth{}
	err := c.restClient.
		Get().
		Resource("Auths").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *AuthClient) Create(ctx context.Context, Auth *v1beta1.Auth) (*v1beta1.Auth, error) {
	result := v1beta1.Auth{}
	err := c.restClient.
		Post().
		Resource("Auths").
		Body(Auth).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *AuthClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Auths").
		Name(name).
		Do(ctx).
		Error()
}

func (c *AuthClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Auths").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *AuthClient) Update(ctx context.Context, o *v1beta1.Auth) (*v1beta1.Auth, error) {
	result := v1beta1.Auth{}
	err := c.restClient.
		Put().
		Resource("Auths").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
