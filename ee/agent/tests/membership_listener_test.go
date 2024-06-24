package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/components/agent/internal/generated"
	. "github.com/formancehq/stack/components/agent/tests/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("Membership listener", func() {
	var (
		membershipClient *internal.MembershipClientMock
		clientInfo       internal.ClientInfo
	)
	BeforeEach(func() {
		membershipClient = internal.NewMembershipClientMock()
		clientInfo = internal.ClientInfo{
			BaseUrl: &url.URL{},
		}
		listener := internal.NewMembershipListener(internal.NewDefaultK8SClient(k8sClient), clientInfo, mapper, membershipClient)
		done := make(chan struct{})
		DeferCleanup(func() {
			<-done
		})
		go func() {
			defer close(done)
			listener.Start(context.Background())
		}()

		DeferCleanup(func() {
			close(membershipClient.Orders())
		})
	})
	Context("When sending an existing stack from membership", func() {
		var (
			membershipStack *generated.Stack
			stack           *v1beta1.Stack
		)
		BeforeEach(func() {
			stackName := uuid.NewString()
			wrong := uuid.NewString()
			By("Creating a wrong client", func() {
				client := &v1beta1.AuthClient{
					ObjectMeta: metav1.ObjectMeta{
						Name: wrong,
						Labels: map[string]string{
							"formance.com/created-by-agent": "true",
							"formance.com/stack":            stackName,
						},
					},
					Spec: v1beta1.AuthClientSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stackName,
						},
						ID:     wrong,
						Public: true,
					},
				}
				Expect(k8sClient.Post().Resource("AuthClients").Body(client).Do(context.Background()).Error()).To(BeNil())
			})

			By("Creating a stack", func() {
				modules := make([]*generated.Module, 0)
				for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
					object := reflect.New(rtype).Interface()
					if _, ok := object.(v1beta1.Module); !ok {
						continue
					}

					if gvk.Kind == "Stargate" {
						continue
					}

					modules = append(modules, &generated.Module{
						Name: gvk.Kind,
					})
				}

				membershipStack = &generated.Stack{
					ClusterName: stackName,
					AuthConfig: &generated.AuthConfig{
						ClientId:     "clientid",
						ClientSecret: "clientsecret",
						Issuer:       "http://example.net",
					},
					AdditionalLabels: map[string]string{
						"foo":     "bar",
						"foo.foo": "bar",
						"foo-foo": "bar",
					},
					AdditionalAnnotations: map[string]string{
						"foo":     "bar",
						"foo.foo": "bar",
						"foo-foo": "bar",
						"foo_foo": "@bar",
					},
					StaticClients: []*generated.AuthClient{
						{
							Id:     "clientid1",
							Public: true,
						},
						{
							Id:     "clientid2",
							Public: true,
						},
					},
					Modules: modules,
				}
				membershipClient.Orders() <- &generated.Order{
					Message: &generated.Order_ExistingStack{
						ExistingStack: membershipStack,
					},
				}
				stack = &v1beta1.Stack{}
				Eventually(func() error {
					return LoadResource("Stacks", membershipStack.ClusterName, stack)
				}).Should(BeNil())
			})
		})
		It("Should have sync auth client", func() {
			clients := &unstructured.UnstructuredList{}
			Eventually(func(g Gomega) []unstructured.Unstructured {
				g.Expect(k8sClient.Get().Resource("AuthClients").VersionedParams(&metav1.ListOptions{
					LabelSelector: "formance.com/created-by-agent=true,formance.com/stack=" + membershipStack.ClusterName,
				}, scheme.ParameterCodec).Do(context.Background()).Into(clients)).To(Succeed())
				return clients.Items
			}).Should(HaveLen(2))
		})
		When("Having an @ in labels", func() {
			var patch []byte
			var err error
			BeforeEach(func() {
				patch, err = json.Marshal(map[string]any{
					"metadata": map[string]any{
						"labels": map[string]any{
							"formance.com/owner-email": "example@example.net",
						},
					},
				})
				Expect(err).To(BeNil())

			})
			It("Should through an error", func() {
				err := k8sClient.Patch(types.MergePatchType).
					Resource("Stacks").
					Name(stack.Name).
					Body(patch).
					Do(context.Background()).
					Error()

				Expect(err).To(HaveOccurred())
			})
		})
		It("Should have additional labels", func() {
			Expect(stack.Labels).To(HaveKeyWithValue("formance.com/foo", "bar"))
			Expect(stack.Labels).To(HaveKeyWithValue("formance.com/foo.foo", "bar"))
			Expect(stack.Labels).To(HaveKeyWithValue("formance.com/foo-foo", "bar"))
		})
		It("Should have additional annotations", func() {
			Expect(stack.Annotations).To(HaveKeyWithValue("formance.com/foo", "bar"))
			Expect(stack.Annotations).To(HaveKeyWithValue("formance.com/foo.foo", "bar"))
			Expect(stack.Annotations).To(HaveKeyWithValue("formance.com/foo-foo", "bar"))
			Expect(stack.Annotations).To(HaveKeyWithValue("formance.com/foo_foo", "@bar"))
		})
		It("Should create all required crds cluster side", func() {
			auth := &v1beta1.Auth{}
			Eventually(func() error {
				return LoadResource("Auths", membershipStack.ClusterName, auth)
			}).Should(BeNil())
			Expect(auth).To(BeOwnedBy(stack))
			Expect(auth).To(TargetStack(stack))
			Expect(auth.Spec.DelegatedOIDCServer).NotTo(BeNil())
			Expect(auth.Spec.DelegatedOIDCServer.ClientSecret).To(Equal(membershipStack.AuthConfig.ClientSecret))
			Expect(auth.Spec.DelegatedOIDCServer.ClientID).To(Equal(membershipStack.AuthConfig.ClientId))
			Expect(auth.Spec.DelegatedOIDCServer.Issuer).To(Equal(membershipStack.AuthConfig.Issuer))

			gateway := &v1beta1.Gateway{}
			Eventually(func() error {
				return LoadResource("Gateways", membershipStack.ClusterName, gateway)
			}).Should(BeNil())
			Expect(gateway).To(BeOwnedBy(stack))
			Expect(gateway).To(TargetStack(stack))
			Expect(gateway.Spec.Ingress).NotTo(BeNil())
			Expect(gateway.Spec.Ingress.Host).To(Equal(fmt.Sprintf("%s.%s", stack.GetName(), clientInfo.BaseUrl.Host)))
			Expect(gateway.Spec.Ingress.Scheme).To(Equal(clientInfo.BaseUrl.Scheme))

			for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
				object := reflect.New(rtype).Interface()
				if _, ok := object.(v1beta1.Module); !ok {
					continue
				}

				if gvk.Kind == "Stargate" {
					continue
				}

				restMapping, err := mapper.RESTMapping(gvk.GroupKind())
				Expect(err).To(BeNil())

				u := &unstructured.Unstructured{}
				Eventually(func() error {
					return LoadResource(restMapping.Resource.Resource, membershipStack.ClusterName, u)
				}).Should(Succeed())
				Expect(u).To(BeOwnedBy(stack))
				Expect(u).To(TargetStack(stack))
			}
		})
		When("removing modules", func() {
			var (
				modulesToRemove map[string]struct{}
			)
			BeforeEach(func() {

				modulesToRemove = map[string]struct{}{}
				modulesToRemove["Webhooks"] = struct{}{}
				modulesToRemove["Search"] = struct{}{}

				membershipStack.Modules = collectionutils.Filter(membershipStack.Modules, func(module *generated.Module) bool {
					_, exist := modulesToRemove[module.Name]
					return !exist
				})

				membershipClient.Orders() <- &generated.Order{
					Message: &generated.Order_ExistingStack{
						ExistingStack: membershipStack,
					},
				}
			})
			It("modules should be removed", func() {
				for moduleName := range modulesToRemove {
					Eventually(func(g Gomega) error {
						u := &unstructured.Unstructured{}
						return LoadResource(moduleName, membershipStack.ClusterName, u)
					}).Should(HaveOccurred())
				}
			})
		})
		Context("then when disabling the stack", func() {
			BeforeEach(func() {
				membershipClient.Orders() <- &generated.Order{
					Message: &generated.Order_DisabledStack{
						DisabledStack: &generated.DisabledStack{
							ClusterName: membershipStack.ClusterName,
						},
					},
				}
			})
			shouldBeDisabled := func() {
				stack := &v1beta1.Stack{}
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("Stacks", membershipStack.ClusterName, stack)).To(Succeed())
					return stack.Spec.Disabled
				}).Should(BeTrue())
			}
			It("Should disable the stack on the cluster", shouldBeDisabled)
			Context("Then re enabling the stack", func() {
				BeforeEach(func() {
					shouldBeDisabled()
					membershipClient.Orders() <- &generated.Order{
						Message: &generated.Order_EnabledStack{
							EnabledStack: &generated.EnabledStack{
								ClusterName: membershipStack.ClusterName,
							},
						},
					}
				})
				It("Should enable the stack on the cluster", func() {
					stack := &v1beta1.Stack{}
					Eventually(func(g Gomega) bool {
						g.Expect(LoadResource("Stacks", membershipStack.ClusterName, stack)).To(Succeed())
						return stack.Spec.Disabled
					}).Should(BeFalse())
				})
			})
		})
	})
})

func LoadResource(resource, name string, into runtime.Object) error {
	return k8sClient.Get().
		Resource(resource).
		Name(name).
		Do(context.Background()).
		Into(into)
}
