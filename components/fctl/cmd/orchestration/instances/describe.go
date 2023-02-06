package instances

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDescribeCommand() *cobra.Command {
	return fctl.NewCommand("describe <instance-id>",
		fctl.WithShortDescription("Describe a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return errors.Wrap(err, "retrieving config")
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			ret, _, err := client.OrchestrationApi.GetInstanceHistory(cmd.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			defaultWriter := fctl.BasicText.WithWriter(cmd.OutOrStdout())
			greenWriter := fctl.BasicTextGreen.WithWriter(cmd.OutOrStdout())
			redWriter := fctl.BasicTextRed.WithWriter(cmd.OutOrStdout())
			cyanWriter := fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout())

			for i, history := range ret.Data {
				switch {
				case history.Input.StageSend != nil:
					printHistoryBaseInfo(cmd.OutOrStdout(), "send", i, history)
					cyanWriter.Printfln("Send %d %s from %s to %s", history.Input.StageSend.Amount.Amount,
						history.Input.StageSend.Amount.Asset, stageSourceName(history.Input.StageSend.Source),
						stageDestinationName(history.Input.StageSend.Destination))
					fctl.Println()
					defaultWriter.Println("Activities :")

					stageResponse, _, err := client.OrchestrationApi.GetInstanceStageHistory(cmd.Context(), args[0], int32(i)).Execute()
					if err != nil {
						return err
					}

					for _, historyStage := range stageResponse.Data {
						switch {
						case historyStage.Input.StripeTransfer != nil:
							greenWriter.Printfln("Send %d %s to Stripe connected account: %s",
								*historyStage.Input.StripeTransfer.Amount,
								*historyStage.Input.StripeTransfer.Asset,
								*historyStage.Input.StripeTransfer.Destination,
							)
						case historyStage.Input.CreateTransaction != nil:
							greenWriter.Printfln("Send %d %s from account %s to account %s",
								historyStage.Input.CreateTransaction.Data.Postings[0].Amount,
								historyStage.Input.CreateTransaction.Data.Postings[0].Asset,
								historyStage.Input.CreateTransaction.Data.Postings[0].Source,
								historyStage.Input.CreateTransaction.Data.Postings[0].Destination,
							)
							if historyStage.Error == nil {
								//defaultWriter.Printfln("\tCreated transaction: %d", historyStage.Output.CreateTransaction.Data[0].Txid)
								if historyStage.Input.CreateTransaction.Data.Reference != nil {
									defaultWriter.Printfln("\tReference: %s", *historyStage.Output.CreateTransaction.Data[0].Reference)
								}
								if len(historyStage.Input.CreateTransaction.Data.Metadata) > 0 {
									printMetadata(defaultWriter, historyStage.Input.CreateTransaction.Data.Metadata)
								}
							}
						case historyStage.Input.ConfirmHold != nil:
							greenWriter.Printfln("Confirm debit hold %s", historyStage.Input.ConfirmHold.Id)
						case historyStage.Input.CreditWallet != nil:
							greenWriter.Printfln("Credit wallet %s (balance: %s) of %d %s from %s",
								*historyStage.Input.CreditWallet.Id,
								*historyStage.Input.CreditWallet.Data.Balance,
								historyStage.Input.CreditWallet.Data.Amount.Amount,
								historyStage.Input.CreditWallet.Data.Amount.Asset,
								subjectName(historyStage.Input.CreditWallet.Data.Sources[0]),
							)
							if historyStage.Error == nil {
								if len(historyStage.Input.CreditWallet.Data.Metadata) > 0 {
									printMetadata(defaultWriter, historyStage.Input.CreditWallet.Data.Metadata)
								}
							}
						case historyStage.Input.DebitWallet != nil:
							destination := "@world"
							if historyStage.Input.DebitWallet.Data.Destination != nil {
								destination = subjectName(*historyStage.Input.DebitWallet.Data.Destination)
							}
							greenWriter.Printfln("Debit wallet %s (balance: %s) of %d %s to %s",
								*historyStage.Input.DebitWallet.Id,
								historyStage.Input.DebitWallet.Data.Balances[0],
								historyStage.Input.DebitWallet.Data.Amount.Amount,
								historyStage.Input.DebitWallet.Data.Amount.Asset,
								destination,
							)
							if historyStage.Error == nil {
								if len(historyStage.Input.DebitWallet.Data.Metadata) > 0 {
									printMetadata(defaultWriter, historyStage.Input.DebitWallet.Data.Metadata)
								}
							}
						case historyStage.Input.GetAccount != nil:
							greenWriter.Printfln("Read account %s of ledger %s",
								historyStage.Input.GetAccount.Id,
								historyStage.Input.GetAccount.Ledger,
							)
						case historyStage.Input.GetPayment != nil:
							greenWriter.Printfln("Read payment %s", historyStage.Input.GetPayment.Id)
						case historyStage.Input.GetWallet != nil:
							greenWriter.Printfln("Read wallet '%s'", historyStage.Input.GetWallet.Id)
						case historyStage.Input.RevertTransaction != nil:
							greenWriter.Printfln("Revert transaction %s", historyStage.Input.RevertTransaction.Id)
							if historyStage.Error == nil {
								defaultWriter.Printfln("\tCreated transaction: %d", historyStage.Output.RevertTransaction.Data.Txid)
							}
						case historyStage.Input.VoidHold != nil:
							greenWriter.Printfln("Cancel debit hold %s", historyStage.Input.VoidHold.Id)
						}
						if historyStage.Error != nil {
							redWriter.WithWriter(cmd.OutOrStdout()).Printfln("\t%s", *historyStage.Error)
						}
					}
				case history.Input.StageDelay != nil:
					printHistoryBaseInfo(cmd.OutOrStdout(), "delay", i, history)
					switch {
					case history.Input.StageDelay.Duration != nil:
						cyanWriter.Printfln("Pause workflow for a delay of %s", *history.Input.StageDelay.Duration)
					case history.Input.StageDelay.Until != nil:
						cyanWriter.Printfln("Pause workflow until %s", *history.Input.StageDelay.Until)
					}
				case history.Input.StageWaitEvent != nil:
					printHistoryBaseInfo(cmd.OutOrStdout(), "wait_event", i, history)
					cyanWriter.Printfln("Waiting event '%s'", history.Input.StageWaitEvent.Event)
					if history.Error == nil {
						if history.Terminated {
							defaultWriter.Printfln("\tEvent received!")
						} else {
							defaultWriter.Printfln("\tStill waiting event...")
						}
					}
				default:
					// Display error?
				}
				if history.Error != nil {
					redWriter.WithWriter(cmd.OutOrStdout()).Printfln("Stage terminated with error: %s", *history.Error)
				}
			}

			return nil
		}),
	)
}

