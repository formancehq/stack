package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("TopicController", func() {
	Context("When creating a Topic", func() {
		var (
			stack               *v1beta1.Stack
			topic               *v1beta1.Topic
			brokerConfiguration *v1beta1.BrokerConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			brokerConfiguration = &v1beta1.BrokerConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/stack": stack.Name,
					},
				},
				Spec: v1beta1.BrokerConfigurationSpec{},
			}
			Expect(Create(brokerConfiguration)).To(Succeed())
			topic = &v1beta1.Topic{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: v1beta1.TopicSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Queries: []string{"orchestration"},
				},
			}
			Expect(Create(topic)).To(Succeed())
		})
		It("Should be set to ready status", func() {
			t := &v1beta1.Topic{}
			Eventually(func(g Gomega) bool {
				g.Expect(Get(internal.GetResourceName(topic.Name), t)).To(Succeed())
				return t.Status.Ready
			}).Should(BeTrue())
		})
		Context("Then updating removing all queries", func() {
			BeforeEach(func() {
				topic.Spec.Queries = []string{}
				Expect(Update(topic))
			})
			It("Should trigger the removal of the topic object", func() {
				Eventually(func(g Gomega) bool {
					return errors.IsNotFound(Get(internal.GetResourceName(topic.Name), topic))
				}).Should(BeTrue())
			})
		})
	})
})
