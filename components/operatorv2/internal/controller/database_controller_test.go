package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("DatabaseController", func() {
	Context("When creating a Database", func() {
		var (
			stack                 *v1beta1.Stack
			database              *v1beta1.Database
			databaseConfiguration *v1beta1.DatabaseConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/service": "any",
						"formance.com/stack":   stack.Name,
					},
				},
				Spec: v1beta1.DatabaseConfigurationSpec{},
			}
			Expect(Create(databaseConfiguration)).To(Succeed())
			database = &v1beta1.Database{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.DatabaseSpec{
					Stack:   stack.Name,
					Service: "ledger",
				},
			}
			Expect(Create(database)).To(Succeed())
		})
		AfterEach(func() {
			Expect(client.IgnoreNotFound(Delete(database))).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should be set to ready status", func() {
			d := &v1beta1.Database{}
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", database.Name, d)).To(Succeed())
				return d.Status.Ready
			}).Should(BeTrue())
		})
		Context("Then when deleting the Database object", func() {
			BeforeEach(func() {
				Expect(Delete(database)).To(Succeed())
			})
			It("Should drop the underlying object", func() {
				Eventually(LoadResource(stack.Name, database.Name, &v1beta1.Database{})).
					Should(BeNotFound())
			})
		})
	})
})
