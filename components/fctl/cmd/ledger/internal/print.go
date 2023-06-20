package internal

import (
	"fmt"
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func printCommonInformation(
	out io.Writer,
	txID int64,
	reference string,
	postings []shared.Posting,
	timestamp time.Time,
) error {
	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(txID)})
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), reference})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), timestamp.Format(time.RFC3339)})

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
	for _, posting := range postings {
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

func PrintExpandedTransaction(out io.Writer, transaction shared.ExpandedTransaction) error {

	if err := printCommonInformation(
		out,
		transaction.Txid,
		*transaction.Reference,
		transaction.Postings,
		transaction.Timestamp,
	); err != nil {
		return err
	}

	fctl.Section.WithWriter(out).Println("Summary")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{"Account", "Asset", "Movement", "Final balance"})
	for account, postCommitVolume := range transaction.PostCommitVolumes {
		for asset, volumes := range postCommitVolume {
			movement := *volumes.Balance - *(transaction.PreCommitVolumes)[account][asset].Balance
			movementStr := fmt.Sprint(movement)
			if movement > 0 {
				movementStr = "+" + movementStr
			}
			tableData = append(tableData, []string{
				account, asset, movementStr, fmt.Sprint(*volumes.Balance),
			})
		}
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fmt.Fprintln(out, "")
	if err := fctl.PrintMetadata(out, transaction.Metadata); err != nil {
		return err
	}
	return nil
}

func PrintTransaction(out io.Writer, transaction shared.Transaction) error {

	if err := printCommonInformation(
		out,
		transaction.Txid,
		*transaction.Reference,
		transaction.Postings,
		transaction.Timestamp,
	); err != nil {
		return err
	}

	if err := fctl.PrintMetadata(out, transaction.Metadata); err != nil {
		return err
	}
	return nil
}
