package components

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/controllerutils"
	. "github.com/formancehq/operator/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Payments controller", func() {
	mutator := NewPaymentsMutator(GetClient(), GetScheme())
	WithMutator(mutator, func() {
		WithNewNamespace(func() {
			Context("When creating a payments server", func() {
				var (
					payments *componentsv1beta2.Payments
				)
				BeforeEach(func() {
					payments = &componentsv1beta2.Payments{
						ObjectMeta: metav1.ObjectMeta{
							Name: "payments",
						},
						Spec: componentsv1beta2.PaymentsSpec{
							Collector: &componentsv1beta2.CollectorConfig{
								KafkaConfig: NewDumpKafkaConfig(),
								Topic:       "xxx",
							},
							Postgres: componentsv1beta2.PostgresConfigCreateDatabase{
								PostgresConfigWithDatabase: apisv1beta2.PostgresConfigWithDatabase{
									Database:       "payments",
									PostgresConfig: NewDumpPostgresConfig(),
								},
								CreateDatabase: false,
							}},
					}
					Expect(Create(payments)).To(BeNil())
					Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeReady)).Should(Equal(metav1.ConditionTrue))
				})
				It("Should create a deployment", func() {
					Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeDeploymentReady)).Should(Equal(metav1.ConditionTrue))
					deployment := &appsv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{
							Name:      payments.Name,
							Namespace: payments.Namespace,
						},
					}
					Expect(Exists(deployment)()).To(BeTrue())
					Expect(deployment.OwnerReferences).To(HaveLen(1))
					Expect(deployment.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(payments)))
				})
				It("Should create a service", func() {
					Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeServiceReady)).Should(Equal(metav1.ConditionTrue))
					service := &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:      payments.Name,
							Namespace: payments.Namespace,
						},
					}
					Expect(Exists(service)()).To(BeTrue())
					Expect(service.OwnerReferences).To(HaveLen(1))
					Expect(service.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(payments)))
				})
				Context("Then enable ingress", func() {
					BeforeEach(func() {
						payments.Spec.Ingress = &apisv1beta2.IngressSpec{
							Path: "/payments",
							Host: "localhost",
						}
						Expect(Update(payments)).To(BeNil())
					})
					It("Should create a ingress", func() {
						Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeIngressReady)).Should(Equal(metav1.ConditionTrue))
						ingress := &networkingv1.Ingress{
							ObjectMeta: metav1.ObjectMeta{
								Name:      payments.Name,
								Namespace: payments.Namespace,
							},
						}
						Expect(Exists(ingress)()).To(BeTrue())
						Expect(ingress.OwnerReferences).To(HaveLen(1))
						Expect(ingress.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(payments)))
					})
					Context("Then disabling ingress support", func() {
						BeforeEach(func() {
							Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeIngressReady)).
								Should(Equal(metav1.ConditionTrue))
							payments.Spec.Ingress = nil
							Expect(Update(payments)).To(BeNil())
							Eventually(ConditionStatus(payments, apisv1beta2.ConditionTypeIngressReady)).
								Should(Equal(metav1.ConditionUnknown))
						})
						It("Should remove the ingress", func() {
							Eventually(NotFound(&networkingv1.Ingress{
								ObjectMeta: metav1.ObjectMeta{
									Name:      payments.Name,
									Namespace: payments.Namespace,
								},
							})).Should(BeTrue())
						})
					})
				})
			})
		})
	})
})
