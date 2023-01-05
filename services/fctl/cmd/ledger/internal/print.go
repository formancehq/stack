package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pterm/pterm"
)

func PrintMetadata(out io.Writer, metadata map[string]any) error {
	if len(metadata) == 0 {
		fctl.Highlightln(out, "Metadata : <empty>")
		return nil
	}
	fctl.Highlightln(out, "Metadata :")
	tableData := pterm.TableData{}
	for k, v := range metadata {
		data, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		tableData = append(tableData, []string{pterm.LightCyan(k), string(data)})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}

func PrintTransaction(out io.Writer, transaction formance.Transaction) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(transaction.Txid)})
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), fctl.StringPointerToString(transaction.Reference)})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), transaction.Timestamp.Format(time.RFC3339)})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	fmt.Fprintln(out, "")

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

	tableData = pterm.TableData{}
	tableData = append(tableData, []string{"Account", "Asset", "Movement", "Final balance"})
	for account, postCommitVolume := range *transaction.PostCommitVolumes {
		for asset, volumes := range postCommitVolume {
			movement := *volumes.Balance - *(*transaction.PreCommitVolumes)[account][asset].Balance
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

	if err := PrintMetadata(out, transaction.Metadata); err != nil {
		return err
	}
	return nil
}
