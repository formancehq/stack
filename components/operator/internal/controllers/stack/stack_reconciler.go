package stack

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/formancehq/operator/internal/controllers/stack/delete"
	"github.com/formancehq/operator/internal/controllers/stack/storage/s3"
	"github.com/formancehq/operator/internal/controllerutils"
	"github.com/formancehq/operator/internal/modules"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	pkgError "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	_ "github.com/formancehq/operator/internal/handlers"
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
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions,verbs=get;list;watch

type Configuration struct {
	// Cloud region where the stack is deployed
	Region string
	// Cloud environment where the stack is deployed: staging, production,
	// sandbox, etc.
	Environment string
}

// Reconciler reconciles a Stack object
type Reconciler struct {
	configuration Configuration
	client        client.Client
	scheme        *runtime.Scheme
	stackDeployer *modules.StackDeployer
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

	res, err := r.HandleFinalizer(ctx, log, stack, req)
	if res != nil || err != nil {
		return *res, err
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

// Neet to be able to be called multiple times !!!
// Need to be idempotent
func (r *Reconciler) deleteStack(ctx context.Context, key types.NamespacedName, stack *stackv1beta3.Stack, log logr.Logger) error {
	conf := &stackv1beta3.Configuration{}
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: "",
		Name:      stack.Spec.Seed,
	}, conf); err != nil {
		return err
	}
	s3Client, err := s3.NewClient(
		"formance",
		"formance",
		"localhost:9000",
		"toto",
		true,
		true,
	)
	if err != nil {
		log.Error(err, "Cannot create s3 client")
		return err
	}

	bucket := "backups"
	storage := s3.NewS3Storage(s3Client, bucket)

	log.Info("start backup for " + stack.Name)
	if err := delete.BackupServicesData(conf, stack, storage, log); err != nil {
		log.Error(err, "Error during backups")
	}

	log.Info("start deleting databases " + stack.Name)
	if err := delete.DeleteServiceData(conf, stack.Name, log); err != nil {
		log.Error(err, "Error during deleting databases")
	}

	log.Info("start deleting brokers subjects " + stack.Name)
	if err := delete.DeleteBrokersData(conf, stack.Name, []string{"ledger", "payments"}, log); err != nil {
		log.Error(err, "Error during deleting brokers subjects")
	}

	return nil

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
		Owns(&appsv1.Deployment{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&stackv1beta3.Migration{}).
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
		Complete(r)
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

	deployer := modules.NewDeployer(r.client, r.scheme, stack, configuration)
	resolveContext := modules.Context{
		Context:       ctx,
		Region:        r.configuration.Region,
		Environment:   r.configuration.Environment,
		Stack:         stack,
		Configuration: configuration,
		Versions:      versions,
	}

	return r.stackDeployer.HandleStack(resolveContext, deployer)
}

func NewReconciler(client client.Client, scheme *runtime.Scheme, stackDeployer *modules.StackDeployer, configuration Configuration) *Reconciler {
	return &Reconciler{
		configuration: configuration,
		client:        client,
		scheme:        scheme,
		stackDeployer: stackDeployer,
	}
}

func watch(mgr ctrl.Manager, field string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		stacks := &stackv1beta3.StackList{}
		listOps := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(field, object.GetName()),
			Namespace:     object.GetNamespace(),
		}
		err := mgr.GetClient().List(context.TODO(), stacks, listOps)
		if err != nil {
			return []reconcile.Request{}
		}

		return collectionutils.Map(stacks.Items, func(s stackv1beta3.Stack) reconcile.Request {
			return reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      s.GetName(),
					Namespace: s.GetNamespace(),
				},
			}
		})
	})
}
