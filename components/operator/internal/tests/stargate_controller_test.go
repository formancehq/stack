package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/operator/internal/tests/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("StargateController", func() {
	Context("When creating a Stargate object", func() {
		var (
			stack    *v1beta1.Stack
			stargate *v1beta1.Stargate
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			stargate = &v1beta1.Stargate{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.StargateSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					ServerURL:      "server:8080",
					OrganizationID: "orgID",
					StackID:        "stackID",
					Auth: v1beta1.StargateAuthSpec{
						ClientID:     "client0",
						ClientSecret: "client0",
						Issuer:       "http://server:8081",
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(stargate)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stargate)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should add an owner reference on the stack", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", stargate.Name, stargate)).To(Succeed())
				reference, err := core.HasOwnerReference(TestContext(), stack, stargate)
				g.Expect(err).To(BeNil())
				return reference
			}).Should(BeTrue())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "stargate", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(stargate))
		})
	})
})
