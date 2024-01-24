package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/operator/internal/tests/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("StreamController", func() {
	Context("When creating a BenthosStream", func() {
		var (
			stream *v1beta1.BenthosStream
			stack  *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			stream = &v1beta1.BenthosStream{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.BenthosStreamSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			Expect(Create(stream)).To(Succeed())
		})
		It("Should create a ConfigMap", func() {
			t := &corev1.ConfigMap{}
			Eventually(func() error {
				return Get(core.GetNamespacedResourceName(stack.Name, "stream-"+stream.Name), t)
			}).Should(BeNil())
		})
	})
})
