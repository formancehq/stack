package tests_test

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	core "github.com/formancehq/operator/internal/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("OrchestrationController", func() {
	Context("When creating a Orchestration object", func() {
		var (
			stack                     *v1beta1.Stack
			gateway                   *v1beta1.Gateway
			auth                      *v1beta1.Auth
			ledger                    *v1beta1.Ledger
			databaseSettings          *v1beta1.Settings
			orchestration             *v1beta1.Orchestration
			brokerKindSettings        *v1beta1.Settings
			brokerNatsDSNSettings     *v1beta1.Settings
			temporalAddressSettings   *v1beta1.Settings
			temporalNamespaceSettings *v1beta1.Settings
			temporalTLSCrtSettings    *v1beta1.Settings
			temporalTLSKeySettings    *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
			brokerKindSettings = settings.New(uuid.NewString(), "broker.kind", "nats", stack.Name)
			brokerNatsDSNSettings = settings.New(uuid.NewString(), "broker.nats.dsn", "nats://localhost:1234", stack.Name)
			temporalAddressSettings = settings.New(uuid.NewString(), "temporal.address", "localhost", stack.Name)
			temporalNamespaceSettings = settings.New(uuid.NewString(), "temporal.namespace", "namespace", stack.Name)
			temporalTLSCrtSettings = settings.New(uuid.NewString(), "temporal.tls.crt", "crt", stack.Name)
			temporalTLSKeySettings = settings.New(uuid.NewString(), "temporal.tls.key", "key", stack.Name)
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
			Expect(Create(temporalAddressSettings)).To(Succeed())
			Expect(Create(temporalNamespaceSettings)).To(Succeed())
			Expect(Create(temporalTLSCrtSettings)).To(Succeed())
			Expect(Create(temporalTLSKeySettings)).To(Succeed())
			Expect(Create(brokerKindSettings)).To(BeNil())
			Expect(Create(brokerNatsDSNSettings)).To(BeNil())
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
			Expect(Create(orchestration)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(auth)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(orchestration)).To(Succeed())
			Expect(Delete(brokerNatsDSNSettings)).To(Succeed())
			Expect(Delete(brokerKindSettings)).To(Succeed())
			Expect(Delete(temporalAddressSettings)).To(Succeed())
			Expect(Delete(temporalNamespaceSettings)).To(Succeed())
			Expect(Delete(temporalTLSCrtSettings)).To(Succeed())
			Expect(Delete(temporalTLSKeySettings)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "orchestration", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(orchestration))
			Expect(deployment.Spec.Template.Spec.Containers[0].Env).To(ContainElements(
				core.Env("WORKER", "true"),
				core.Env("TOPICS", fmt.Sprintf("%s-ledger", stack.Name)),
			))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "orchestration"), httpService)
			}).Should(Succeed())
		})
		It("Should create a new AuthClient object", func() {
			authClient := &v1beta1.AuthClient{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "orchestration"), authClient)
			}).Should(Succeed())
		})
		It("Should create a new BrokerTopicConsumer object for the ledger", func() {
			authClient := &v1beta1.BrokerTopicConsumer{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "orchestration-ledger"), authClient)
			}).Should(Succeed())
		})
	})
})
