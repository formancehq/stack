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
	"k8s.io/client-go/dynamic"
)

var _ = Describe("Stacks informer", func() {
	var (
		membershipClientMock  *internal.MembershipClientMock
		inMemoryStacksModules map[string][]string
		startListener         func()
	)
	BeforeEach(func() {
		inMemoryStacksModules = map[string][]string{}
		membershipClientMock = internal.NewMembershipClientMock()
		dynamicClient, err := dynamic.NewForConfig(restConfig)
		Expect(err).To(Succeed())

		factory := internal.NewDynamicSharedInformerFactory(dynamicClient)
		Expect(internal.CreateStacksInformer(factory, logging.Testing(), membershipClientMock, inMemoryStacksModules)).To(Succeed())
		startListener = func() {
			stopCh := make(chan struct{})
			factory.Start(stopCh)
			DeferCleanup(func() {
				close(stopCh)
			})
		}
	})
	When("a stack is created on the cluster disabled", func() {
		var stack *v1beta1.Stack
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: v1beta1.StackSpec{
					Disabled: true,
				},
			}
			Expect(k8sClient.Post().
				Resource("Stacks").
				Body(stack).
				Do(context.Background()).
				Into(stack)).To(Succeed())

			inMemoryStacksModules[stack.Name] = []string{}

			startListener()

			DeferCleanup(func() {
				Expect(k8sClient.Delete().
					Resource("Stacks").
					Name(stack.Name).
					Do(context.Background()).Error()).To(Succeed())
			})
		})
		It("Should be disabled and have sent a Status_Disabled", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Disabled {
						return membershipClientMock.GetMessages()
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
		When("the stack is re-enabled", func() {
			BeforeEach(func() {
				stack.Spec.Disabled = false
				Expect(k8sClient.Put().
					Resource("Stacks").
					Name(stack.Name).
					Body(stack).
					Do(context.Background()).
					Into(stack)).To(Succeed())
			})
			It("should have sent a Status_Progressing", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Progressing {
							return membershipClientMock.GetMessages()
						}
					}
					return []*generated.Message{}
				}).ShouldNot(BeEmpty())
			})
			When("the stack is reconcilled", func() {
				BeforeEach(func() {
					stack.Status.Ready = true
					stack.Status.Modules = internal.GetExpectedModules(stack.Name, inMemoryStacksModules)
					Expect(
						k8sClient.Put().
							Resource("Stacks").
							SubResource("status").
							Name(stack.Name).
							Body(stack).
							Do(context.Background()).
							Error(),
					).To(Succeed())
				})
				It("should have sent a Status_Ready", func() {
					Eventually(func() []*generated.Message {
						for _, message := range membershipClientMock.GetMessages() {
							if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Ready && message.GetStatusChanged().ClusterName == stack.Name {
								return membershipClientMock.GetMessages()
							}
						}
						return nil
					}).ShouldNot(BeEmpty())
				})
			})
		})

	})
	When("Stack is created", func() {
		var stack *v1beta1.Stack
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
			}
			Expect(k8sClient.Post().
				Resource("Stacks").
				Body(stack).
				Do(context.Background()).
				Into(stack)).To(Succeed())

			inMemoryStacksModules[stack.Name] = []string{}

			startListener()

		})
		It("should have sent a Status_Progressing", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Progressing {
						return membershipClientMock.GetMessages()
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
		When("all stack dependent are ready", func() {
			BeforeEach(func() {
				stack.Status.Ready = true
				Expect(
					k8sClient.Put().
						Resource("Stacks").
						SubResource("status").
						Name(stack.Name).
						Body(stack).
						Do(context.Background()).
						Error(),
				).To(Succeed())
			})
			It("should have sent a Status_Ready", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Ready {
							return membershipClientMock.GetMessages()
						}
					}
					return nil
				}).ShouldNot(BeEmpty())
			})
		})
	})
	When("Stack is deleted", func() {
		BeforeEach(func() {
			stack := &v1beta1.Stack{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
			}
			Expect(k8sClient.Post().
				Resource("Stacks").
				Body(stack).
				Do(context.Background()).
				Into(stack)).To(Succeed())

			startListener()
			inMemoryStacksModules[stack.Name] = []string{}

			Expect(k8sClient.Delete().
				Resource("Stacks").
				Name(stack.Name).
				Do(context.Background()).Error()).To(Succeed())
		})
		It("should have sent a Stack_Deleted", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStackDeleted() != nil {
						return membershipClientMock.GetMessages()
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
	})
})
