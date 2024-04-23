package stacks

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/operator/internal/core"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=versions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=versions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=versions/finalizers,verbs=update

func areDependentReady(ctx Context, stack *v1beta1.Stack) error {
	pendingResources := make([]string, 0)
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		var r v1beta1.Dependent
		r, ok := v.(v1beta1.Dependent)
		if !ok {
			continue
		}

		gvk, err := apiutil.GVKForObject(r, ctx.GetScheme())
		if err != nil {
			return err
		}
		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
			"stack": stack.Name,
		}); err != nil {

			return err
		}

		for _, item := range l.Items {
			content := item.UnstructuredContent()
			if content["status"] == nil {
				pendingResources = append(pendingResources, fmt.Sprintf("%s: %s", item.GetObjectKind().GroupVersionKind().Kind, item.GetName()))
				continue
			}

			status := content["status"].(map[string]interface{})
			if status["ready"] == nil {
				pendingResources = append(pendingResources, fmt.Sprintf("%s: %s", item.GetObjectKind().GroupVersionKind().Kind, item.GetName()))
				continue
			}

			isReady := status["ready"].(bool)
			if !isReady {
				pendingResources = append(pendingResources, fmt.Sprintf("%s: %s", item.GetObjectKind().GroupVersionKind().Kind, item.GetName()))
				continue
			}
		}

	}

	if len(pendingResources) > 0 {
		return NewApplicationError("Still pending dependent: %s ", strings.Join(pendingResources, ","))
	}

	return nil
}

func retrieveReferenceModules(ctx Context, stack *v1beta1.Stack) error {
	setKind := map[string]interface{}{}
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		r, ok := v.(v1beta1.Module)
		if !ok {
			continue
		}

		gvk, err := apiutil.GVKForObject(r, ctx.GetScheme())
		if err != nil {
			return err
		}
		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
			"stack": stack.Name,
		}); err != nil {
			return err
		}

		for _, item := range l.Items {
			content := item.UnstructuredContent()
			if content["kind"] != nil {
				kind := content["kind"].(string)
				if setKind[kind] == nil {
					setKind[kind] = []string{}
				}
			}
		}

	}

	modules := []string{}
	for k := range setKind {
		modules = append(modules, k)
	}

	sort.Strings(modules)

	stack.Status.Modules = modules

	return nil
}

func Reconcile(ctx Context, stack *v1beta1.Stack) error {
	errAlreadyExist := errors.New("namespace already exists")
	if _, _, err := CreateOrUpdate(ctx, types.NamespacedName{
		Name: stack.Name,
	},
		func(ns *corev1.Namespace) error {
			_, stackCreatedByAgent := stack.GetLabels()[v1beta1.CreatedByAgentLabel]
			if ns.ResourceVersion == "" || stackCreatedByAgent {
				return nil
			}

			return errAlreadyExist
		}, core.WithController[*corev1.Namespace](ctx.GetScheme(), stack)); err != nil {
		if !errors.Is(err, errAlreadyExist) {
			return err
		}
	}

	if err := retrieveReferenceModules(ctx, stack); err != nil {
		return err
	}

	if err := areDependentReady(ctx, stack); err != nil {
		return err
	}

	return nil
}

func Clean(ctx Context, t *v1beta1.Stack) error {
	logger := log.FromContext(ctx)
	logger = logger.WithValues("stack", t.Name)
	logger.Info("Clean stack")

	if err := deleteModules(ctx, t, logger); err != nil {
		return err
	}

	logger.Info("All modules removed")

	if err := deleteResources(ctx, t, logger); err != nil {
		return err
	}

	logger.Info("All dependencies removed")

	return nil
}

