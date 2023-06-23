package views

import (
	"fmt"
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintBalance(out io.Writer, balance shared.BalanceWithAssets) error {
	fctl.Section.Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name"), balance.Name})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.Println("Assets")
	if len(balance.Assets) == 0 {
		fctl.Println("No assets found.")
		return nil
	}
	tableData = pterm.TableData{}
	tableData = append(tableData, []string{"Asset", "Amount"})
	for asset, amount := range balance.Assets {
		tableData = append(tableData, []string{asset, fmt.Sprint(amount)})
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}
