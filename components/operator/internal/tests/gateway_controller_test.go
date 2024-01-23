package tests_test

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/google/uuid"

	. "github.com/formancehq/operator/internal/tests/internal"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/httpapis"
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
			stack   *v1beta1.Stack
			gateway *v1beta1.Gateway
			httpAPI *v1beta1.HTTPAPI
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
						Host:   "example.net",
						Scheme: "https",
					},
				},
			}
			httpAPI = &v1beta1.HTTPAPI{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.HTTPAPISpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Name: "ledger",
					Service: &v1beta1.ServiceConfiguration{
						Annotations: map[string]string{
							"foo": "bar",
						},
					},
					Rules: []v1beta1.HTTPAPIRule{httpapis.RuleSecured()},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(BeNil())
			Expect(Create(httpAPI)).To(Succeed())
			Expect(Create(gateway)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(BeNil())
			Expect(Delete(httpAPI)).To(Succeed())
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
			}).Should(ContainElements(httpAPI.Spec.Name))

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
						Name: "another",
						Rules: []v1beta1.HTTPAPIRule{
							{
								Path:    "/webhooks",
								Methods: []string{"POST"},
								Secured: true,
							},
							httpapis.RuleSecured(),
						},
					},
				}
				Expect(Create(anotherHttpService)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(anotherHttpService)).To(Succeed())
			})
			It("Should trigger deployment gateway with the new service", func() {
				Eventually(func(g Gomega) []string {
					g.Expect(LoadResource("", gateway.Name, gateway))

					return gateway.Status.SyncHTTPAPIs
				}).Should(ContainElements(httpAPI.Spec.Name, anotherHttpService.Spec.Name))

				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-ledger-and-another-service.yaml"))
			})
		})
		Context("Then creating a Auth object", func() {
			var auth *v1beta1.Auth
			BeforeEach(func() {
				auth = &v1beta1.Auth{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.AuthSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
					},
				}
				Expect(Create(auth)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(auth)).To(Succeed())
			})
			It("Should redeploy the gateway with auth configuration", func() {
				Eventually(func(g Gomega) []string {
					g.Expect(LoadResource("", gateway.Name, gateway))
					return gateway.Status.SyncHTTPAPIs
				}).Should(ContainElements("ledger", "auth"))
				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-ledger-and-auth.yaml"))
			})
		})
		Context("With audit enabled", func() {
			var brokerNatsDSNSettings *v1beta1.Settings
			BeforeEach(func() {
				stack.Spec.EnableAudit = true
				brokerNatsDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
			})
			JustBeforeEach(func() {
				Expect(Create(brokerNatsDSNSettings)).To(BeNil())
			})
			JustAfterEach(func() {
				Expect(Delete(brokerNatsDSNSettings)).To(Succeed())

			})
			It("Should create a topic", func() {
				Eventually(func() error {
					topic := &v1beta1.BrokerTopic{}
					return LoadResource("", fmt.Sprintf("%s-audit", stack.GetName()), topic)
				}).Should(Succeed())
			})
			It("Should adapt the Caddyfile", func() {
				cm := &corev1.ConfigMap{}
				Eventually(func(g Gomega) error {
					return LoadResource(stack.Name, "gateway", cm)
				}).Should(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-audit.yaml"))
			})
			//It("Should add env vars to the deployment", func() {
			//	Eventually(func(g Gomega) []corev1.EnvVar {
			//		d := &appsv1.Deployment{}
			//		g.Expect(LoadResource(stack.Name, "gateway", d)).To(Succeed())
			//		return d.Spec.Template.Spec.Containers[0].Env
			//	}).Should(ContainElements(
			//		brokerconfigurations.BrokerEnvVars(brokerConfiguration.Spec, stack.Name, "gateway"),
			//	))
			//})
		})
		Context("With otlp enabled", func() {
			var otelTracesDSNSetting *v1beta1.Settings
			JustBeforeEach(func() {
				otelTracesDSNSetting = settings.New(uuid.NewString(), "opentelemetry.traces.dsn", "grpc://collector", stack.Name)
				Expect(Create(otelTracesDSNSetting)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(otelTracesDSNSetting)).To(Succeed())
			})
			It("Should adapt the Caddyfile", func() {
				cm := &corev1.ConfigMap{}
				Eventually(func(g Gomega) string {
					g.Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
					return cm.Data["Caddyfile"]
				}).Should(MatchGoldenFile("gateway-controller", "configmap-with-opentelemetry.yaml"))
			})
			It("Should add env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					d := &appsv1.Deployment{}
					g.Expect(LoadResource(stack.Name, "gateway", d)).To(Succeed())
					return d.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(corev1.EnvVar{
					Name:  "OTEL_SERVICE_NAME",
					Value: "gateway",
				}))
			})
		})
	})
})
