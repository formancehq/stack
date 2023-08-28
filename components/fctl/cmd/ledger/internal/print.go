package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
)

func printCommonInformation(
	out io.Writer,
	transaction *ExportTransaction,
) error {
	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(transaction.Info.ID)})
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), transaction.Info.Reference})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), transaction.Info.Date.Format(time.RFC3339)})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	fmt.Fprintln(out, "")
	fctl.Section.WithWriter(out).Println("Postings")
	tableData = pterm.TableData{}
	tableData = append(tableData, []string{"Source", "Destination", "Asset", "Amount"})
	for _, posting := range transaction.Postings {
		tableData = append(tableData, []string{
			posting.Source, posting.Destination, posting.Asset, fmt.Sprint(posting.Amount),
		})
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

func PrintExpandedTransaction(out io.Writer, export *ExportTransaction) error {
	if err := printCommonInformation(
		out,
		export,
	); err != nil {
		return err
	}

	fctl.Section.WithWriter(out).Println("Summary")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{"Account", "Asset", "Movement", "Final balance"})
	for _, summary := range export.Summary {
		movementStr := fmt.Sprint(summary.Movement)
		if summary.Movement.Cmp(big.NewInt(0)) > 0 {
			movementStr = "+" + movementStr
		}

		tableData = append(tableData, []string{
			summary.Account, summary.Asset, movementStr, summary.FinalBalance.String(),
		})
	}

	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fmt.Fprintln(out, "")
	if err := PrintMetadata(out, export.Metadata); err != nil {
		return err
	}
	return nil
}

func PrintTransaction(out io.Writer, export *ExportTransaction) error {

	if err := printCommonInformation(
		out,
		export,
	); err != nil {
		return err
	}

	if err := PrintMetadata(out, export.Metadata); err != nil {
		return err
	}
	return nil
}
func PrintMetadata(out io.Writer, metadata Metadata) error {
	fctl.Section.WithWriter(out).Println("Metadata")
	if len(metadata) == 0 {
		fmt.Fprintln(out, "No metadata.")
		return nil
	}
	tableData := pterm.TableData{}
	for k, v := range metadata {
		asJson, err := json.Marshal(v)
		if err != nil {
			return err
		}
		tableData = append(tableData, []string{pterm.LightCyan(k), string(asJson)})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}
