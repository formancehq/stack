package core

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Setup(mgr ctrl.Manager, platform Platform) error {

	if err := indexStackDependentsObjects(mgr); err != nil {
		return err
	}

	if err := setupConfigurationObjects(mgr); err != nil {
		return err
	}

	wrappedMgr := NewDefaultManager(mgr, platform)
	for _, initializer := range initializers {
		if err := initializer(wrappedMgr); err != nil {
			return err
		}
	}

	return nil
}

// indexStackDependentsObjects automatically add an index on `stack` property for all stack dependents objects
func indexStackDependentsObjects(mgr ctrl.Manager) error {
	for _, rtype := range mgr.GetScheme().AllKnownTypes() {

		object, ok := reflect.New(rtype).Interface().(client.Object)
		if !ok {
			continue
		}

		_, ok = object.(v1beta1.Dependent)
		if !ok {
			continue
		}

		mgr.GetLogger().Info("Detect stack dependency object, automatically index field", "type", rtype)
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), object, "stack", func(object client.Object) []string {
				return []string{object.(v1beta1.Dependent).GetStack()}
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", rtype)
			return err
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(object)
		if err != nil {
			return err
		}
		us := &unstructured.Unstructured{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), us, "stack", func(object client.Object) []string {
				stack := object.(*unstructured.Unstructured).Object["spec"].(map[string]any)["stack"]
				if stack == nil {
					return []string{}
				}
				return []string{stack.(string)}
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", &unstructured.Unstructured{})
			return err
		}
	}
	return nil
}

func setupConfigurationObjects(mgr ctrl.Manager) error {
	for _, rtype := range mgr.GetScheme().AllKnownTypes() {

		object, ok := reflect.New(rtype).Interface().(client.Object)
		if !ok {
			continue
		}

		_, ok = object.(v1beta1.ConfigurationObject)
		if !ok {
			continue
		}

		mgr.GetLogger().Info("Detect configuration object, automatically index field", "type", rtype)
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), object, "stack", func(object client.Object) []string {
				return object.(v1beta1.ConfigurationObject).GetStacks()
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", rtype)
			return err
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(object)
		if err != nil {
			return err
		}
		us := &unstructured.Unstructured{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), us, "stack", func(object client.Object) []string {
				u := object.(*unstructured.Unstructured)
				spec := u.UnstructuredContent()["spec"].(map[string]any)
				if stacks, ok := spec["stacks"]; !ok {
					return []string{}
				} else {
					return collectionutils.Map(stacks.([]any), func(v any) string {
						return v.(string)
					})
				}
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", &unstructured.Unstructured{})
			return err
		}
	}
	return nil
}
