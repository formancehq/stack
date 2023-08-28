package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintHistoryBaseInfo(out io.Writer, name string, ind int, history shared.WorkflowInstanceHistory) {
	fctl.Section.WithWriter(out).Printf("Stage %d : %s\n", ind, name)
	fctl.BasicText.WithWriter(out).Printfln("Started at: %s", history.StartedAt.Format(time.RFC3339))
	if history.Terminated {
		fctl.BasicText.WithWriter(out).Printfln("Terminated at: %s", history.StartedAt.Format(time.RFC3339))
	}
}

func StageSourceName(src *shared.StageSendSource) string {
	switch {
	case src.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", src.Wallet.ID, *src.Wallet.Balance)
	case src.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", src.Account.ID, *src.Account.Ledger)
	case src.Payment != nil:
		return fmt.Sprintf("payment '%s'", src.Payment.ID)
	default:
		return "unknown_source_type"
	}
}

func StageDestinationName(dst *shared.StageSendDestination) string {
	switch {
	case dst.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", dst.Wallet.ID, *dst.Wallet.Balance)
	case dst.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", dst.Account.ID, *dst.Account.Ledger)
	case dst.Payment != nil:
		return dst.Payment.Psp
	default:
		return "unknown_source_type"
	}
}

func SubjectName(src shared.Subject) string {
	switch {
	case src.WalletSubject != nil:
		return fmt.Sprintf("wallet %s (balance: %s)", src.WalletSubject.Identifier, *src.WalletSubject.Balance)
	case src.LedgerAccountSubject != nil:
		return fmt.Sprintf("account %s", src.LedgerAccountSubject.Identifier)
	default:
		return "unknown_subject_type"
	}
}

func PrintMetadata(metadata map[string]string) []pterm.BulletListItem {
	ret := make([]pterm.BulletListItem, 0)
	ret = append(ret, HistoryItemDetails("Added metadata:"))
	for k, v := range metadata {
		ret = append(ret, pterm.BulletListItem{
			Level: 2,
			Text:  fmt.Sprintf("%s: %s", k, v),
		})
	}
	return ret
}

