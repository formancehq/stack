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
			stack            *v1beta1.Stack
			gateway          *v1beta1.Gateway
			ledger           *v1beta1.Ledger
			databaseSettings *v1beta1.Settings
			wallets          *v1beta1.Wallets
			auth             *v1beta1.Auth
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			gateway = &v1beta1.Gateway{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewaySpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Ingress: &v1beta1.GatewayIngress{},
				},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
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
			Expect(Create(stack)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(databaseSettings)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
			Expect(Create(wallets)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(auth)).To(Succeed())
			Expect(Delete(wallets)).To(Succeed())
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "wallets", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(wallets))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
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
