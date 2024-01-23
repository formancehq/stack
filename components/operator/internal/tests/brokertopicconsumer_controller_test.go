package tests_test

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("BrokerTopicConsumer", func() {
	Context("When creating a BrokerTopicConsumer", func() {
		var (
			brokerTopicConsumer   *v1beta1.BrokerTopicConsumer
			brokerNatsDSNSettings *v1beta1.Settings
			stack                 *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			brokerNatsDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
			Expect(Create(brokerNatsDSNSettings)).To(BeNil())
			brokerTopicConsumer = &v1beta1.BrokerTopicConsumer{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.BrokerTopicConsumerSpec{
					Service:   "ledger",
					QueriedBy: "orchestration",
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			Expect(Create(brokerTopicConsumer)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(brokerNatsDSNSettings)).To(Succeed())
			Expect(client.IgnoreNotFound(Delete(brokerTopicConsumer))).To(Succeed())
		})
		It("Should create a BrokerTopic", func() {
			t := &v1beta1.BrokerTopic{}
			Eventually(func(g Gomega) *v1beta1.BrokerTopic {
				g.Expect(Get(core.GetResourceName(
					core.GetObjectName(stack.Name, brokerTopicConsumer.Spec.Service)), t)).To(Succeed())
				return t
			}).Should(BeOwnedBy(brokerTopicConsumer))
		})
		Context("Then when the BrokerTopic is ready", func() {
			t := &v1beta1.BrokerTopic{}
			BeforeEach(func() {
				Eventually(func(g Gomega) bool {
					g.Expect(Get(core.GetResourceName(
						core.GetObjectName(stack.Name, brokerTopicConsumer.Spec.Service)), t)).To(Succeed())
					return t.Status.Ready
				}).Should(BeTrue())
			})
			It("Should set the BrokerTopicConsumer to ready status", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", brokerTopicConsumer.Name, brokerTopicConsumer)).To(Succeed())

					return brokerTopicConsumer.Status.Ready
				}).Should(BeTrue())
			})
			Context("Then create a new BrokerTopicConsumer on the same service", func() {
				brokerTopicConsumer2 := &v1beta1.BrokerTopicConsumer{}
				BeforeEach(func() {
					brokerTopicConsumer2 = &v1beta1.BrokerTopicConsumer{
						ObjectMeta: RandObjectMeta(),
						Spec: v1beta1.BrokerTopicConsumerSpec{
							Service:   brokerTopicConsumer.Spec.Service,
							QueriedBy: "webhooks",
							StackDependency: v1beta1.StackDependency{
								Stack: stack.Name,
							},
						},
					}
					Expect(Create(brokerTopicConsumer2)).To(Succeed())
				})
				AfterEach(func() {
					Expect(client.IgnoreNotFound(Delete(brokerTopicConsumer2))).To(Succeed())
				})
				It("Should be set to ready too", func() {
					Eventually(func(g Gomega) bool {
						g.Expect(LoadResource("", brokerTopicConsumer2.Name, brokerTopicConsumer2)).To(Succeed())

						return brokerTopicConsumer2.Status.Ready
					}).Should(BeTrue())
				})
				Context("Then first BrokerTopicConsumer object", func() {
					BeforeEach(func() {
						Expect(Delete(brokerTopicConsumer)).To(Succeed())
					})
					It("Should remove the service from the queries of the topic", func() {
						Eventually(func(g Gomega) *v1beta1.BrokerTopic {
							topic := &v1beta1.BrokerTopic{}
							g.Expect(Get(core.GetResourceName(core.GetObjectName(stack.Name, brokerTopicConsumer.Spec.Service)), topic)).To(Succeed())
							return topic
						}).ShouldNot(BeControlledBy(brokerTopicConsumer))
					})
					Context("Then removing the last BrokerTopicConsumer", func() {
						BeforeEach(func() {
							Expect(Delete(brokerTopicConsumer2)).To(Succeed())
						})
						It("Should completely remove the BrokerTopic object", func() {
							Eventually(func(g Gomega) bool {
								t := &v1beta1.BrokerTopic{}
								err := Get(core.GetResourceName(core.GetObjectName(stack.Name, brokerTopicConsumer.Spec.Service)), t)

								return errors.IsNotFound(err)
							}).Should(BeTrue())
						})
					})
				})
			})
		})
	})
})
