package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("GatewayController", func() {
	Context("When creating a Gateway", func() {
		var (
			stack       *v1beta1.Stack
			gateway     *v1beta1.Gateway
			httpService *v1beta1.HTTPAPI
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
					Ingress: &v1beta1.GatewayIngress{
						Host: "example.net",
					},
				},
			}
			httpService = &v1beta1.HTTPAPI{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.HTTPAPISpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Name:               "ledger",
					PortName:           "http",
					HasVersionEndpoint: true,
					Liveness:           v1beta1.LivenessDefault,
					Annotations: map[string]string{
						"foo": "bar",
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(BeNil())
			Expect(Create(httpService)).To(Succeed())
			Expect(Create(gateway)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(BeNil())
			Expect(Delete(httpService)).To(Succeed())
			Expect(Delete(gateway)).To(Succeed())
		})
		It("Should create a deployment", func() {
			Eventually(func() error {
				return LoadResource(stack.Name, "gateway", &appsv1.Deployment{})
			}).Should(Succeed())
		})
		It("Should create a service", func() {
			Eventually(func() error {
				return LoadResource(stack.Name, "gateway", &corev1.Service{})
			}).Should(Succeed())
		})
		It("Should create a config map with the Caddyfile", func() {
			Eventually(func(g Gomega) []string {
				g.Expect(LoadResource("", gateway.Name, gateway))

				return gateway.Status.SyncHTTPAPIs
			}).Should(ContainElements(httpService.Spec.Name))

			cm := &corev1.ConfigMap{}
			Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
			Expect(cm.Data["Caddyfile"]).To(
				MatchGoldenFile("gateway-controller", "configmap-with-ledger-only.yaml"))
		})
		Context("with a host defined", func() {
			JustBeforeEach(func() {
				patch := client.MergeFrom(gateway.DeepCopy())
				gateway.Spec.Ingress = &v1beta1.GatewayIngress{
					Host:   "example.com",
					Scheme: "https",
				}
				Expect(Patch(gateway, patch)).To(Succeed())
			})
			It("Should create an ingress", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, "gateway", &networkingv1.Ingress{})
				}).Should(Succeed())
			})
			Context("Then removing the hostname from the gateway", func() {
				var ingress *networkingv1.Ingress
				JustBeforeEach(func() {
					ingress = &networkingv1.Ingress{}
					Eventually(func() error {
						return LoadResource(stack.Name, "gateway", ingress)
					}).Should(Succeed())
					patch := client.MergeFrom(gateway.DeepCopy())
					gateway.Spec.Ingress = nil
					Expect(Patch(gateway, patch)).To(Succeed())
				})
				It("should delete the ingress", func() {
					Eventually(func() error {
						return LoadResource(stack.Name, "gateway", &networkingv1.Ingress{})
					}).Should(BeNotFound())
				})
			})
		})
		Context("Then adding a new HTTPService", func() {
			var anotherHttpService *v1beta1.HTTPAPI
			BeforeEach(func() {
				anotherHttpService = &v1beta1.HTTPAPI{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.HTTPAPISpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						Name:               "another",
						PortName:           "http",
						HasVersionEndpoint: true,
						Liveness:           v1beta1.LivenessDefault,
					},
				}
				Expect(Create(anotherHttpService)).To(Succeed())
			})
			It("Should trigger deployment gateway with the new service", func() {
				Eventually(func(g Gomega) []string {
					g.Expect(LoadResource("", gateway.Name, gateway))

					return gateway.Status.SyncHTTPAPIs
				}).Should(ContainElements(httpService.Spec.Name, anotherHttpService.Spec.Name))

				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-ledger-and-another-service.yaml"))
			})
		})
		Context("Then creating a Auth object", func() {
			BeforeEach(func() {
				auth := &v1beta1.Auth{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.AuthSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
					},
				}
				Expect(Create(auth)).To(Succeed())
			})
			It("Should redeploy the gateway with auth configuration", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", gateway.Name, gateway))

					return gateway.Status.AuthEnabled
				}).Should(BeTrue())
				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-ledger-and-auth.yaml"))
			})
		})
	})
})
