package tests

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/google/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

var _ = Describe("Stacks informer", func() {
	var (
		membershipClientMock *internal.MembershipClientMock
		startListener        func()
	)
	BeforeEach(func() {
		membershipClientMock = internal.NewMembershipClientMock()
		dynamicClient, err := dynamic.NewForConfig(restConfig)
		Expect(err).To(Succeed())

		factory := internal.NewDynamicSharedInformerFactory(dynamicClient, 5*time.Minute)
		Expect(internal.CreateStacksInformer(factory, logging.Testing(), membershipClientMock)).To(Succeed())
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
			By("Creating a disabled stack", func() {
				Expect(k8sClient.Post().
					Resource("Stacks").
					Body(stack).
					Do(context.Background()).
					Into(stack)).To(Succeed())
			})

			By("Adding ready", func() {
				stack.Status.Ready = true
				patch, err := json.Marshal(struct {
					Status v1beta1.StackStatus `json:"status"`
				}{
					Status: stack.Status,
				})
				Expect(err).To(Succeed())

				Expect(k8sClient.Patch(types.MergePatchType).
					Resource("Stacks").
					SubResource("status").
					Name(stack.Name).
					Body(patch).
					Do(context.Background()).
					Error()).To(Succeed())
			})

			startListener()

			DeferCleanup(func() {
				Expect(k8sClient.Delete().
					Resource("Stacks").
					Name(stack.Name).
					Do(context.Background()).Error()).To(Succeed())
			})
		})
		It("Should have sent a Status_Changed", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStatusChanged() != nil {
						if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
							isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
							if isReady {
								return membershipClientMock.GetMessages()
							}
						}
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
		When("The stack is ready", func() {
			BeforeEach(func() {
				stack.Status.Ready = true
				patch, err := json.Marshal(struct {
					Status v1beta1.StackStatus `json:"status,omitempty"`
				}{
					Status: stack.Status,
				})
				Expect(err).To(Succeed())
				Expect(k8sClient.Patch(types.MergePatchType).
					Resource("Stacks").
					SubResource("status").
					Name(stack.Name).
					Body(patch).
					Do(context.Background()).
					Error()).To(Succeed())
			})
			It("Should have sent a Status_Changed", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil {
							if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
								isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
								if isReady {
									return membershipClientMock.GetMessages()
								}
							}
						}
					}
					return nil
				}).ShouldNot(BeEmpty())
			})
		})
		When("the stack is re-enabled", func() {
			BeforeEach(func() {
				By("setting the status ready", func() {
					stack.Status.Ready = false
					patch, err := json.Marshal(struct {
						Status v1beta1.StackStatus `json:"status"`
					}{
						Status: stack.Status,
					})
					Expect(err).To(Succeed())

					Expect(k8sClient.Patch(types.MergePatchType).
						Resource("Stacks").
						SubResource("status").
						Name(stack.Name).
						Body(patch).
						Do(context.Background()).
						Error()).To(Succeed())
				})
				By("Enabling the stack", func() {
					stack.Spec.Disabled = false
					path, err := json.Marshal(struct {
						Spec v1beta1.StackSpec `json:"spec"`
					}{
						Spec: stack.Spec,
					})
					Expect(err).To(Succeed())
					Expect(k8sClient.Patch(types.MergePatchType).
						Resource("Stacks").
						Name(stack.Name).
						Body(path).
						Do(context.Background()).
						Error()).To(Succeed())
				})
			})
			It("should have sent a Status_Changed", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil {
							if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
								isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
								if !isReady {
									return membershipClientMock.GetMessages()
								}
							}
						}
					}
					return []*generated.Message{}
				}).ShouldNot(BeEmpty())
			})
			When("the stack is reconcilled", func() {
				BeforeEach(func() {
					By("Setting the status ready", func() {
						stack.Status.Ready = true
						stack.Status.Modules = []string{}
						patch, err := json.Marshal(struct {
							Status v1beta1.StackStatus `json:"status"`
						}{
							Status: stack.Status,
						})
						Expect(err).To(Succeed())

						Expect(k8sClient.Patch(types.MergePatchType).
							Resource("Stacks").
							SubResource("status").
							Name(stack.Name).
							Body(patch).
							Do(context.Background()).
							Error()).To(Succeed())
					})
				})
				It("should have sent a Status_Changed", func() {
					Eventually(func() []*generated.Message {
						for _, message := range membershipClientMock.GetMessages() {
							if message.GetStatusChanged() != nil {
								if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
									isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
									if isReady {
										return membershipClientMock.GetMessages()
									}
								}
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

			By("Disabling the status ready", func() {
				stack.Status.Ready = false
				patch, err := json.Marshal(struct {
					Status v1beta1.StackStatus `json:"status"`
				}{
					Status: stack.Status,
				})
				Expect(err).To(Succeed())

				Expect(k8sClient.Patch(types.MergePatchType).
					Resource("Stacks").
					SubResource("status").
					Name(stack.Name).
					Body(patch).
					Do(context.Background()).
					Error()).To(Succeed())
			})

			startListener()
		})
		AfterEach(func() {
			Expect(k8sClient.Delete().
				Resource("Stacks").
				Name(stack.Name).
				Do(context.Background()).Error()).To(Succeed())
		})
		It("should have sent a Status_Changed", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStatusChanged() != nil && message.GetStatusChanged().Status == generated.StackStatus_Progressing && stack.Name == message.GetStatusChanged().ClusterName {
						if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
							isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
							if !isReady {
								return membershipClientMock.GetMessages()
							}
						}
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
		When("all stack dependent are ready", func() {
			BeforeEach(func() {
				By("setting the status ready", func() {
					stack.Status.Ready = true
					patch, err := json.Marshal(struct {
						Status v1beta1.StackStatus `json:"status"`
					}{
						Status: stack.Status,
					})
					Expect(err).To(Succeed())

					Expect(k8sClient.Patch(types.MergePatchType).
						Resource("Stacks").
						SubResource("status").
						Name(stack.Name).
						Body(patch).
						Do(context.Background()).
						Error()).To(Succeed())
				})
			})
			It("should have sent a Status_Ready", func() {
				Eventually(func() []*generated.Message {
					for _, message := range membershipClientMock.GetMessages() {
						if message.GetStatusChanged() != nil {
							if _, ok := message.GetStatusChanged().GetStatuses().Fields["ready"]; ok {
								isReady := message.GetStatusChanged().GetStatuses().Fields["ready"].GetBoolValue()
								if isReady {
									return membershipClientMock.GetMessages()
								}
							}
						}
					}
					return nil
				}).ShouldNot(BeEmpty())
			})
		})
	})
	When("Stack is deleted", func() {
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

			startListener()

			Expect(k8sClient.Delete().
				Resource("Stacks").
				Name(stack.Name).
				Do(context.Background()).Error()).To(Succeed())
		})
		It("should have sent a Stack_Deleted", func() {
			Eventually(func() []*generated.Message {
				for _, message := range membershipClientMock.GetMessages() {
					if message.GetStackDeleted() != nil && message.GetStackDeleted().ClusterName == stack.Name {
						return membershipClientMock.GetMessages()
					}
				}
				return nil
			}).ShouldNot(BeEmpty())
		})
	})
})