func printHistoryBaseInfo(out io.Writer, name string, ind int, history formance.WorkflowInstanceHistory) {
	fctl.Section.WithWriter(out).Printf("Stage %d : %s\n", ind, name)
	fctl.BasicText.WithWriter(out).Printfln("Started at: %s", history.StartedAt.Format(time.RFC3339))
	if history.Terminated {
		fctl.BasicText.WithWriter(out).Printfln("Terminated at: %s", history.StartedAt.Format(time.RFC3339))
	}
}

func stageSourceName(src *formance.StageSendSource) string {
	switch {
	case src.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", src.Wallet.Id, *src.Wallet.Balance)
	case src.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", src.Account.Id, *src.Account.Ledger)
	case src.Payment != nil:
		return fmt.Sprintf("payment '%s'", src.Payment.Id)
	default:
		return "unknown_source_type"
	}
}

func stageDestinationName(dst *formance.StageSendDestination) string {
	switch {
	case dst.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", dst.Wallet.Id, *dst.Wallet.Balance)
	case dst.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", dst.Account.Id, *dst.Account.Ledger)
	case dst.Payment != nil:
		return dst.Payment.Psp
	default:
		return "unknown_source_type"
	}
}

func subjectName(src formance.Subject) string {
	switch {
	case src.WalletSubject != nil:
		return fmt.Sprintf("wallet %s (balance: %s)", src.WalletSubject.Identifier, *src.WalletSubject.Balance)
	case src.LedgerAccountSubject != nil:
		return fmt.Sprintf("account %s", src.LedgerAccountSubject.Identifier)
	default:
		return "unknown_subject_type"
	}
}

func printMetadata(writer *pterm.BasicTextPrinter, metadata map[string]any) {
	writer.Printfln("\tAdded metadata:")
	for k, v := range metadata {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		writer.Printfln("\t- %s: %s", k, string(jsonValue))
	}
}
