package controllers_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controllers/testing"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("TopicQueryController", func() {
	Context("When creating a TopicQuery", func() {
		var (
			topicQuery   *v1beta1.TopicQuery
			brokerConfig *v1beta1.BrokerConfiguration
			stack        *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			brokerConfig = &v1beta1.BrokerConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel: stack.Name,
					},
				},
				Spec: v1beta1.BrokerConfigurationSpec{},
			}
			Expect(Create(brokerConfig)).To(Succeed())
			topicQuery = &v1beta1.TopicQuery{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.TopicQuerySpec{
					Service:   "ledger",
					QueriedBy: "orchestration",
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			Expect(Create(topicQuery)).To(Succeed())
		})
		It("Should create a Topic", func() {
			t := &v1beta1.Topic{}
			Eventually(func(g Gomega) *v1beta1.Topic {
				g.Expect(Get(core.GetResourceName(
					core.GetObjectName(stack.Name, topicQuery.Spec.Service)), t)).To(Succeed())
				return t
			}).Should(BeOwnedBy(topicQuery))
		})
		Context("Then when the Topic is ready", func() {
			t := &v1beta1.Topic{}
			BeforeEach(func() {
				Eventually(func(g Gomega) bool {
					g.Expect(Get(core.GetResourceName(
						core.GetObjectName(stack.Name, topicQuery.Spec.Service)), t)).To(Succeed())
					return t.Status.Ready
				}).Should(BeTrue())
			})
			It("Should set the TopicQuery to ready status", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", topicQuery.Name, topicQuery)).To(Succeed())

					return topicQuery.Status.Ready
				}).Should(BeTrue())
			})
			Context("Then create a new TopicQuery on the same service", func() {
				topicQuery2 := &v1beta1.TopicQuery{}
				BeforeEach(func() {
					topicQuery2 = &v1beta1.TopicQuery{
						ObjectMeta: RandObjectMeta(),
						Spec: v1beta1.TopicQuerySpec{
							Service:   topicQuery.Spec.Service,
							QueriedBy: "webhooks",
							StackDependency: v1beta1.StackDependency{
								Stack: stack.Name,
							},
						},
					}
					Expect(Create(topicQuery2)).To(Succeed())
				})
				It("Should be set to ready too", func() {
					Eventually(func(g Gomega) bool {
						g.Expect(LoadResource("", topicQuery2.Name, topicQuery2)).To(Succeed())

						return topicQuery2.Status.Ready
					}).Should(BeTrue())
				})
				Context("Then first TopicQuery object", func() {
					BeforeEach(func() {
						Expect(Delete(topicQuery)).To(Succeed())
					})
					It("Should remove the service from the queries of the topic", func() {
						Eventually(func(g Gomega) *v1beta1.Topic {
							topic := &v1beta1.Topic{}
							g.Expect(Get(core.GetResourceName(core.GetObjectName(stack.Name, topicQuery.Spec.Service)), topic)).To(Succeed())
							return topic
						}).ShouldNot(BeControlledBy(topicQuery))
					})
					Context("Then removing the last TopicQuery", func() {
						BeforeEach(func() {
							Expect(Delete(topicQuery2)).To(Succeed())
						})
						It("Should completely remove the Topic object", func() {
							Eventually(func(g Gomega) bool {
								t := &v1beta1.Topic{}
								err := Get(core.GetResourceName(core.GetObjectName(stack.Name, topicQuery.Spec.Service)), t)
								spew.Dump(t)

								return errors.IsNotFound(err)
							}).Should(BeTrue())
						})
					})
				})
			})
		})
	})
})
