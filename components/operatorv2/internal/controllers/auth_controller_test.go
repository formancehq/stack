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

var _ = Describe("AuthController", func() {
	Context("When creating a Auth object", func() {
		var (
			stack                 *v1beta1.Stack
			gateway               *v1beta1.Gateway
			auth                  *v1beta1.Auth
			databaseConfiguration *v1beta1.DatabaseConfiguration
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
			Expect(Create(gateway)).To(Succeed())
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(auth)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "auth", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeOwnedBy(auth))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", common.GetObjectName(stack.Name, "auth"), httpService)
			}).Should(Succeed())
		})
		Context("Then when create an AuthClient object", func() {
			var (
				authClient *v1beta1.AuthClient
			)
			JustBeforeEach(func() {
				authClient = &v1beta1.AuthClient{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.AuthClientSpec{
						ID: "client0",
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						Secret: "secret",
					},
				}
				Expect(Create(authClient)).To(Succeed())
				Eventually(func(g Gomega) []string {
					g.Expect(LoadResource("", auth.Name, auth)).To(Succeed())
					return auth.Status.Clients
				}).Should(ContainElements(authClient.Name))
			})
			It("Should configure the config map with the auth client", func() {
				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "auth-configuration", cm)).To(Succeed())
				Expect(cm.Data["config.yaml"]).To(MatchGoldenFile("auth-controller", "config-with-auth-client.yaml"))
			})
		})
	})
})
