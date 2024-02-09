package tests_test

import (
	"fmt"

	"github.com/formancehq/operator/internal/resources/settings"

	. "github.com/formancehq/operator/internal/tests/internal"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("BrokerTopicController", func() {
	Context("When creating a BrokerTopic", func() {
		var (
			stack             *v1beta1.Stack
			brokerTopic       *v1beta1.BrokerTopic
			brokerDSNSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			brokerDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
			Expect(Create(brokerDSNSettings)).To(BeNil())
			brokerTopic = &v1beta1.BrokerTopic{
				ObjectMeta: v1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: v1beta1.BrokerTopicSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Service: "ledger",
				},
			}
			// notes(gfyrag): add a "fake" owner reference to prevent automatic deletion
			Expect(controllerutil.SetOwnerReference(brokerDSNSettings, brokerTopic, GetScheme())).To(Succeed())
			Expect(Create(brokerTopic)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
			Expect(client.IgnoreNotFound(Delete(brokerTopic))).To(Succeed())
			Expect(Delete(brokerDSNSettings)).To(Succeed())
		})
		It("Should be set to ready status", func() {
			t := &v1beta1.BrokerTopic{}
			Eventually(func(g Gomega) bool {
				g.Expect(Get(core.GetResourceName(brokerTopic.Name), t)).To(Succeed())
				return t.Status.Ready
			}).Should(BeTrue())
		})
		Context("Then updating removing all owner references", func() {
			BeforeEach(func() {
				Eventually(func(g Gomega) bool {
					t := &v1beta1.BrokerTopic{}
					g.Expect(Get(core.GetResourceName(brokerTopic.Name), t)).To(Succeed())
					return t.Status.Ready
				}).Should(BeTrue())

				patch := client.MergeFrom(brokerTopic.DeepCopy())
				Expect(controllerutil.RemoveOwnerReference(brokerDSNSettings, brokerTopic, GetScheme())).To(Succeed())
				Expect(Patch(brokerTopic, patch)).To(Succeed())
			})
			It("Should trigger the deletion of the brokerTopic object", func() {
				Eventually(func(g Gomega) error {
					return LoadResource("", brokerTopic.Name, brokerTopic)
				}).Should(BeNotFound())
			})
			It("Should create a brokerTopic deletion job", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, fmt.Sprintf("%s-delete-topic", brokerTopic.GetUID()), &batchv1.Job{})
				}).Should(Succeed())
			})
		})
	})
})
