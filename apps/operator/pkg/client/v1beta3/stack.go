package v1beta3

import (
	"context"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type StackInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta3.StackList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error)
	Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type stackClient struct {
	restClient rest.Interface
}

func (c *stackClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta3.StackList, error) {
	result := v1beta3.StackList{}
	err := c.restClient.
		Get().
		Resource("stacks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *stackClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta3.Stack, error) {
	result := v1beta3.Stack{}
	err := c.restClient.
		Get().
		Resource("stacks").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *stackClient) Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	result := v1beta3.Stack{}
	err := c.restClient.
		Post().
		Resource("stacks").
		Body(stack).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *stackClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("stacks").
		Name(name).
		Do(ctx).
		Error()
}

func (c *stackClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("stacks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *stackClient) Update(ctx context.Context, o *v1beta3.Stack) (*v1beta3.Stack, error) {
	result := v1beta3.Stack{}
	err := c.restClient.
		Put().
		Resource("stacks").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
