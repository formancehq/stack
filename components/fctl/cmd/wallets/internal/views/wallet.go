package views

import (
	"fmt"
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintWallet(out io.Writer, wallet shared.WalletWithBalances) error {
	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(wallet.ID)})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), wallet.Name})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.WithWriter(out).Println("Balances")
	if len(wallet.Balances.Main.Assets) == 0 {
		fmt.Fprintln(out, "No balances found.")
		return nil
	}
	tableData = pterm.TableData{}
	tableData = append(tableData, []string{"Asset", "Amount"})
	for asset, amount := range wallet.Balances.Main.Assets {
		tableData = append(tableData, []string{asset, fmt.Sprint(amount)})
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fmt.Fprintln(out, "")

	return nil
}
func PrintWalletWithMetadata(out io.Writer, wallet shared.WalletWithBalances) error {
	err := PrintWallet(out, wallet)
	if err != nil {
		return err
	}

	return fctl.PrintMetadata(out, wallet.Metadata)
}
