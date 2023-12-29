package controllers_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/common"
	. "github.com/formancehq/operator/v2/internal/controllers/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("StreamProcessorController", func() {

	Context("When creating a stream", func() {
		var (
			streamProcessor            *v1beta1.StreamProcessor
			stack                      *v1beta1.Stack
			elasticSearchConfiguration *v1beta1.ElasticSearchConfiguration
			brokerConfiguration        *v1beta1.BrokerConfiguration
		)
		BeforeEach(func() {

			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(Succeed())

			elasticSearchConfiguration = &v1beta1.ElasticSearchConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/stack": stack.Name,
					},
				},
				Spec: v1beta1.ElasticSearchConfigurationSpec{},
			}
			Expect(Create(elasticSearchConfiguration)).To(Succeed())

			brokerConfiguration = &v1beta1.BrokerConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/stack": stack.Name,
					},
				},
				Spec: v1beta1.BrokerConfigurationSpec{},
			}
			Expect(Create(brokerConfiguration)).To(Succeed())

			streamProcessor = &v1beta1.StreamProcessor{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.StreamProcessorSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			Expect(Create(streamProcessor)).To(Succeed())
		})
		It("Should create a deployment", func() {
			t := &appsv1.Deployment{}
			Eventually(func() error {
				return Get(common.GetNamespacedResourceName(stack.Name, "stream-processor"), t)
			}).Should(BeNil())
		})
		It("Should create a ConfigMap for templates configuration", func() {
			t := &corev1.ConfigMap{}
			Eventually(func() error {
				return Get(common.GetNamespacedResourceName(stack.Name, "stream-processor-templates"), t)
			}).Should(BeNil())
		})
		It("Should create a ConfigMap for resources configuration", func() {
			t := &corev1.ConfigMap{}
			Eventually(func() error {
				return Get(common.GetNamespacedResourceName(stack.Name, "stream-processor-resources"), t)
			}).Should(BeNil())
		})
	})
})
