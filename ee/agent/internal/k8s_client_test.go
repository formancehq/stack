package internal

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestEnsureNotExistBySelector(t *testing.T) {
	test(t, func(ctx context.Context, testConfig *testConfig) {
		k8sClient := NewDefaultK8SClient(testConfig.client)

		for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
			gvk := gvk

			object := reflect.New(rtype).Interface()
			if _, ok := object.(v1beta1.Module); !ok {
				continue
			}

			t.Run(fmt.Sprintf("EnsureNotExistBySelector %s", gvk.Kind), func(t *testing.T) {
				t.Parallel()
				name := uuid.NewString()
				module := unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": gvk.GroupVersion().String(),
						"kind":       gvk.Kind,
						"metadata": map[string]interface{}{
							"name": name,
							"labels": map[string]interface{}{
								"formance.com/created-by-agent": "true",
								"formance.com/stack":            name,
							},
						},
					},
				}

				resources, err := testConfig.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
				require.NoError(t, err)

				require.NoError(t, testConfig.client.Post().Resource(resources.Resource.Resource).Body(&module).Do(ctx).Error())

				require.NoError(t, k8sClient.EnsureNotExistsBySelector(ctx, resources.Resource.Resource, stackLabels(module.GetName())))
				require.NoError(t, client.IgnoreNotFound(testConfig.client.Get().Resource(resources.Resource.Resource).Name(module.GetName()).Do(ctx).Error()))
			})

		}
	})
}
