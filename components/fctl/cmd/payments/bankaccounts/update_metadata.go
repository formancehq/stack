package bankaccounts

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateMetadataStore struct {
	Success bool `json:"success"`
}
type UpdateMetadataController struct {
	PaymentsVersion versions.Version

	store *UpdateMetadataStore
}

func (c *UpdateMetadataController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*UpdateMetadataStore] = (*UpdateMetadataController)(nil)

func NewUpdateMetadataStore() *UpdateMetadataStore {
	return &UpdateMetadataStore{}
}

func NewUpdateMetadataController() *UpdateMetadataController {
	return &UpdateMetadataController{
		store: NewUpdateMetadataStore(),
	}
}

func NewUpdateMetadataCommand() *cobra.Command {
	c := NewUpdateMetadataController()
	return fctl.NewCommand("update-metadata <bankAccountID> [<key>=<value>...]",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Set metadata on bank account"),
		fctl.WithAliases("um", "update-meta"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*UpdateMetadataStore](c),
	)
}

func (c *UpdateMetadataController) GetStore() *UpdateMetadataStore {
	return c.store
}

func (c *UpdateMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("bank accounts are only supported in >= v1.0.0")
	}

	metadata, err := fctl.ParseMetadata(args[1:])
	if err != nil {
		return nil, err
	}

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	bankAccountID := args[0]

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to set a metadata on bank account '%s'", bankAccountID) {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.UpdateBankAccountMetadataRequest{
		UpdateBankAccountMetadataRequest: shared.UpdateBankAccountMetadataRequest{
			Metadata: metadata,
		},
		BankAccountID: bankAccountID,
	}

	response, err := client.Payments.UpdateBankAccountMetadata(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 204

	return c, nil
}

func (c *UpdateMetadataController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata added!")
	return nil
}
