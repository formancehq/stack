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
)

var _ = Describe("WalletsController", func() {
	Context("When creating a Wallets object", func() {
		var (
			stack                    *v1beta1.Stack
			gateway                  *v1beta1.Gateway
			ledger                   *v1beta1.Ledger
			wallets                  *v1beta1.Wallets
			auth                     *v1beta1.Auth
			databaseSettings         *v1beta1.Settings
			resourceLimitsSettings   *v1beta1.Settings
			resourceRequestsSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
			resourceLimitsSettings = settings.New(uuid.NewString(),
				"deployments.*.containers.*.resource-requirements.limits", "cpu=500m", stack.Name)
			resourceRequestsSettings = settings.New(uuid.NewString(),
				"deployments.*.containers.*.resource-requirements.requests", "cpu=250m", stack.Name)
			gateway = &v1beta1.Gateway{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewaySpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Ingress: &v1beta1.GatewayIngress{},
				},
			}
			wallets = &v1beta1.Wallets{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.WalletsSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			ledger = &v1beta1.Ledger{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.LedgerSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			auth = &v1beta1.Auth{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.AuthSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(
				stack, databaseSettings, resourceLimitsSettings, resourceRequestsSettings,
				auth, gateway, ledger, wallets,
			)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack, databaseSettings, resourceLimitsSettings, resourceRequestsSettings)).To(Succeed())
		})
		It("Should add an owner reference on the stack", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", wallets.Name, wallets)).To(Succeed())
				reference, err := core.HasOwnerReference(TestContext(), stack, wallets)
				g.Expect(err).To(BeNil())
				return reference
			}).Should(BeTrue())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "wallets", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(wallets))
			By("should have set proper resource requirements", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).NotTo(BeEmpty())
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Requests).NotTo(BeEmpty())
			})
		})
		It("Should create a new GatewayHTTPAPI object", func() {
			httpService := &v1beta1.GatewayHTTPAPI{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "wallets"), httpService)
			}).Should(Succeed())
		})
		It("Should create a new AuthClient object", func() {
			authClient := &v1beta1.AuthClient{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "wallets"), authClient)
			}).Should(Succeed())
		})
	})
})
