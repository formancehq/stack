package stack_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	osRuntime "runtime"
	"testing"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack"
	"github.com/formancehq/operator/internal/modules"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//+kubebuilder:scaffold:imports
)

func TestAPIs(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Controller Suite")
}

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

var _ = ginkgo.BeforeSuite(func() {
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
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(restConfig).NotTo(gomega.BeNil())

	gomega.Expect(v1beta3.AddToScheme(scheme.Scheme)).To(gomega.Succeed())

	k8sClient, err = client.New(restConfig, client.Options{Scheme: scheme.Scheme})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(k8sClient).NotTo(gomega.BeNil())
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(testEnv.Stop())
})

var (
	done chan struct{}
)

var _ = ginkgo.BeforeEach(func() {
	ctx, cancel = context.WithCancel(context.Background())
	done = make(chan struct{})
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:             GetScheme(),
		MetricsBindAddress: "0",
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	config := stack.Configuration{
		Region:      "us-west-1",
		Environment: "staging",
	}
	err = stack.NewReconciler(mgr.GetClient(), mgr.GetScheme(), modules.NewStackDeployer(http.DefaultTransport), config).SetupWithManager(mgr)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1beta3.Migration{}).
		Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
			log := log.FromContext(ctx, "migration", request.NamespacedName)
			log.Info("Starting reconciliation")

			migration := &v1beta3.Migration{}
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			}, migration); err != nil {
				return reconcile.Result{}, err
			}
			migration.Status.Terminated = true
			if err := mgr.GetClient().Status().Update(ctx, migration); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}))
	gomega.Expect(err)

	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
			log := log.FromContext(ctx, "deployment", request.NamespacedName)
			log.Info("Starting reconciliation")

			deployment := &v1.Deployment{}
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			}, deployment); err != nil {
				return reconcile.Result{}, err
			}
			deployment.Status.ObservedGeneration = deployment.Generation
			if len(deployment.Status.Conditions) == 0 {
				deployment.Status.Conditions = append(deployment.Status.Conditions, v1.DeploymentCondition{})
			}
			deployment.Status.Conditions[0] = v1.DeploymentCondition{
				Type:               v1.DeploymentAvailable,
				Status:             "True",
				LastUpdateTime:     metav1.Time{Time: time.Now()},
				LastTransitionTime: metav1.Time{Time: time.Now()},
			}
			if err := mgr.GetClient().Status().Update(ctx, deployment); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}))
	gomega.Expect(err)

	go func() {
		defer ginkgo.GinkgoRecover()
		err := mgr.Start(ctx)
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "failed to run manager")
		close(done)
	}()
})

var _ = ginkgo.AfterEach(func() {
	cancel()
	<-done
})

func Create(ob client.Object) error {
	return k8sClient.Create(ctx, ob)
}

func Delete(ob client.Object) error {
	return k8sClient.Delete(ctx, ob)
}

func Get(key types.NamespacedName, ob client.Object) error {
	return k8sClient.Get(ctx, key, ob)
}
