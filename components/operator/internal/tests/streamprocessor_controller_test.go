package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("StreamProcessorController", func() {

	Context("When creating a stream processor", func() {
		var (
			streamProcessor            *v1beta1.StreamProcessor
			stack                      *v1beta1.Stack
			elasticSearchConfiguration *v1beta1.ElasticSearchConfiguration
			brokerKindSettings         *v1beta1.Settings
			brokerNatsEndpointSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			elasticSearchConfiguration = &v1beta1.ElasticSearchConfiguration{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.ElasticSearchConfigurationSpec{
					ConfigurationProperties: v1beta1.ConfigurationProperties{
						Stacks: []string{stack.Name},
					},
				},
			}
			brokerKindSettings = settings.New(uuid.NewString(), "broker.kind", "nats", stack.Name)
			brokerNatsEndpointSettings = settings.New(uuid.NewString(), "broker.nats.endpoint", "localhost:1234", stack.Name)
			streamProcessor = &v1beta1.StreamProcessor{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.StreamProcessorSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(brokerKindSettings)).To(BeNil())
			Expect(Create(brokerNatsEndpointSettings)).To(BeNil())
			Expect(Create(stack)).To(Succeed())
			Expect(Create(elasticSearchConfiguration)).To(Succeed())
			Expect(Create(streamProcessor)).To(Succeed())
		})
		JustAfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(elasticSearchConfiguration)).To(Succeed())
			Expect(Delete(brokerNatsEndpointSettings)).To(Succeed())
			Expect(Delete(brokerKindSettings)).To(Succeed())
			Expect(Delete(streamProcessor)).To(Succeed())
		})
		It("Should create a deployment", func() {
			t := &appsv1.Deployment{}
			Eventually(func() error {
				return Get(core.GetNamespacedResourceName(stack.Name, "stream-processor"), t)
			}).Should(BeNil())
		})
		It("Should create a ConfigMap for templates configuration", func() {
			t := &corev1.ConfigMap{}
			Eventually(func() error {
				return Get(core.GetNamespacedResourceName(stack.Name, "stream-processor-templates"), t)
			}).Should(BeNil())
		})
		It("Should create a ConfigMap for resources configuration", func() {
			t := &corev1.ConfigMap{}
			Eventually(func() error {
				return Get(core.GetNamespacedResourceName(stack.Name, "stream-processor-resources"), t)
			}).Should(BeNil())
		})
		Context("with audit enabled on stack", func() {
			BeforeEach(func() {
				stack.Spec.EnableAudit = true
			})
			It("should add a config map for the stream", func() {
				Eventually(func() error {
					cm := &corev1.ConfigMap{}
					return LoadResource(stack.Name, "stream-processor-audit", cm)
				}).Should(Succeed())
			})
			It("should add a cmd args to the deployment", func() {
				t := &appsv1.Deployment{}
				Eventually(func(g Gomega) []string {
					g.Expect(LoadResource(stack.Name, "stream-processor", t)).To(Succeed())
					return t.Spec.Template.Spec.Containers[0].Command
				}).Should(ContainElement("/audit/gateway_audit.yaml"))
			})
			Context("then disabling audit", func() {
				JustBeforeEach(func() {
					Eventually(func() error {
						cm := &corev1.ConfigMap{}
						return LoadResource(stack.Name, "stream-processor-audit", cm)
					}).Should(Succeed())
					patch := client.MergeFrom(stack.DeepCopy())
					stack.Spec.EnableAudit = false
					Expect(Patch(stack, patch)).To(Succeed())
				})
				It("should remove the associated config map", func() {
					Eventually(func() error {
						return LoadResource(stack.Name, "stream-processor-audit", &corev1.ConfigMap{})
					}).Should(BeNotFound())
				})
			})
		})
	})
})
