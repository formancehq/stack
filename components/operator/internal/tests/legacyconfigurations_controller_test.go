package tests_test

import (
	. "github.com/formancehq/operator/internal/tests/internal"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigurationsController (legacy)", func() {
	When("Creating a legacy Configuration", func() {
		var (
			configuration *v1beta3.Configuration
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
					Registries: map[string]v1beta3.RegistryConfig{
						"ghcr.io": {
							Endpoint: "http://localhost:8080",
						},
					},
				},
			}
			Expect(Create(configuration)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(configuration)).To(Succeed())
		})
		It("Should create the settings objects", func() {
			//By("DatabaseConfiguration", func() {
			//	for _, service := range []string{
			//		"auth", "ledger", "payments", "orchestration", "webhooks", "reconciliation",
			//	} {
			//		Eventually(func(g Gomega) error {
			//			return LoadResource("", fmt.Sprintf("%s-%s", configuration.Name, service), &v1beta1.DatabaseConfiguration{})
			//		}).Should(Succeed())
			//	}
			//})
			By("OpenTelemetryConfiguration", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", configuration.Name, &v1beta1.OpenTelemetryConfiguration{})
				}).Should(Succeed())
			})
			By("BrokerConfiguration", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", configuration.Name, &v1beta1.BrokerConfiguration{})
				}).Should(Succeed())
			})
			By("ElasticSearchConfiguration", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", configuration.Name, &v1beta1.ElasticSearchConfiguration{})
				}).Should(Succeed())
			})
			By("RegistriesConfiguration", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", configuration.Name, &v1beta1.RegistriesConfiguration{})
				}).Should(Succeed())
			})
			By("SearchBatchingConfiguration", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", configuration.Name, &v1beta1.SearchBatchingConfiguration{})
				}).Should(Succeed())
			})
		})
	})
})
