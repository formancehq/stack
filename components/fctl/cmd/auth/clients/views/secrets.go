package views

import (
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintSecrets(out io.Writer, secrets []shared.ClientSecret) error {
	fctl.Section.WithWriter(out).Println("Secrets :")

	return pterm.DefaultTable.
		WithWriter(out).
		WithHasHeader(true).
		WithData(fctl.Prepend(
			fctl.Map(secrets, func(secret shared.ClientSecret) []string {
				return []string{
					secret.ID, secret.Name, secret.LastDigits,
				}
			}),
			[]string{"ID", "Name", "Last digits"},
		)).
		Render()
}
