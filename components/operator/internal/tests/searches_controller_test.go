package tests_test

import (
	"fmt"
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
			stack                                 *v1beta1.Stack
			search                                *v1beta1.Search
			elasticSearchConfigurationHostSetting *v1beta1.Settings
			brokerKindSettings                    *v1beta1.Settings
			brokerNatsDSNSettings                 *v1beta1.Settings
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
			elasticSearchConfigurationHostSetting = settings.New(uuid.NewString(),
				"elasticsearch.host", "localhost", stack.Name)
			brokerKindSettings = settings.New(uuid.NewString(),
				"broker.kind", "nats", stack.Name)
			brokerNatsDSNSettings = settings.New(uuid.NewString(), "broker.nats.dsn", "nats://localhost:1234", stack.Name)
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(elasticSearchConfigurationHostSetting)).To(Succeed())
			Expect(Create(brokerKindSettings)).To(Succeed())
			Expect(Create(brokerNatsDSNSettings)).To(Succeed())
			Expect(Create(search)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(search)).To(Succeed())
			Expect(Delete(elasticSearchConfigurationHostSetting)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(brokerKindSettings)).To(Succeed())
			Expect(Delete(brokerNatsDSNSettings)).To(Succeed())
		})
		It("Should create a stream processor", func() {
			streamProcessor := &v1beta1.StreamProcessor{}
			Eventually(func() error {
				return LoadResource(stack.Name, fmt.Sprintf("%s-stream-processor", stack.Name), streamProcessor)
			}).Should(Succeed())
			Expect(streamProcessor).To(BeControlledBy(search))
		})
		Context("Then when creating a SearchBatchingConfiguration object", func() {
			var searchBatchingCountSettings *v1beta1.Settings
			JustBeforeEach(func() {
				searchBatchingCountSettings = settings.New(uuid.NewString(), "search.batching.count", "10", stack.Name)
				Expect(Create(searchBatchingCountSettings)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(searchBatchingCountSettings)).To(Succeed())
			})
			It("Should update the stream processor with the new batching configuration", func() {
				streamProcessor := &v1beta1.StreamProcessor{}
				Eventually(func(g Gomega) v1beta1.Batching {
					g.Expect(LoadResource(stack.Name, fmt.Sprintf("%s-stream-processor", stack.Name), streamProcessor)).To(Succeed())
					g.Expect(streamProcessor.Spec.Batching).NotTo(BeNil())
					return *streamProcessor.Spec.Batching
				}).Should(Equal(v1beta1.Batching{
					Count: 10,
				}))
			})
		})
	})
})
