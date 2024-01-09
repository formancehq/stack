package formance_com_test

import (
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/formance.com/testing"
	"github.com/formancehq/operator/v2/internal/core"
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
			auth                  *v1beta1.Auth
			databaseConfiguration *v1beta1.DatabaseConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel:   stack.Name,
						core.ServiceLabel: "any",
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
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(auth)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(auth)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "auth", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(auth))
			Expect(deployment.Spec.Template.Spec.Containers[0].Env).To(ContainElements(
				core.Env("BASE_URL", "http://auth:8080"),
			))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "auth"), httpService)
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
		Context("with a Gateway", func() {
			var (
				gateway *v1beta1.Gateway
			)
			BeforeEach(func() {
				gateway = &v1beta1.Gateway{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.GatewaySpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
					},
				}
			})
			JustBeforeEach(func() {
				Expect(Create(gateway)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(gateway)).To(Succeed())
			})
			It("Should create a deployment with proper BASE_URL env var", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func(g Gomega) []corev1.EnvVar {
					g.Expect(LoadResource(stack.Name, "auth", deployment)).To(Succeed())
					return deployment.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(
					core.Env("BASE_URL", "http://gateway:8080/api/auth"),
				))
			})
			Context("with an ingress", func() {
				BeforeEach(func() {
					gateway.Spec.Ingress = &v1beta1.GatewayIngress{
						Host:   "example.net",
						Scheme: "https",
					}
				})
				It("Should create a deployment with proper BASE_URL env var", func() {
					deployment := &appsv1.Deployment{}
					Eventually(func(g Gomega) []corev1.EnvVar {
						g.Expect(LoadResource(stack.Name, "auth", deployment)).To(Succeed())
						return deployment.Spec.Template.Spec.Containers[0].Env
					}).Should(ContainElements(
						core.Env("BASE_URL", "https://example.net/api/auth"),
					))
				})
			})
		})
	})
})
