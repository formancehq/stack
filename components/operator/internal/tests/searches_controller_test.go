package tests_test

import (
	"fmt"

	"github.com/formancehq/operator/internal/core"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SearchesController", func() {
	Context("When creating a Search object", func() {
		var (
			stack                   *v1beta1.Stack
			search                  *v1beta1.Search
			elasticSearchDSNSetting *v1beta1.Settings
			brokerDSNSettings       *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			search = &v1beta1.Search{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.SearchSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			elasticSearchDSNSetting = settings.New(uuid.NewString(), "elasticsearch.dsn", "https://localhost", stack.Name)
			brokerDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(elasticSearchDSNSetting)).To(Succeed())
			Expect(Create(brokerDSNSettings)).To(Succeed())
			Expect(Create(search)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(search)).To(Succeed())
			Expect(Delete(elasticSearchDSNSetting)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(brokerDSNSettings)).To(Succeed())
		})
		It("Should add an owner reference on the stack", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", search.Name, search)).To(Succeed())
				reference, err := core.HasOwnerReference(TestContext(), stack, search)
				g.Expect(err).To(BeNil())
				return reference
			}).Should(BeTrue())
		})
		It("Should create a Benthos", func() {
			benthos := &v1beta1.Benthos{}
			Eventually(func() error {
				return LoadResource(stack.Name, fmt.Sprintf("%s-benthos", stack.Name), benthos)
			}).Should(Succeed())
			Expect(benthos).To(BeControlledBy(search))
		})
		Context("Then when creating a SearchBatchingConfiguration object", func() {
			var searchBatchingCountSettings *v1beta1.Settings
			JustBeforeEach(func() {
				searchBatchingCountSettings = settings.New(uuid.NewString(), "search.batching", "count=10", stack.Name)
				Expect(Create(searchBatchingCountSettings)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(searchBatchingCountSettings)).To(Succeed())
			})
			It("Should update the Benthos with the new batching configuration", func() {
				benthos := &v1beta1.Benthos{}
				Eventually(func(g Gomega) v1beta1.Batching {
					g.Expect(LoadResource(stack.Name, fmt.Sprintf("%s-benthos", stack.Name), benthos)).To(Succeed())
					g.Expect(benthos.Spec.Batching).NotTo(BeNil())
					return *benthos.Spec.Batching
				}).Should(Equal(v1beta1.Batching{
					Count: 10,
				}))
			})
		})
	})
})
