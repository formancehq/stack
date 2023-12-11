package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func printCommonInformation(
	out io.Writer,
	txID *big.Int,
	reference string,
	postings []shared.V2Posting,
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

func PrintExpandedTransaction(out io.Writer, transaction shared.V2ExpandedTransaction) error {

	if err := printCommonInformation(
		out,
		transaction.ID,
		func() string {
			if transaction.Reference == nil {
				return ""
			}
			return *transaction.Reference
		}(),
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
			movement := big.NewInt(0)
			movement = movement.Sub(volumes.Balance, (transaction.PreCommitVolumes)[account][asset].Balance)
			movementStr := fmt.Sprint(movement)
			if movement.Cmp(big.NewInt(0)) > 0 {
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

func PrintTransaction(out io.Writer, transaction shared.V2Transaction) error {

	reference := ""
	if transaction.Reference != nil {
		reference = *transaction.Reference
	}
	if err := printCommonInformation(
		out,
		transaction.ID,
		reference,
		transaction.Postings,
		transaction.Timestamp,
	); err != nil {
		return err
	}

	if err := PrintMetadata(out, transaction.Metadata); err != nil {
		return err
	}
	return nil
}
func PrintMetadata(out io.Writer, metadata map[string]string) error {
	fctl.Section.WithWriter(out).Println("Metadata")
	if len(metadata) == 0 {
		fmt.Println("No metadata.")
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
