package tests_test

import (
	"fmt"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("DatabaseController", func() {
	Context("When creating a Database", func() {
		var (
			stack            *v1beta1.Stack
			database         *v1beta1.Database
			databaseSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(BeNil())
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
			Expect(Create(databaseSettings)).Should(Succeed())
			database = &v1beta1.Database{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.DatabaseSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Service: "ledger",
				},
			}
			Expect(Create(database)).To(Succeed())
		})
		AfterEach(func() {
			Expect(client.IgnoreNotFound(Delete(database))).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		shouldBeReady := func() {
			d := &v1beta1.Database{}
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", database.Name, d)).To(Succeed())
				return d.Status.Ready
			}).Should(BeTrue())
		}
		It("Should be set to ready status", shouldBeReady)
		Context("Then when deleting the Database object", func() {
			BeforeEach(func() {
				shouldBeReady()
				Expect(Delete(database)).To(Succeed())
			})
			It("Should create a deletion job", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, fmt.Sprintf("%s-drop-database", database.Spec.Service), &batchv1.Job{})
				}).Should(Succeed())
			})
			It("Should eventually be deleted", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, database.Name, &v1beta1.Database{})
				}).Should(BeNotFound())
			})
		})
		Context("Then when updating the DatabaseConfiguration object", func() {
			BeforeEach(func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", database.Name, database)).To(Succeed())
					return database.Status.Ready
				}).Should(BeTrue())

				patch := client.MergeFrom(databaseSettings.DeepCopy())
				databaseSettings.Spec.Value = "postgresql://xxx"
				Expect(Patch(databaseSettings, patch)).To(Succeed())
			})
			It("Should declare the Database object as out of sync", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", database.Name, database)).To(Succeed())

					return database.Status.OutOfSync
				}).Should(BeTrue())
			})
		})
	})
})
