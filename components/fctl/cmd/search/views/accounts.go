package views

import (
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
)

func DisplayAccounts(out io.Writer, accounts []map[string]interface{}) error {
	tableData := make([][]string, 0)
	for _, account := range accounts {
		tableData = append(tableData, []string{
			// TODO: Missing property 'ledger' on api response
			//account["ledger"].(string),
			account["address"].(string),
		})
	}
	tableData = fctl.Prepend(tableData, []string{ /*"Ledger",*/ "Address"})

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(out).
		WithData(tableData).
		Render()
}
