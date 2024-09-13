package internal

import (
	"context"
	"os"
	"path/filepath"
	osRuntime "runtime"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	_ "github.com/formancehq/operator/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	controllerruntime "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//+kubebuilder:scaffold:imports
)

var (
	ctx        context.Context
	cancel     func()
	testEnv    *envtest.Environment
	restConfig *rest.Config
	k8sClient  client.Client
	coreMgr    core.Manager
)

func GetScheme() *runtime.Scheme {
	return scheme.Scheme
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(os.Stdout), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.Background())

	SetDefaultEventuallyTimeout(10 * time.Second)

	_, filename, _, _ := osRuntime.Caller(0)

	apiServer := envtest.APIServer{}
	apiServer.Configure().
		Set("service-cluster-ip-range", "10.0.0.0/20")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join(filepath.Dir(filename), "..", "..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
		ControlPlane: envtest.ControlPlane{
			APIServer: &apiServer,
		},
	}

	var err error
	restConfig, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(restConfig).NotTo(BeNil())

	Expect(v1beta1.AddToScheme(scheme.Scheme)).To(Succeed())

	k8sClient, err = client.New(restConfig, client.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	ctx, cancel = context.WithCancel(context.Background())
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme: GetScheme(),
		Metrics: server.Options{
			BindAddress: "0",
		},
		Client: client.Options{
			Cache: &client.CacheOptions{
				Unstructured: true,
			},
		},
	})
	Expect(err).ToNot(HaveOccurred())
	coreMgr = core.NewDefaultManager(mgr, core.Platform{
		Region:      "testing",
		Environment: "testing",
	})

	Expect(core.Setup(mgr, core.Platform{
		Region:      "us-west-1",
		Environment: "staging",
	})).To(Succeed())

	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithOptions(controllerruntime.Options{
			MaxConcurrentReconciles: 100,
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
			MaxConcurrentReconciles: 100,
		}).
		Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

			job := &batchv1.Job{}
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			}, job); err != nil {
				if client.IgnoreNotFound(err) == nil {
					return reconcile.Result{}, nil
				}
				return reconcile.Result{}, err
			}
			if !job.DeletionTimestamp.IsZero() {
				patch := client.MergeFrom(job.DeepCopy())
				if controllerutil.RemoveFinalizer(job, "orphan") {
					if err := mgr.GetClient().Patch(ctx, job, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
				return reconcile.Result{}, nil
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

var _ = AfterSuite(func() {
	cancel()
	if done != nil {
		<-done
	}
	Expect(testEnv.Stop())
})

var (
	done chan struct{}
)

func Create(objects ...client.Object) error {
	for _, object := range objects {
		if err := k8sClient.Create(ctx, object); err != nil {
			return err
		}
	}
	return nil
}

func Delete(objects ...client.Object) error {
	for _, object := range objects {
		if err := k8sClient.Delete(ctx, object); err != nil {
			return err
		}
	}
	return nil
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

func List(list client.ObjectList, opts ...client.ListOption) error {
	return k8sClient.List(ctx, list, opts...)
}

func Client() client.Client {
	return k8sClient
}

func TestContext() core.Context {
	return core.NewContext(coreMgr, ctx)
}
