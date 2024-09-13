package tests_test

import (
	"fmt"

	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	core "github.com/formancehq/operator/internal/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrchestrationController", func() {
	Context("When creating a Orchestration object", func() {
		var (
			stack               *v1beta1.Stack
			gateway             *v1beta1.Gateway
			auth                *v1beta1.Auth
			ledger              *v1beta1.Ledger
			orchestration       *v1beta1.Orchestration
			databaseSettings    *v1beta1.Settings
			brokerDSNSettings   *v1beta1.Settings
			temporalDSNSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
			brokerDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
			temporalDSNSettings = settings.New(uuid.NewString(), "temporal.dsn", "temporal://localhost/namespace", stack.Name)
			gateway = &v1beta1.Gateway{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewaySpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Ingress: &v1beta1.GatewayIngress{},
				},
			}
			auth = &v1beta1.Auth{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.AuthSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			ledger = &v1beta1.Ledger{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.LedgerSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			orchestration = &v1beta1.Orchestration{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.OrchestrationSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(databaseSettings)).To(Succeed())
			Expect(Create(temporalDSNSettings)).To(Succeed())
			Expect(Create(brokerDSNSettings)).To(BeNil())
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
			Expect(Create(orchestration)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(brokerDSNSettings)).To(Succeed())
			Expect(Delete(temporalDSNSettings)).To(Succeed())
		})
		It("Should create appropriate components", func() {
			By("Should set the status to ready", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", orchestration.Name, orchestration)).To(Succeed())
					return orchestration.Status.Ready
				}).Should(BeTrue())
			})
			By("Should add an owner reference on the stack", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", orchestration.Name, orchestration)).To(Succeed())
					reference, err := core.HasOwnerReference(TestContext(), stack, orchestration)
					g.Expect(err).To(BeNil())
					return reference
				}).Should(BeTrue())
			})
			By("Should create a deployment", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "orchestration", deployment)
				}).Should(Succeed())
				Expect(deployment).To(BeControlledBy(orchestration))
				Expect(deployment.Spec.Template.Spec.Containers[0].Env).To(ContainElements(
					core.Env("WORKER", "true"),
					core.Env("TOPICS", fmt.Sprintf("%s.ledger", stack.Name)),
				))
			})
			By("Should create a new GatewayHTTPAPI object", func() {
				httpService := &v1beta1.GatewayHTTPAPI{}
				Eventually(func() error {
					return LoadResource("", core.GetObjectName(stack.Name, "orchestration"), httpService)
				}).Should(Succeed())
			})
			By("Should create a new AuthClient object", func() {
				authClient := &v1beta1.AuthClient{}
				Eventually(func() error {
					return LoadResource("", core.GetObjectName(stack.Name, "orchestration"), authClient)
				}).Should(Succeed())
			})
			By("Should create a new BrokerConsumer object", func() {
				consumer := &v1beta1.BrokerConsumer{}
				Eventually(func() error {
					return LoadResource("", orchestration.Name+"-orchestration", consumer)
				}).Should(Succeed())
			})
		})
	})
})
