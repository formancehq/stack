package tests

import (
	"context"
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/components/agent/internal/generated"
	. "github.com/formancehq/stack/components/agent/tests/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"net/url"
	"reflect"
)

var _ = Describe("Membership listener", func() {
	var (
		orders     internal.OrdersChan
		clientInfo internal.ClientInfo
	)
	BeforeEach(func() {
		orders = make(internal.OrdersChan)
		clientInfo = internal.ClientInfo{
			BaseUrl: &url.URL{},
		}
		listener := internal.NewMembershipListener(k8sClient, clientInfo, mapper, orders)
		done := make(chan struct{})
		DeferCleanup(func() {
			<-done
		})
		go func() {
			defer close(done)
			listener.Start(context.Background())
		}()

		DeferCleanup(func() {
			close(orders)
		})
	})
	Context("When sending an existing stack from membership", func() {
		It("Should create all required crds cluster side", func() {
			membershipStack := &generated.Stack{
				ClusterName: "stack0",
				AuthConfig: &generated.AuthConfig{
					ClientId:     "clientid",
					ClientSecret: "clientsecret",
					Issuer:       "http://example.net",
				},
			}
			orders <- &generated.Order{
				Message: &generated.Order_ExistingStack{
					ExistingStack: membershipStack,
				},
			}

			stack := &v1beta1.Stack{}
			Eventually(func() error {
				return LoadResource("Stacks", "stack0", stack)
			}).Should(BeNil())

			auth := &v1beta1.Auth{}
			Eventually(func() error {
				return LoadResource("Auths", "stack0", auth)
			}).Should(BeNil())
			Expect(auth).To(BeOwnedBy(stack))
			Expect(auth).To(TargetStack(stack))
			Expect(auth.Spec.DelegatedOIDCServer).NotTo(BeNil())
			Expect(auth.Spec.DelegatedOIDCServer.ClientSecret).To(Equal(membershipStack.AuthConfig.ClientSecret))
			Expect(auth.Spec.DelegatedOIDCServer.ClientID).To(Equal(membershipStack.AuthConfig.ClientId))
			Expect(auth.Spec.DelegatedOIDCServer.Issuer).To(Equal(membershipStack.AuthConfig.Issuer))

			gateway := &v1beta1.Gateway{}
			Eventually(func() error {
				return LoadResource("Gateways", "stack0", gateway)
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
				if gvk.Kind == "Reconciliation" || gvk.Kind == "Stargate" { // EE modules, not actually enabled
					continue
				}

				restMapping, err := mapper.RESTMapping(gvk.GroupKind())
				Expect(err).To(BeNil())

				u := &unstructured.Unstructured{}
				Eventually(func() error {
					return LoadResource(restMapping.Resource.Resource, "stack0", u)
				}).Should(Succeed())
				Expect(u).To(BeOwnedBy(stack))
				Expect(u).To(TargetStack(stack))
			}
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
