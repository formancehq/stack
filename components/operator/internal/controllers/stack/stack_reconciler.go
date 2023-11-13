package stack

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/formancehq/operator/internal/controllerutils"
	"github.com/formancehq/operator/internal/modules"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	_ "github.com/formancehq/operator/internal/modules/auth"
	_ "github.com/formancehq/operator/internal/modules/control"
	_ "github.com/formancehq/operator/internal/modules/gateway"
	_ "github.com/formancehq/operator/internal/modules/ledger"
	_ "github.com/formancehq/operator/internal/modules/orchestration"
	_ "github.com/formancehq/operator/internal/modules/payments"
	_ "github.com/formancehq/operator/internal/modules/search"
	_ "github.com/formancehq/operator/internal/modules/stargate"
	_ "github.com/formancehq/operator/internal/modules/wallets"
	_ "github.com/formancehq/operator/internal/modules/webhooks"
	pkgError "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	DefaultVersions = "default"
)

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions,verbs=get;list;watch

// Reconciler reconciles a Stack object
type Reconciler struct {
	client                 client.Client
	scheme                 *runtime.Scheme
	stackReconcilerFactory *modules.StackReconcilerFactory

	enableStackFinalizer bool
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx, "stack", req.NamespacedName)
	stack := &stackv1beta3.Stack{}
	if err := r.client.Get(ctx, req.NamespacedName, stack); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Object not found, skip")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, pkgError.Wrap(err, "Reading target")
	}
	log.Info("Starting reconciliation")

	stack.SetProgressing()

	conf := &stackv1beta3.Configuration{}
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: "",
		Name:      stack.Spec.Seed,
	}, conf); err != nil {
		return ctrl.Result{}, err
	}

	platform := r.stackReconcilerFactory.Platform()

	configuration := &modules.ReconciliationConfig{
		Stack:         stack,
		Configuration: conf,
		Versions:      nil,
		Platform:      platform,
	}

	stackFinalizer := NewStackFinalizer(
		r.client,
		log,
		configuration,
	)

	if r.enableStackFinalizer {
		deleted, err := stackFinalizer.HandleFinalizer(ctx, req.Name)
		if err != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, err
		}
		if deleted {
			return ctrl.Result{}, nil
		}
	} else {
		if err := stackFinalizer.RemoveFinalizer(ctx); err != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, err
		}
	}

	var (
		reconcileError error
		ready          bool
	)
	func() {
		ready, reconcileError = r.reconcileStack(ctx, stack)
		if reconcileError != nil {
			log.Info("reconciliation terminated with error", "error", reconcileError)
			stack.SetError(reconcileError)
		} else {
			log.Info("reconciliation terminated with success")
			if ready {
				stack.SetReady()
			}
		}
	}()

	if reconcileError != nil {
		log.Info("reconcile failed with error", "error", reconcileError)
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, nil
	}

	if patchErr := r.client.Status().Update(ctx, stack); patchErr != nil {
		log.Info("unable to update status", "error", patchErr)
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, nil
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) reconcileStack(ctx context.Context, stack *stackv1beta3.Stack) (bool, error) {

	configuration := &stackv1beta3.Configuration{}
	if err := r.client.Get(ctx, types.NamespacedName{
		Name: stack.Spec.Seed,
	}, configuration); err != nil {
		if errors.IsNotFound(err) {
			return false, pkgError.New("Configuration object not found")
		}
		return false, fmt.Errorf("error retrieving Configuration object: %s", err)
	}

	if err := configuration.Validate(); err != nil {
		return false, err
	}

	versionsString := stack.Spec.Versions
	if versionsString == "" {
		versionsString = DefaultVersions
	}

	versions := &stackv1beta3.Versions{}
	if err := r.client.Get(ctx, types.NamespacedName{Name: versionsString}, versions); err != nil {
		if errors.IsNotFound(err) {
			return false, pkgError.New("Versions object not found")
		}
		return false, fmt.Errorf("error retrieving Versions object: %s", err)
	}

	_, _, err := controllerutils.CreateOrUpdate(ctx, r.client, types.NamespacedName{
		Name: stack.Name,
	}, controllerutils.WithController[*corev1.Namespace](stack, r.scheme), func(ns *corev1.Namespace) {})
	if err != nil {
		return false, err
	}

	return r.stackReconcilerFactory.
		NewDeployer(stack, configuration, versions).
		Reconcile(ctx)
}

func listStacksAndReconcile(mgr ctrl.Manager, opts ...client.ListOption) []reconcile.Request {
	stacks := &stackv1beta3.StackList{}
	err := mgr.GetClient().List(context.TODO(), stacks, opts...)
	if err != nil {
		panic(err)
	}

	return collectionutils.Map(stacks.Items, func(s stackv1beta3.Stack) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      s.GetName(),
				Namespace: s.GetNamespace(),
			},
		}
	})
}

func watch(mgr ctrl.Manager, field string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		return listStacksAndReconcile(mgr, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(field, object.GetName()),
			Namespace:     object.GetNamespace(),
		})
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().
		IndexField(context.Background(), &stackv1beta3.Stack{}, ".spec.seed", func(rawObj client.Object) []string {
			return []string{rawObj.(*stackv1beta3.Stack).Spec.Seed}
		}); err != nil {
		return err
	}
	if err := mgr.GetFieldIndexer().
		IndexField(context.Background(), &stackv1beta3.Stack{}, ".spec.versions", func(rawObj client.Object) []string {
			return []string{rawObj.(*stackv1beta3.Stack).Spec.Versions}
		}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1beta3.Stack{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&networkingv1.Ingress{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Secret{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.ConfigMap{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.CronJob{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.Job{}).
		Owns(&stackv1beta3.Migration{}).
		Owns(&appsv1.Deployment{}).
		Watches(
			&source.Kind{Type: &stackv1beta3.Configuration{}},
			watch(mgr, ".spec.seed"),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&source.Kind{Type: &stackv1beta3.Versions{}},
			watch(mgr, ".spec.versions"),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&source.Kind{Type: &corev1.Secret{}},
			handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
				labels := object.GetLabels()
				partOfConfiguration, ok := labels[modules.PartOfConfigurationLabel]
				if !ok {
					return nil
				}

				var listOptions []client.ListOption
				if partOfConfiguration != modules.PartOfConfigurationAnyValue {
					listOptions = append(listOptions, &client.ListOptions{
						FieldSelector: fields.OneTermEqualSelector(".spec.seed", partOfConfiguration),
					})
				}

				return listStacksAndReconcile(mgr, listOptions...)
			}),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

type ReconcilerOpts func(*Reconciler)

func NewReconciler(client client.Client, scheme *runtime.Scheme, stackReconcilerFactory *modules.StackReconcilerFactory, opts ...ReconcilerOpts) *Reconciler {
	r := &Reconciler{
		client:                 client,
		scheme:                 scheme,
		stackReconcilerFactory: stackReconcilerFactory,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithEnableStackFinalizer(enable bool) ReconcilerOpts {
	return func(r *Reconciler) {
		r.enableStackFinalizer = enable
	}
}
