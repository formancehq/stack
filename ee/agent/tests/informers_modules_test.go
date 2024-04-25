package tests

import (
	"context"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("Informer modules", func() {
	var (
		membershipClientMock *internal.MembershipClientMock
		restMapper           meta.RESTMapper
		err                  error
	)

	BeforeEach(func() {
		membershipClientMock = internal.NewMembershipClientMock()
		restMapper, err = internal.CreateRestMapper(restConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	When("a module is created on the cluster", func() {
		var (
			modules map[schema.GroupVersionKind]*unstructured.Unstructured
		)
		BeforeEach(func() {
			modules = map[schema.GroupVersionKind]*unstructured.Unstructured{}
			for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
				object := reflect.New(rtype).Interface()
				if _, ok := object.(v1beta1.Module); !ok {
					continue
				}

				restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
				Expect(err).ToNot(HaveOccurred())
				var (
					unstructuredObj *unstructured.Unstructured
					resource        string
				)

				resource = restMapping.Resource.Resource
				name := uuid.NewString()

				unstructuredObj = &unstructured.Unstructured{
					Object: map[string]interface{}{},
				}
				unstructuredObj.Object["apiVersion"] = gvk.GroupVersion().String()
				unstructuredObj.Object["kind"] = gvk.Kind
				unstructuredObj.Object["metadata"] = map[string]interface{}{
					"name": name,
				}

				By("Creating the module", func() {
					Expect(k8sClient.Post().
						Resource(resource).
						Body(unstructuredObj).
						Name(name).Do(context.Background()).Error()).To(Succeed())

					DeferCleanup(func() {
						Expect(k8sClient.Delete().Resource(resource).Name(name).Do(context.Background()).Error()).To(Succeed())
					})

				})

				By("Loading then updating status", func() {
					Eventually(func() error {
						return LoadResource(resource, name, unstructuredObj)
					}).Should(Succeed())

					/**
						Those 2 a reset by LoadResource to empty
					**/
					unstructuredObj.Object["apiVersion"] = gvk.GroupVersion().String()
					unstructuredObj.Object["kind"] = gvk.Kind
					/** **/

					unstructuredObj.Object["status"] = map[string]interface{}{
						"info": uuid.NewString(),
					}

					Expect(k8sClient.Put().
						Resource(resource).
						SubResource("status").
						Name(name).
						Body(unstructuredObj).
						Do(context.Background()).
						Error()).To(Succeed())
				})

				modules[gvk] = unstructuredObj
			}
		})
		It("Should have been created", func() {
			for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
				object := reflect.New(rtype).Interface()
				if _, ok := object.(v1beta1.Module); !ok {
					continue
				}

				restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
				Expect(err).ToNot(HaveOccurred())

				Expect(true).To(BeTrue())
				Eventually(func() error {
					return LoadResource(restMapping.Resource.Resource, modules[gvk].GetName(), modules[gvk])
				}).Should(Succeed())
			}
		})
		When("Listening to all modules", func() {
			BeforeEach(func() {
				dynamicClient, err := dynamic.NewForConfig(restConfig)
				Expect(err).ToNot(HaveOccurred())
				factory := internal.NewDynamicSharedInformerFactory(dynamicClient)
				Expect(internal.CreateModulesInformers(factory, restMapper, logging.Testing(), membershipClientMock)).ToNot(HaveOccurred())

				stopCh := make(chan struct{})
				factory.Start(stopCh)
				DeferCleanup(func() {
					close(stopCh)
				})
			})
			It("Should have sent ModuleStatusChanged", func() {
				for gvk, module := range modules {
					Eventually(func(g Gomega) bool {
						for _, message := range membershipClientMock.GetMessages() {
							if msg := message.GetModuleStatusChanged(); msg != nil &&
								msg.Vk.Kind == gvk.Kind &&
								msg.Vk.Version == gvk.Version &&
								msg.ClusterName == module.GetName() {

								status, _, _ := unstructured.NestedMap(module.Object, "status")

								g.Expect(msg.Status.AsMap()["info"]).ToNot(BeNil())
								g.Expect(msg.Status.AsMap()["ready"]).To(BeFalse())

								for k, value := range status {
									g.Expect(msg.Status.AsMap()[k]).To(Equal(value))
								}
								return true
							}
						}
						return false
					}).Should(BeTrue())
				}
			})
		})
	})
})
