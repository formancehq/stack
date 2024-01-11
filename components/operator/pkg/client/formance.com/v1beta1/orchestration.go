package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OrchestrationInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.OrchestrationList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Orchestration, error)
	Create(ctx context.Context, Orchestration *v1beta1.Orchestration) (*v1beta1.Orchestration, error)
	Update(ctx context.Context, Orchestration *v1beta1.Orchestration) (*v1beta1.Orchestration, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type OrchestrationClient struct {
	restClient rest.Interface
}

func (c *OrchestrationClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.OrchestrationList, error) {
	result := v1beta1.OrchestrationList{}
	err := c.restClient.
		Get().
		Resource("Orchestrations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *OrchestrationClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Orchestration, error) {
	result := v1beta1.Orchestration{}
	err := c.restClient.
		Get().
		Resource("Orchestrations").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *OrchestrationClient) Create(ctx context.Context, Orchestration *v1beta1.Orchestration) (*v1beta1.Orchestration, error) {
	result := v1beta1.Orchestration{}
	err := c.restClient.
		Post().
		Resource("Orchestrations").
		Body(Orchestration).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *OrchestrationClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Orchestrations").
		Name(name).
		Do(ctx).
		Error()
}

func (c *OrchestrationClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Orchestrations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *OrchestrationClient) Update(ctx context.Context, o *v1beta1.Orchestration) (*v1beta1.Orchestration, error) {
	result := v1beta1.Orchestration{}
	err := c.restClient.
		Put().
		Resource("Orchestrations").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
