package testing

import (
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller"
	"github.com/formancehq/operator/v2/internal/controller/shared"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"os"
	"path/filepath"
	osRuntime "runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	controllerruntime "sigs.k8s.io/controller-runtime/pkg/controller"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//+kubebuilder:scaffold:imports
)

var (
	ctx        context.Context
	cancel     func()
	testEnv    *envtest.Environment
	restConfig *rest.Config
	k8sClient  client.Client
)

func GetScheme() *runtime.Scheme {
	return scheme.Scheme
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(os.Stdout), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.Background())

	_, filename, _, _ := osRuntime.Caller(0)

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join(filepath.Dir(filename), "..", "..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	restConfig, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(restConfig).NotTo(BeNil())

	Expect(v1beta1.AddToScheme(scheme.Scheme)).To(Succeed())

	k8sClient, err = client.New(restConfig, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	Expect(testEnv.Stop())
})

var (
	done chan struct{}
)

var _ = BeforeEach(func() {
	ctx, cancel = context.WithCancel(context.Background())
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme: GetScheme(),
		Metrics: server.Options{
			BindAddress: "0",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	SetDefaultEventuallyTimeout(5 * time.Second)

	platform := shared.Platform{
		Region:      "us-west-1",
		Environment: "staging",
	}

	Expect(reconcilers.SetupReconcilers(mgr,
		controller.NewStackReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewTopicQueryReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewTopicReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewLedgerReconciler(mgr.GetClient(), mgr.GetScheme()),
		reconcilers.New[*v1beta1.HTTPAPI](mgr.GetClient(), mgr.GetScheme(),
			controller.ForHTTPAPI(mgr.GetClient(), mgr.GetScheme())),
		controller.NewGatewayReconciler(mgr.GetClient(), mgr.GetScheme(), platform),
		controller.NewAuthReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewDatabaseReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewAuthClientReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewWalletReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewOrchestrationReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewPaymentsReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewReconciliationReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewWebhooksReconciler(mgr.GetClient(), mgr.GetScheme()),
		controller.NewSearchReconciler(mgr.GetClient(), mgr.GetScheme()),
	)).To(Succeed())

	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithOptions(controllerruntime.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

			deployment := &appsv1.Deployment{}
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			}, deployment); err != nil {
				return reconcile.Result{}, err
			}
			deployment.Status.ObservedGeneration = deployment.Generation
			deployment.Status.UpdatedReplicas = 1
			deployment.Status.AvailableReplicas = 1
			deployment.Status.Replicas = 1
			deployment.Status.ReadyReplicas = 1
			if err := mgr.GetClient().Status().Update(ctx, deployment); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}))
	Expect(err)

	err = ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		WithOptions(controllerruntime.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

			job := &batchv1.Job{}
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			}, job); err != nil {
				return reconcile.Result{}, err
			}
			job.Status.Succeeded = 1

			if err := mgr.GetClient().Status().Update(ctx, job); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}))
	Expect(err)

	go func() {
		defer GinkgoRecover()
		done = make(chan struct{})
		err := mgr.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
		close(done)
	}()
})

var _ = AfterEach(func() {
	cancel()
	if done != nil {
		<-done
	}
})

func Create(ob client.Object) error {
	return k8sClient.Create(ctx, ob)
}

func Delete(ob client.Object) error {
	return k8sClient.Delete(ctx, ob)
}

func Update(ob client.Object) error {
	return k8sClient.Update(ctx, ob)
}

func Patch(ob client.Object, patch client.Patch) error {
	return k8sClient.Patch(ctx, ob, patch)
}

func Get(key types.NamespacedName, ob client.Object) error {
	return k8sClient.Get(ctx, key, ob)
}

func Client() client.Client {
	return k8sClient
}
