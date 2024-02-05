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
	Context("When creating a stack on the cluster then updating its status", func() {
		It("Should trigger a new update of the status on membership client", func() {
			stack := &v1beta1.Stack{}
			Expect(k8sClient.Post().
				Resource("Stacks").
				Body(&v1beta1.Stack{
					ObjectMeta: v1.ObjectMeta{
						Name: uuid.NewString(),
					},
				}).
				Do(context.Background()).
				Into(stack)).To(Succeed())

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

			Eventually(func() []*generated.Message {
				return membershipClientMock.GetMessages()
			}).ShouldNot(BeEmpty())
		})
	})
})
