package controllers_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controllers/testing"
	"github.com/formancehq/operator/v2/internal/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("StreamController", func() {
	Context("When creating a Stream", func() {
		var (
			stream *v1beta1.Stream
			stack  *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			stream = &v1beta1.Stream{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.StreamSpec{
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