func deleteModules(ctx Context, stack *v1beta1.Stack, logger logr.Logger) error {
	stillExistingModules := false
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		module, ok := v.(v1beta1.Module)
		if !ok {
			continue
		}

		gvk, err := apiutil.GVKForObject(module, ctx.GetScheme())
		if err != nil {
			return err
		}

		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetAPIReader().List(ctx, l); err != nil {
			return err
		}

		items := collectionutils.Filter(l.Items, func(u unstructured.Unstructured) bool {
			return u.Object["spec"].(map[string]any)["stack"].(string) == stack.Name
		})

		stillExistingModules = stillExistingModules || len(items) > 0

		for _, item := range items {
			if item.GetDeletionTimestamp().IsZero() {
				logger.Info(fmt.Sprintf("Delete module %s [%s]", item.GetName(), gvk))
				if err := ctx.GetClient().Delete(ctx, &item); client.IgnoreNotFound(err) != nil {
					return err
				}
			}
		}
	}

	if stillExistingModules {
		logger.Info("Still pending modules")
		return NewPendingError()
	}

	return nil
}

func deleteResources(ctx Context, stack *v1beta1.Stack, logger logr.Logger) error {
	stillExistingResources := false
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		resource, ok := v.(v1beta1.Resource)
		if !ok {
			continue
		}
		gvk, err := apiutil.GVKForObject(resource, ctx.GetScheme())
		if err != nil {
			return err
		}

		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetAPIReader().List(ctx, l); err != nil {
			return err
		}

		items := collectionutils.Filter(l.Items, func(u unstructured.Unstructured) bool {
			return u.Object["spec"].(map[string]any)["stack"].(string) == stack.Name
		})

		stillExistingResources = stillExistingResources || len(items) > 0

		for _, item := range items {
			if item.GetDeletionTimestamp().IsZero() {
				logger.Info(fmt.Sprintf("Delete resource %s [%s]", item.GetName(), gvk))
				if err := ctx.GetClient().Delete(ctx, &item); client.IgnoreNotFound(err) != nil {
					return err
				}
			}
		}
	}

	if stillExistingResources {
		logger.Info("Still pending resources")
		return NewPendingError()
	}

	return nil
}

func init() {
	Init(
		WithSimpleIndex[*v1beta1.Stack](".spec.versionsFromFile", func(t *v1beta1.Stack) string {
			return t.Spec.VersionsFromFile
		}),
		WithStdReconciler(Reconcile,
			WithOwn[*v1beta1.Stack](&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})),
			WithRaw[*v1beta1.Stack](func(ctx Context, b *builder.Builder) error {
				for _, rtype := range ctx.GetScheme().AllKnownTypes() {
					v := reflect.New(rtype).Interface()

					switch v := v.(type) {
					case v1beta1.Module:
						u := &unstructured.Unstructured{}
						gvk, err := apiutil.GVKForObject(v, ctx.GetScheme())
						if err != nil {
							return err
						}
						u.SetGroupVersionKind(gvk)

						b.Watches(u, handler.EnqueueRequestsFromMapFunc(func(watchContext context.Context, object client.Object) []reconcile.Request {
							return []reconcile.Request{{
								NamespacedName: types.NamespacedName{
									Name: object.(*unstructured.Unstructured).Object["spec"].(map[string]any)["stack"].(string),
								},
							}}
						}))
					case v1beta1.Resource:
						u := &unstructured.Unstructured{}
						gvk, err := apiutil.GVKForObject(v, ctx.GetScheme())
						if err != nil {
							return err
						}
						u.SetGroupVersionKind(gvk)

						b.Watches(u, handler.EnqueueRequestsFromMapFunc(func(watchContext context.Context, object client.Object) []reconcile.Request {
							return []reconcile.Request{{
								NamespacedName: types.NamespacedName{
									Name: object.(*unstructured.Unstructured).Object["spec"].(map[string]any)["stack"].(string),
								},
							}}
						}))
					}
				}

				return nil
			}),
			// notes(gfyrag): Some resources need to be properly dropped before the stack is dropped
			WithFinalizer[*v1beta1.Stack]("delete", Clean),
		),
	)
}
