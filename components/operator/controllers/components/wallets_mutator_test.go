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

var _ = Describe("Wallets controller", func() {
	mutator := NewWalletsMutator(GetClient(), GetScheme())
	WithMutator(mutator, func() {
		WithNewNamespace(func() {
			Context("When creating a wallets server", func() {
				var (
					wallets *componentsv1beta2.Wallets
				)
				BeforeEach(func() {
					wallets = &componentsv1beta2.Wallets{
						ObjectMeta: metav1.ObjectMeta{
							Name: "wallets",
						},
						Spec: componentsv1beta2.WalletsSpec{},
					}
					Expect(Create(wallets)).To(BeNil())
					Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeReady)).Should(Equal(metav1.ConditionTrue))
				})
				It("Should create a deployment", func() {
					Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeDeploymentReady)).Should(Equal(metav1.ConditionTrue))
					deployment := &appsv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{
							Name:      wallets.Name,
							Namespace: wallets.Namespace,
						},
					}
					Expect(Exists(deployment)()).To(BeTrue())
					Expect(deployment.OwnerReferences).To(HaveLen(1))
					Expect(deployment.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(wallets)))
				})
				It("Should create a service", func() {
					Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeServiceReady)).Should(Equal(metav1.ConditionTrue))
					service := &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:      wallets.Name,
							Namespace: wallets.Namespace,
						},
					}
					Expect(Exists(service)()).To(BeTrue())
					Expect(service.OwnerReferences).To(HaveLen(1))
					Expect(service.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(wallets)))
				})
				Context("Then enable ingress", func() {
					BeforeEach(func() {
						wallets.Spec.Ingress = &apisv1beta2.IngressSpec{
							Path: "/wallets",
							Host: "localhost",
						}
						Expect(Update(wallets)).To(BeNil())
					})
					It("Should create a ingress", func() {
						Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeIngressReady)).Should(Equal(metav1.ConditionTrue))
						ingress := &networkingv1.Ingress{
							ObjectMeta: metav1.ObjectMeta{
								Name:      wallets.Name,
								Namespace: wallets.Namespace,
							},
						}
						Expect(Exists(ingress)()).To(BeTrue())
						Expect(ingress.OwnerReferences).To(HaveLen(1))
						Expect(ingress.OwnerReferences).To(ContainElement(controllerutils.OwnerReference(wallets)))
					})
					Context("Then disabling ingress support", func() {
						BeforeEach(func() {
							Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeIngressReady)).
								Should(Equal(metav1.ConditionTrue))
							wallets.Spec.Ingress = nil
							Expect(Update(wallets)).To(BeNil())
							Eventually(ConditionStatus(wallets, apisv1beta2.ConditionTypeIngressReady)).
								Should(Equal(metav1.ConditionUnknown))
						})
						It("Should remove the ingress", func() {
							Eventually(NotFound(&networkingv1.Ingress{
								ObjectMeta: metav1.ObjectMeta{
									Name:      wallets.Name,
									Namespace: wallets.Namespace,
								},
							})).Should(BeTrue())
						})
					})
				})
			})
		})
	})
})
