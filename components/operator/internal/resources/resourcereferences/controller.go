package resourcereferences

import (
	"fmt"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//+kubebuilder:rbac:groups=formance.com,resources=resourcereferences,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=resourcereferences/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=resourcereferences/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

func init() {
	Init(
		WithStackDependencyReconciler[*v1beta1.ResourceReference](Reconcile,
			WithWatch[*v1beta1.ResourceReference, *v1.Secret](watchResource[*v1.Secret]),
			WithWatch[*v1beta1.ResourceReference, *v1.ServiceAccount](watchResource[*v1.ServiceAccount]),
		),
	)
}

func watchResource[T client.Object](ctx Context, object T) []reconcile.Request {
	ret := make([]reconcile.Request, 0)

	// Watch resources created by the ResourceReference
	var resourceReference string
	for _, reference := range object.GetOwnerReferences() {
		gvk, err := apiutil.GVKForObject(&v1beta1.ResourceReference{}, ctx.GetScheme())
		if err != nil {
			panic(err)
		}
		apiVersion, kind := gvk.ToAPIVersionAndKind()
		if reference.Kind == kind && reference.APIVersion == apiVersion {
			resourceReference = reference.Name
			break
		}
	}

	if resourceReference != "" {
		return append(ret, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: resourceReference,
			},
		})
	}

	// Watch resources which should be replicated by the ResourceReferences
	if object.GetLabels()[v1beta1.StackLabel] != "any" {
		for _, stack := range strings.Split(object.GetLabels()[v1beta1.StackLabel], ",") {
			ret = append(ret, BuildReconcileRequests(
				ctx,
				ctx.GetClient(),
				ctx.GetScheme(),
				&v1beta1.ResourceReference{},
				client.MatchingFields{
					"stack": stack,
				},
			)...)
		}
	} else {
		ret = append(ret, BuildReconcileRequests(
			ctx,
			ctx.GetClient(),
			ctx.GetScheme(),
			&v1beta1.ResourceReference{},
		)...)
	}

	return ret
}

const (
	AnyValue = "any"

	RewrittenResourceName = "formance.com/referenced-by-name"
)

func Reconcile(ctx Context, stack *v1beta1.Stack, req *v1beta1.ResourceReference) error {
	resource, err := findMatchingResource(ctx, stack.Name, *req.Spec.GroupVersionKind, req.Spec.Name)
	if err != nil {
		return err
	}

	gvk := schema.GroupVersionKind{
		Group:   req.Spec.GroupVersionKind.Group,
		Version: req.Spec.GroupVersionKind.Version,
		Kind:    req.Spec.GroupVersionKind.Kind,
	}

	if req.Status.SyncedResource != "" && req.Spec.Name != req.Status.SyncedResource {
		oldResource := &unstructured.Unstructured{}
		oldResource.SetGroupVersionKind(gvk)
		err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Namespace: stack.Name,
			Name:      req.Status.SyncedResource,
		}, oldResource)
		if client.IgnoreNotFound(err) != nil {
			return err
		}

		if err == nil { // Can be not found, if the resource has been manually deleted
			patch := client.MergeFrom(oldResource.DeepCopy())
			if err := controllerutil.RemoveOwnerReference(req, oldResource, ctx.GetScheme()); err == nil {
				if err := ctx.GetClient().Patch(ctx, oldResource, patch); err != nil {
					return nil
				}
			}
			if len(oldResource.GetOwnerReferences()) == 0 {
				if err := ctx.GetClient().Delete(ctx, oldResource); err != nil {
					return err
				}
			}
		}
	}

	annotations := make(map[string]any)
	originalMetadata := resource.UnstructuredContent()["metadata"]
	if originalMetadata != nil {
		metadata := originalMetadata.(map[string]any)
		originalAnnotations := metadata["annotations"]
		if originalAnnotations != nil {
			annotations = originalAnnotations.(map[string]any)
		}
	}

	unstructured.RemoveNestedField(resource.UnstructuredContent(), "metadata")
	unstructured.RemoveNestedField(resource.UnstructuredContent(), "status")

	newResource := &unstructured.Unstructured{}
	newResource.SetGroupVersionKind(gvk)
	newResource.SetNamespace(stack.Name)
	newResource.SetName(req.Spec.Name)

	_, err = controllerutil.CreateOrUpdate(ctx, ctx.GetClient(), newResource, func() error {
		content := newResource.UnstructuredContent()
		if err := mergo.Merge(&content, resource.UnstructuredContent()); err != nil {
			return err
		}

		if err := unstructured.SetNestedMap(content, annotations, "metadata", "annotations"); err != nil {
			panic(err)
		}

		hasOwnerReference, err := HasOwnerReference(ctx, req, newResource)
		if err != nil {
			return err
		}
		if !hasOwnerReference {
			if err := controllerutil.SetOwnerReference(req, newResource, ctx.GetScheme()); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	req.Status.Hash = HashFromResources(resource)
	req.Status.SyncedResource = req.Spec.Name

	return nil
}

func findMatchingResource(ctx Context, stack string, gvk metav1.GroupVersionKind, name string) (*unstructured.Unstructured, error) {
	requirement, err := labels.NewRequirement(v1beta1.StackLabel, selection.In, []string{stack, AnyValue})
	if err != nil {
		return nil, err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	})
	if err := ctx.GetClient().List(ctx, list, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return nil, errors.Wrap(err, "listing resources")
	}

	foundResources := make([]*unstructured.Unstructured, 0)
	for _, item := range list.Items {
		resourceName, ok := item.GetAnnotations()[RewrittenResourceName]
		if !ok {
			resourceName = item.GetName()
		}

		if resourceName != name {
			continue
		}
		foundResources = append(foundResources, item.DeepCopy())
	}

	if len(foundResources) > 1 {
		return nil, fmt.Errorf("found more than one matching item for '%s': %s", name, collectionutils.Map(foundResources, func(from *unstructured.Unstructured) string {
			return from.GetName()
		}))
	}
	if len(foundResources) == 0 {
		return nil, fmt.Errorf("item not found: %s", name)
	}

	return foundResources[0], nil
}
