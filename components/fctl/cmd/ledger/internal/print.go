package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
)

func printCommonInformation(
	out io.Writer,
	txID *big.Int,
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
			posting.GetSource(), posting.GetDestination(), posting.GetAsset(), posting.GetAmount().String(),
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

func PrintExpandedTransaction(out io.Writer, transaction ExpandedTransaction) error {

	if err := printCommonInformation(
		out,
		transaction.GetID(),
		func() string {
			if transaction.GetReference() == nil {
				return ""
			}
			return *transaction.GetReference()
		}(),
		transaction.GetPostings(),
		transaction.GetTimestamp(),
	); err != nil {
		return err
	}

	fctl.Section.WithWriter(out).Println("Summary")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{"Account", "Asset", "Movement", "Final balance"})
	for account, postCommitVolume := range transaction.GetPostCommitVolumes() {
		for asset, volumes := range postCommitVolume {
			movement := big.NewInt(0)
			movement = movement.Sub(volumes.Balance, (transaction.GetPreCommitVolumes())[account][asset].Balance)
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
	if err := PrintMetadata(out, transaction.GetMetadata()); err != nil {
		return err
	}
	return nil
}

type Transaction interface {
	GetReference() *string
	GetID() *big.Int
	GetPostings() []shared.Posting
	GetTimestamp() time.Time
	GetMetadata() map[string]string
}

type v2Transaction struct {
	shared.V2Transaction
}

func (t v2Transaction) GetPostings() []shared.Posting {
	return collectionutils.Map(t.V2Transaction.GetPostings(), func(from shared.V2Posting) shared.Posting {
		return shared.Posting{
			Amount:      from.GetAmount(),
			Asset:       from.GetAsset(),
			Destination: from.GetDestination(),
			Source:      from.GetSource(),
		}
	})
}

var _ Transaction = (*v2Transaction)(nil)

func WrapV2Transaction(transaction shared.V2Transaction) *v2Transaction {
	return &v2Transaction{
		V2Transaction: transaction,
	}
}

type v1Transaction struct {
	shared.Transaction
}

func (t v1Transaction) GetID() *big.Int {
	return t.Transaction.GetTxid()
}

func (t v1Transaction) GetMetadata() map[string]string {
	return collectionutils.ConvertMap(t.Transaction.Metadata, collectionutils.ToFmtString)
}

var _ Transaction = (*v1Transaction)(nil)

func WrapV1Transaction(transaction shared.Transaction) *v1Transaction {
	return &v1Transaction{
		Transaction: transaction,
	}
}

type ExpandedTransaction interface {
	Transaction
	GetPreCommitVolumes() map[string]map[string]shared.Volume
	GetPostCommitVolumes() map[string]map[string]shared.Volume
}

func PrintTransaction(out io.Writer, transaction Transaction) error {

	reference := ""
	if transaction.GetReference() != nil {
		reference = *transaction.GetReference()
	}
	if err := printCommonInformation(
		out,
		transaction.GetID(),
		reference,
		transaction.GetPostings(),
		transaction.GetTimestamp(),
	); err != nil {
		return err
	}

	// collectionutils.ConvertMap(transaction.Metadata, collectionutils.ToFmtString[string])
	if err := PrintMetadata(out, transaction.GetMetadata()); err != nil {
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
