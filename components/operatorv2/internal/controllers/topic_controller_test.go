package controllers_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controllers/testing"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
				},
			}
			Expect(controllerutil.SetOwnerReference(stack, topic, GetScheme())).To(Succeed())
			Expect(Create(topic)).To(Succeed())
		})
		It("Should be set to ready status", func() {
			t := &v1beta1.Topic{}
			Eventually(func(g Gomega) bool {
				g.Expect(Get(core.GetResourceName(topic.Name), t)).To(Succeed())
				return t.Status.Ready
			}).Should(BeTrue())
		})
		Context("Then updating removing all queries", func() {
			BeforeEach(func() {
				patch := client.MergeFrom(topic.DeepCopy())
				Expect(controllerutil.RemoveOwnerReference(stack, topic, GetScheme())).To(Succeed())
				Expect(Patch(topic, patch)).To(Succeed())
			})
			It("Should trigger the deletion of the topic object", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", topic.Name, topic)
				}).Should(BeNotFound())
			})
		})
	})
})
