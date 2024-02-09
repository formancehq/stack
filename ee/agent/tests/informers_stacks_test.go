package tests

import (
	"context"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

var _ = Describe("Stacks informer", func() {
	var (
		membershipClientMock *internal.MembershipClientMock
	)
	BeforeEach(func() {
		membershipClientMock = internal.NewMembershipClientMock()
		dynamicClient, err := dynamic.NewForConfig(restConfig)
		Expect(err).To(Succeed())

		factory := internal.NewDynamicSharedInformerFactory(dynamicClient)
		stacksInformer, err := internal.CreateStacksInformer(factory, logging.Testing(), membershipClientMock)
		Expect(err).To(Succeed())
		stopCh := make(chan struct{})
		go stacksInformer.Run(stopCh)
		DeferCleanup(func() {
			close(stopCh)
		})
	})
	When("a stack is created on the cluster then disabling it", func() {
		stack := &v1beta1.Stack{}
		BeforeEach(func() {
			stack.ObjectMeta = v1.ObjectMeta{
				Name: uuid.NewString(),
			}
			Expect(k8sClient.Post().
				Resource("Stacks").
				Body(&v1beta1.Stack{
					ObjectMeta: v1.ObjectMeta{
						Name: uuid.NewString(),
					},
				}).
				Do(context.Background()).
				Into(stack)).To(Succeed())
		})
		It("Should be disabled and have sent a Status_Progressing", func() {
			Expect(
				k8sClient.Patch(types.MergePatchType).
					Resource("Stacks").
					Name(stack.Name).
					Body([]byte(`{"spec": {"disabled": true}}`)).
					Do(context.Background()).
					Error(),
			).To(Succeed())

			Expect(k8sClient.Get().Resource("Stacks").Name(stack.Name).Do(context.Background()).Into(stack)).To(Succeed())
			Expect(stack.Spec.Disabled).To(BeTrue())

			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Progressing {
						return membershipClientMock.GetMessages()
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
		When("Stack is fully disabled, mean reconcille and ready ", func() {
			BeforeEach(func() {
				stack.Spec.Disabled = true
				Expect(
					k8sClient.Put().
						Resource("Stacks").
						Name(stack.Name).
						Body(stack).
						Do(context.Background()).
						Error(),
				).To(Succeed())
			})
			JustBeforeEach(func() {
				Expect(
					k8sClient.Patch(types.MergePatchType).
						Resource("Stacks").
						SubResource("status").
						Name(stack.Name).
						Body([]byte(`{"status": {"ready": true}}`)).
						Do(context.Background()).
						Error(),
				).To(Succeed())
			})
			It("should have sent a Status_disabled", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Disabled {
							return membershipClientMock.GetMessages()
						}
					}
					return nil
				}).ShouldNot(BeEmpty())
			})
		})
	})
})
