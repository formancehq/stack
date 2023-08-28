package views

import (
	"fmt"
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintHold(out io.Writer, hold shared.ExpandedDebitHold) error {
	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(hold.ID)})
	tableData = append(tableData, []string{pterm.LightCyan("Wallet ID"), hold.WalletID})
	tableData = append(tableData, []string{pterm.LightCyan("Original amount"), fmt.Sprint(hold.OriginalAmount)})
	tableData = append(tableData, []string{pterm.LightCyan("Remaining"), fmt.Sprint(hold.Remaining)})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(out, hold.Metadata)
}
