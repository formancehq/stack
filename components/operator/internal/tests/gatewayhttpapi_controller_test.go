package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("GatewayHTTPAPI", func() {
	Context("When creating an GatewayHTTPAPI", func() {
		var (
			stack   *v1beta1.Stack
			httpAPI *v1beta1.GatewayHTTPAPI
		)
		BeforeEach(func() {

			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			httpAPI = &v1beta1.GatewayHTTPAPI{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewayHTTPAPISpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Name:  "ledger",
					Rules: []v1beta1.GatewayHTTPAPIRule{gatewayhttpapis.RuleSecured()},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(BeNil())
			Expect(Create(httpAPI)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(httpAPI)).To(Succeed())
			Expect(Delete(stack)).To(BeNil())
		})
		It("Should create a k8s service", func() {
			service := &corev1.Service{}
			Eventually(func() error {
				return LoadResource(stack.Name, "ledger", service)
			}).Should(BeNil())
			Expect(service).To(BeControlledBy(httpAPI))
			Expect(service.Spec.Selector).To(Equal(map[string]string{
				"app.kubernetes.io/name": httpAPI.Spec.Name,
			}))
		})
		Context("With user defined annotations", func() {
			var (
				annotationsSettings *v1beta1.Settings
			)
			JustBeforeEach(func() {
				annotationsSettings = settings.New(uuid.NewString(), "services.*.annotations", "foo=bar", stack.Name)
				Expect(Create(annotationsSettings)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(annotationsSettings)).To(Succeed())
			})
			It("should add annotations to the service", func() {
				Eventually(func(g Gomega) bool {
					service := &corev1.Service{}
					g.Expect(LoadResource(stack.Name, "ledger", service)).To(Succeed())
					g.Expect(service.Annotations).To(Equal(map[string]string{
						"foo": "bar",
					}))
					return true
				}).Should(BeTrue())
			})
		})
	})
})
