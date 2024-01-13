package stack_formance_com_test

import (
	"fmt"
	"math/rand"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/controller/testing"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("StackController (legacy)", func() {
	When("Creating a legacy stack with a Configuration and a Versions", func() {
		var (
			configuration *v1beta3.Configuration
			versions      *v1beta3.Versions
			stack         *v1beta3.Stack
		)
		BeforeEach(func() {
			configuration = &v1beta3.Configuration{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta3.ConfigurationSpec{
					Services: v1beta3.ConfigurationServicesSpec{
						Reconciliation: v1beta3.ReconciliationSpec{
							CommonServiceProperties: v1beta3.CommonServiceProperties{
								Disabled: pointer.For(false),
							},
						},
					},
					Monitoring: &v1beta3.MonitoringSpec{
						Traces: &v1beta3.TracesSpec{
							Otlp: &v1beta3.OtlpSpec{
								Endpoint:           "collector",
								Port:               4317,
								Insecure:           false,
								Mode:               "grpc",
								ResourceAttributes: "foo=bar",
							},
						},
					},
					Registries: map[string]v1beta1.RegistryConfigurationSpec{
						"ghcr.io": {
							Endpoint: "http://localhost:8080",
						},
					},
				},
			}
			Expect(Create(configuration)).To(Succeed())
			versions = &v1beta3.Versions{
				ObjectMeta: RandObjectMeta(),
			}
			Expect(Create(versions)).To(Succeed())
			stack = &v1beta3.Stack{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%d-%d", rand.Int31(), rand.Int31()),
				},
				Spec: v1beta3.StackSpec{
					Seed:     configuration.Name,
					Versions: versions.Name,
					Auth: v1beta3.StackAuthSpec{
						StaticClients: []*v1beta1.AuthClientSpec{{
							ID:     "client0",
							Secret: "client0",
						}},
					},
					Stargate: &v1beta3.StackStargateConfig{},
				},
			}
			Expect(Create(stack)).To(Succeed())
		})
		It("Should create the DatabaseConfiguration objects", func() {
			Eventually(func(g Gomega) int {
				list := &v1beta1.DatabaseConfigurationList{}
				g.Expect(List(list)).To(Succeed())
				return len(list.Items)
			}).Should(BeNumerically("==", 5))
		})
		It("Should create the OpenTelemetryConfiguration object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.OpenTelemetryConfiguration{})
			}).Should(Succeed())
		})
		It("Should create the BrokerConfiguration object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.BrokerConfiguration{})
			}).Should(Succeed())
		})
		It("Should create the ElasticSearchConfiguration object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.ElasticSearchConfiguration{})
			}).Should(Succeed())
		})
		It("Should create the Ledger object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Ledger{})
			}).Should(Succeed())
		})
		It("Should create the Auth object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Auth{})
			}).Should(Succeed())
		})
		It("Should create the Gateway object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Gateway{})
			}).Should(Succeed())
		})
		It("Should create the Orchestration object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Orchestration{})
			}).Should(Succeed())
		})
		It("Should create the Payments object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Payments{})
			}).Should(Succeed())
		})
		It("Should create the Search object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Search{})
			}).Should(Succeed())
		})
		It("Should create the Stargate object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Stargate{})
			}).Should(Succeed())
		})
		It("Should create the Wallets object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Wallets{})
			}).Should(Succeed())
		})
		It("Should create the Webhooks object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Webhooks{})
			}).Should(Succeed())
		})
		It("Should create the Reconciliation object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.Reconciliation{})
			}).Should(Succeed())
		})
		It("Should create a AuthClient object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", fmt.Sprintf("%s-client0", stack.Name), &v1beta1.AuthClient{})
			}).Should(Succeed())
		})
		It("Should create a RegistriesConfiguration object", func() {
			Eventually(func(g Gomega) error {
				return LoadResource("", stack.Name, &v1beta1.RegistriesConfiguration{})
			}).Should(Succeed())
		})
	})
})
