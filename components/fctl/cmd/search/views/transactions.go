package views

import (
	"fmt"
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
)

func DisplayTransactions(out io.Writer, txs []map[string]interface{}) error {
	tableData := make([][]string, 0)
	for _, tx := range txs {
		tableData = append(tableData, []string{
			tx["ledger"].(string),
			fmt.Sprint(tx["id"].(float64)),
			tx["reference"].(string),
			tx["timestamp"].(string),
		})
	}
	tableData = fctl.Prepend(tableData, []string{"Ledger", "ID", "Reference", "Date"})

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(out).
		WithData(tableData).
		Render()
}
