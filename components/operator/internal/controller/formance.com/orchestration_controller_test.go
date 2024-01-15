package formance_com_test

import (
	"fmt"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/controller/testing"
	core "github.com/formancehq/operator/internal/core"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("OrchestrationController", func() {
	Context("When creating a Auth object", func() {
		var (
			stack                 *v1beta1.Stack
			gateway               *v1beta1.Gateway
			auth                  *v1beta1.Auth
			ledger                *v1beta1.Ledger
			databaseConfiguration *v1beta1.DatabaseConfiguration
			orchestration         *v1beta1.Orchestration
			brokerConfiguration   *v1beta1.BrokerConfiguration
			temporalConfiguration *v1beta1.TemporalConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			gateway = &v1beta1.Gateway{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewaySpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Ingress: &v1beta1.GatewayIngress{},
				},
			}
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel:   stack.Name,
						core.ServiceLabel: "any",
					},
				},
				Spec: v1beta1.DatabaseConfigurationSpec{},
			}
			temporalConfiguration = &v1beta1.TemporalConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel:   stack.Name,
						core.ServiceLabel: "any",
					},
				},
				Spec: v1beta1.TemporalConfigurationSpec{},
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
			brokerConfiguration = &v1beta1.BrokerConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel: "any",
					},
				},
				Spec: v1beta1.BrokerConfigurationSpec{
					Nats: &v1beta1.BrokerNatsConfig{},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
			Expect(Create(orchestration)).To(Succeed())
			Expect(Create(brokerConfiguration)).To(Succeed())
			Expect(Create(temporalConfiguration)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(auth)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(orchestration)).To(Succeed())
			Expect(Delete(brokerConfiguration)).To(Succeed())
			Expect(Delete(temporalConfiguration)).To(Succeed())
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
