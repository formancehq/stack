package formance_com_test

import (
	"fmt"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/controller/testing"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/httpapis"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
					g.Expect(gateway.Status.SyncHTTPAPIs).To(ContainElements("ledger", "auth"))
					return gateway.Status.AuthEnabled
				}).Should(BeTrue())
				cm := &corev1.ConfigMap{}
				Expect(LoadResource(stack.Name, "gateway", cm)).To(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-ledger-and-auth.yaml"))
			})
		})
		Context("With audit enabled", func() {
			var brokerConfiguration *v1beta1.BrokerConfiguration
			BeforeEach(func() {
				stack.Spec.EnableAudit = true
				brokerConfiguration = &v1beta1.BrokerConfiguration{
					ObjectMeta: metav1.ObjectMeta{
						Name: uuid.NewString(),
						Labels: map[string]string{
							core.StackLabel: stack.Name,
						},
					},
					Spec: v1beta1.BrokerConfigurationSpec{
						Nats: &v1beta1.NatsConfig{
							URL:      "nats://localhost:4321",
							Replicas: 10,
						},
					},
				}
			})
			JustBeforeEach(func() {
				Expect(Create(brokerConfiguration)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(brokerConfiguration)).To(Succeed())
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
			It("Should add env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					d := &appsv1.Deployment{}
					g.Expect(LoadResource(stack.Name, "gateway", d)).To(Succeed())
					return d.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(
					brokerconfigurations.BrokerEnvVars(brokerConfiguration.Spec, stack.Name, "gateway"),
				))
			})
		})
		Context("With otlp enabled", func() {
			var openTelemetryConfiguration *v1beta1.OpenTelemetryConfiguration
			BeforeEach(func() {
				openTelemetryConfiguration = &v1beta1.OpenTelemetryConfiguration{
					ObjectMeta: metav1.ObjectMeta{
						Name: uuid.NewString(),
						Labels: map[string]string{
							core.StackLabel: stack.Name,
						},
					},
					Spec: v1beta1.OpenTelemetryConfigurationSpec{
						Traces: &v1beta1.TracesSpec{
							Otlp: &v1beta1.OtlpSpec{
								Endpoint: "otlp",
								Port:     4317,
								Insecure: false,
								Mode:     "grpc",
							},
						},
					},
				}
			})
			JustBeforeEach(func() {
				Expect(Create(openTelemetryConfiguration)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(openTelemetryConfiguration)).To(Succeed())
			})
			It("Should adapt the Caddyfile", func() {
				cm := &corev1.ConfigMap{}
				Eventually(func(g Gomega) error {
					return LoadResource(stack.Name, "gateway", cm)
				}).Should(Succeed())
				Expect(cm.Data["Caddyfile"]).To(
					MatchGoldenFile("gateway-controller", "configmap-with-opentelemetry.yaml"))
			})
			It("Should add env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					d := &appsv1.Deployment{}
					g.Expect(LoadResource(stack.Name, "gateway", d)).To(Succeed())
					return d.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(
					opentelemetryconfigurations.GetEnvVars(openTelemetryConfiguration, "gateway"),
				))
			})
		})
	})
})
