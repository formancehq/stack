package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("WalletController", func() {
	Context("When creating a Wallet object", func() {
		var (
			stack                 *v1beta1.Stack
			gateway               *v1beta1.Gateway
			ledger                *v1beta1.Ledger
			databaseConfiguration *v1beta1.DatabaseConfiguration
			wallet                *v1beta1.Wallet
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			gateway = &v1beta1.Gateway{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.GatewaySpec{
					Stack:   stack.Name,
					Ingress: &v1beta1.GatewayIngress{},
				},
			}
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/stack":   stack.Name,
						"formance.com/service": "any",
					},
				},
				Spec: v1beta1.DatabaseConfigurationSpec{},
			}
			wallet = &v1beta1.Wallet{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.WalletSpec{
					Stack: stack.Name,
				},
			}
			ledger = &v1beta1.Ledger{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.LedgerSpec{
					Stack: stack.Name,
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
			Expect(Create(wallet)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(wallet)).To(Succeed())
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "wallets", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeOwnedBy(wallet))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", internal.GetObjectName(stack.Name, "wallets"), httpService)
			}).Should(Succeed())
		})
		It("Should create a new AuthClient object", func() {
			authClient := &v1beta1.AuthClient{}
			Eventually(func() error {
				return LoadResource("", internal.GetObjectName(stack.Name, "wallets"), authClient)
			}).Should(Succeed())
		})
	})
})