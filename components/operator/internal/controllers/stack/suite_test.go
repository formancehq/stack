package stack_test

import (
	"context"
	"os"
	"path/filepath"
	osRuntime "runtime"
	"testing"

	v1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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

	gomega.Expect(pgtesting.CreatePostgresServer()).To(gomega.BeNil())

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
	gomega.Expect(pgtesting.DestroyPostgresServer()).To(gomega.BeNil())
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

	err = stack.NewReconciler(mgr.GetClient(), mgr.GetScheme(), "us-west-1", "staging").SetupWithManager(mgr)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

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

func NewDumbVersions() *v1beta3.Versions {
	return &v1beta3.Versions{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: v1beta3.VersionsSpec{
			Control:       "latest",
			Ledger:        "latest",
			Payments:      "latest",
			Search:        "latest",
			Auth:          "latest",
			Webhooks:      "latest",
			Wallets:       "latest",
			Stargate:      "latest",
			Orchestration: "latest",
		},
	}
}

func NewDumbConfiguration() *v1beta3.Configuration {
	return &v1beta3.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: v1beta3.ConfigurationSpec{
			Services: v1beta3.ConfigurationServicesSpec{
				Auth: v1beta3.AuthSpec{
					Postgres: NewPostgresConfig(),
				},
				Control: v1beta3.ControlSpec{},
				Ledger: v1beta3.LedgerSpec{
					Postgres: NewPostgresConfig(),
				},
				Payments: v1beta3.PaymentsSpec{
					Postgres: NewPostgresConfig(),
				},
				Search: v1beta3.SearchSpec{
					ElasticSearchConfig: NewDumpElasticSearchConfig(),
				},
				Webhooks: v1beta3.WebhooksSpec{
					Postgres: NewPostgresConfig(),
				},
				Stargate: v1beta3.StargateSpec{},
				Wallets:  v1beta3.WalletsSpec{},
				Orchestration: v1beta3.OrchestrationSpec{
					Postgres: NewPostgresConfig(),
				},
			},
			Broker:     NewDumbBrokerConfig(),
			Monitoring: NewDumbMonitoring(),
		},
	}
}

func NewDumpKafkaConfig() v1beta3.KafkaConfig {
	return v1beta3.KafkaConfig{
		Brokers: []string{"kafka:1234"},
	}
}

func NewDumbBrokerConfig() v1beta3.Broker {
	return v1beta3.Broker{
		Kafka: func() *v1beta3.KafkaConfig {
			ret := NewDumpKafkaConfig()
			return &ret
		}(),
	}
}

func NewDumbMonitoring() *v1beta3.MonitoringSpec {
	return &v1beta3.MonitoringSpec{
		Traces: &v1beta3.TracesSpec{
			Otlp: &v1beta3.TracesOtlpSpec{
				Endpoint: "localhost",
				Port:     4317,
				Insecure: true,
				Mode:     "grpc",
			},
		},
	}
}

func NewDumpElasticSearchConfig() v1beta3.ElasticSearchConfig {
	return v1beta3.ElasticSearchConfig{
		Scheme: "http",
		Host:   "elasticsearch",
		Port:   9200,
	}
}

func NewPostgresConfig() v1beta3.PostgresConfig {
	return v1beta3.PostgresConfig{
		Port:           pgtesting.Server().GetPort(),
		Host:           pgtesting.Server().GetHost(),
		Username:       pgtesting.Server().GetUsername(),
		Password:       pgtesting.Server().GetPassword(),
		DisableSSLMode: true,
	}
}
