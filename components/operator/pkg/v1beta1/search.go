package v1beta1

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SearchInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.SearchList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Search, error)
	Create(ctx context.Context, Search *v1beta1.Search) (*v1beta1.Search, error)
	Update(ctx context.Context, Search *v1beta1.Search) (*v1beta1.Search, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
}

type SearchClient struct {
	restClient rest.Interface
}

func (c *SearchClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.SearchList, error) {
	result := v1beta1.SearchList{}
	err := c.restClient.
		Get().
		Resource("Searchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *SearchClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.Search, error) {
	result := v1beta1.Search{}
	err := c.restClient.
		Get().
		Resource("Searchs").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *SearchClient) Create(ctx context.Context, Search *v1beta1.Search) (*v1beta1.Search, error) {
	result := v1beta1.Search{}
	err := c.restClient.
		Post().
		Resource("Searchs").
		Body(Search).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *SearchClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Resource("Searchs").
		Name(name).
		Do(ctx).
		Error()
}

func (c *SearchClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("Searchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *SearchClient) Update(ctx context.Context, o *v1beta1.Search) (*v1beta1.Search, error) {
	result := v1beta1.Search{}
	err := c.restClient.
		Put().
		Resource("Searchs").
		Name(o.Name).
		Body(o).
		Do(ctx).
		Into(&result)

	return &result, err
}