func PrintStage(out io.Writer, ctx context.Context, i int, client *formance.Formance, id string, history shared.WorkflowInstanceHistory) error {
	cyanWriter := fctl.BasicTextCyan
	defaultWriter := fctl.BasicText

	listItems := make([]pterm.BulletListItem, 0)

	buf, err := json.Marshal(history.Input)
	if err != nil {
		return err
	}

	var (
		stageSend      shared.StageSend
		stageDelay     shared.StageDelay
		stageWaitEvent shared.StageWaitEvent
	)

	err = json.Unmarshal(buf, &stageSend)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &stageDelay)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &stageWaitEvent)
	if err != nil {
		return err
	}

	switch {
	case stageSend.Amount != nil && stageSend.Source != nil && stageSend.Destination != nil:
		PrintHistoryBaseInfo(out, "send", i, history)
		cyanWriter.Printfln("Send %d %s from %s to %s", stageSend.Amount.Amount,
			stageSend.Amount.Asset, StageSourceName(stageSend.Source),
			StageDestinationName(stageSend.Destination))
		fmt.Fprintln(out)

		stageResponse, err := client.Orchestration.GetInstanceStageHistory(ctx, operations.GetInstanceStageHistoryRequest{
			InstanceID: id,
			Number:     int64(i),
		})
		if err != nil {
			return err
		}

		for _, historyStage := range stageResponse.GetWorkflowInstanceHistoryStageResponse.Data {
			switch {
			case historyStage.Input.StripeTransfer != nil:
				listItems = append(listItems, HistoryItemTitle("Send %d %s to Stripe connected account: %s",
					*historyStage.Input.StripeTransfer.Amount,
					*historyStage.Input.StripeTransfer.Asset,
					*historyStage.Input.StripeTransfer.Destination,
				))
			case historyStage.Input.CreateTransaction != nil:
				listItems = append(listItems, HistoryItemTitle("Send %d %s from account %s to account %s (ledger %s)",
					historyStage.Input.CreateTransaction.Data.Postings[0].Amount,
					historyStage.Input.CreateTransaction.Data.Postings[0].Asset,
					historyStage.Input.CreateTransaction.Data.Postings[0].Source,
					historyStage.Input.CreateTransaction.Data.Postings[0].Destination,
					*historyStage.Input.CreateTransaction.Ledger,
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					listItems = append(listItems, HistoryItemDetails("Created transaction: %d", historyStage.Output.CreateTransaction.Data.Txid))
					if historyStage.Input.CreateTransaction.Data.Reference != nil {
						listItems = append(listItems, HistoryItemDetails("Reference: %s", *historyStage.Output.CreateTransaction.Data.Reference))
					}
					if len(historyStage.Input.CreateTransaction.Data.Metadata) > 0 {
						listItems = append(listItems, PrintMetadata(historyStage.Input.CreateTransaction.Data.Metadata)...)
					}
				}
			case historyStage.Input.ConfirmHold != nil:
				listItems = append(listItems, HistoryItemTitle("Confirm debit hold %s", historyStage.Input.ConfirmHold.ID))
			case historyStage.Input.CreditWallet != nil:
				listItems = append(listItems, HistoryItemTitle("Credit wallet %s (balance: %s) of %d %s from %s",
					*historyStage.Input.CreditWallet.ID,
					*historyStage.Input.CreditWallet.Data.Balance,
					historyStage.Input.CreditWallet.Data.Amount.Amount,
					historyStage.Input.CreditWallet.Data.Amount.Asset,
					SubjectName(historyStage.Input.CreditWallet.Data.Sources[0]),
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					if len(historyStage.Input.CreditWallet.Data.Metadata) > 0 {
						listItems = append(listItems, PrintMetadata(historyStage.Input.CreditWallet.Data.Metadata)...)
					}
				}
			case historyStage.Input.DebitWallet != nil:
				destination := "@world"
				if historyStage.Input.DebitWallet.Data.Destination != nil {
					destination = SubjectName(*historyStage.Input.DebitWallet.Data.Destination)
				}

				listItems = append(listItems, HistoryItemTitle("Debit wallet %s (balance: %s) of %d %s to %s",
					*historyStage.Input.DebitWallet.ID,
					historyStage.Input.DebitWallet.Data.Balances[0],
					historyStage.Input.DebitWallet.Data.Amount.Amount,
					historyStage.Input.DebitWallet.Data.Amount.Asset,
					destination,
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					if len(historyStage.Input.DebitWallet.Data.Metadata) > 0 {
						listItems = append(listItems, PrintMetadata(historyStage.Input.DebitWallet.Data.Metadata)...)
					}
				}
			case historyStage.Input.GetAccount != nil:
				listItems = append(listItems, HistoryItemTitle("Read account %s of ledger %s",
					historyStage.Input.GetAccount.ID,
					historyStage.Input.GetAccount.Ledger,
				))
			case historyStage.Input.GetPayment != nil:
				listItems = append(listItems, HistoryItemTitle("Read payment %s",
					historyStage.Input.GetPayment.ID))
			case historyStage.Input.GetWallet != nil:
				listItems = append(listItems, HistoryItemTitle("Read wallet '%s'", historyStage.Input.GetWallet.ID))
			case historyStage.Input.RevertTransaction != nil:
				listItems = append(listItems, HistoryItemTitle("Revert transaction %s", historyStage.Input.RevertTransaction.ID))
				if historyStage.Error == nil {
					listItems = append(listItems, HistoryItemTitle("Created transaction: %d", historyStage.Output.RevertTransaction.Data.Txid))
				}
			case historyStage.Input.VoidHold != nil:
				listItems = append(listItems, HistoryItemTitle("Cancel debit hold %s", historyStage.Input.VoidHold.ID))
			}
			if historyStage.LastFailure != nil {
				listItems = append(listItems, HistoryItemError(*historyStage.LastFailure))
				if historyStage.NextExecution != nil {
					listItems = append(listItems, HistoryItemError("Next try: %s", historyStage.NextExecution.Format(time.RFC3339)))
					listItems = append(listItems, HistoryItemError("Attempt: %d", historyStage.Attempt))
				}
			}
			if historyStage.Error != nil {
				listItems = append(listItems, HistoryItemError(*historyStage.Error))
			}
		}
	case stageDelay.Duration != nil && stageDelay.Until != nil:
		PrintHistoryBaseInfo(out, "delay", i, history)
		switch {
		case stageDelay.Duration != nil:
			listItems = append(listItems, HistoryItemTitle("Pause workflow for a delay of %s", *stageDelay.Duration))
		case stageDelay.Until != nil:
			listItems = append(listItems, HistoryItemTitle("Pause workflow until %s", *stageDelay.Until))
		}
	case stageWaitEvent.Event != "":
		PrintHistoryBaseInfo(out, "wait_event", i, history)
		listItems = append(listItems, HistoryItemTitle("Waiting event '%s'", stageWaitEvent.Event))
		if history.Error == nil {
			if history.Terminated {
				listItems = append(listItems, HistoryItemDetails("Event received!"))
			} else {
				listItems = append(listItems, HistoryItemDetails("Still waiting event..."))
			}
		}
	default:
		// Display error?
	}
	if history.Error != nil {
		fctl.BasicTextRed.WithWriter(out).Printfln("Stage terminated with error: %s", *history.Error)
	}

	if len(listItems) > 0 {
		defaultWriter.Print("History :\n")
		return pterm.DefaultBulletList.WithWriter(out).WithItems(listItems).Render()
	}
	return nil
}

func HistoryItemTitle(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level:     0,
		TextStyle: fctl.StyleGreen,
		Text:      fmt.Sprintf(format, args...),
	}
}

func HistoryItemDetails(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level: 1,
		Text:  fmt.Sprintf(format, args...),
	}
}

func HistoryItemError(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level:     1,
		TextStyle: fctl.StyleRed,
		Text:      fmt.Sprintf(format, args...),
	}
}
