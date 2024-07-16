package suite

import (
	"net/http"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {

	When("Trying to insert Hooks", func(){
		It("inserting a valid V2 Hook", func(){
			name := "ValidV2Hook"
			retry := true
			hookBodyParams := shared.V2HookBodyParams{
				Endpoint: "https://example.com",
				Events: []string{
					"ledger.committed_transactions",
				},
				Name: &name,
				Retry: &retry,
			}
	
			response , err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBodyParams,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
			newV2Hook := response.V2HookResponse.Data
			Expect(newV2Hook.Endpoint).To(Equal(hookBodyParams.Endpoint))
			Expect(newV2Hook.Events).To(Equal(hookBodyParams.Events))
			Expect(newV2Hook.Status).To(Equal(shared.V2HookStatusDisabled))
		})

		It("inserting an invalid V2 Hook", func(){
			name := "ValidV2Hook"
			retry := true
			hookBodyParams := shared.V2HookBodyParams{
				Endpoint: "https://example.com",
				Events: []string{
				},
				Name: &name,
				Retry: &retry,
			}

			_ , err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBodyParams,
			)
			Expect(err).To(HaveOccurred())
		
		})
		
		It("inserting a V2 Hook without endpoint", func(){
			name := "ValidV2Hook"
			retry := true
			hookBodyParams := shared.V2HookBodyParams{
				Endpoint: "",
				Events: []string{
				},
				Name: &name,
				Retry: &retry,
			}
				
			_ , err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBodyParams,
			)
			Expect(err).To(HaveOccurred())
		
		})

		It("inserting a V2 Hook with invalid secret", func(){
			name := "ValidV2Hook"
			retry := true
			secret := "invalid"
			hookBodyParams := shared.V2HookBodyParams{
				Endpoint: "",
				Events: []string{
				},
				Name: &name,
				Retry: &retry,
				Secret: &secret,
			}
				
			_ , err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBodyParams,
			)
			Expect(err).To(HaveOccurred())
		
		})

	})
	
	
})
