package instances

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
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

			for i, history := range ret.Data {
				switch {
				case history.Input.StageSend != nil:
					fctl.Printf("Send %d %s from %s to %s", history.Input.StageSend.Amount.Amount,
						history.Input.StageSend.Amount.Asset, stageSourceName(history.Input.StageSend.Source),
						stageDestinationName(history.Input.StageSend.Destination))
					terminated := history.Error != nil
					fctl.Printf("---------------------------------------------")

					stageResponse, _, err := client.OrchestrationApi.GetInstanceStageHistory(cmd.Context(), args[0], int32(i)).Execute()
					if err != nil {
						return err
					}

					for _, historyStage := range stageResponse.Data {
						switch {
						case historyStage.Input.StripeTransferRequest != nil:
							fctl.Printf("Send %d %s (from @world) to Stripe connected account: %s",
								*historyStage.Input.StripeTransferRequest.Amount,
								*historyStage.Input.StripeTransferRequest.Asset,
								*historyStage.Input.StripeTransferRequest.Destination,
							)
						case historyStage.Input.ActivityCreateTransaction != nil:
							fctl.Printf("Send %d %s from account %s to account %s",
								historyStage.Input.ActivityCreateTransaction.Data.Postings[0].Amount,
								historyStage.Input.ActivityCreateTransaction.Data.Postings[0].Asset,
								historyStage.Input.ActivityCreateTransaction.Data.Postings[0].Source,
								historyStage.Input.ActivityCreateTransaction.Data.Postings[0].Destination,
							)
						case historyStage.Input.ActivityConfirmHold != nil:
							fctl.Printf("Confirm debit hold %s", historyStage.Input.ActivityConfirmHold.Id)
						case historyStage.Input.ActivityCreditWallet != nil:
							fctl.Printf("Credit wallet %s (balance: %s) of %d %s from %s",
								*historyStage.Input.ActivityCreditWallet.Id,
								*historyStage.Input.ActivityCreditWallet.Data.Balance,
								historyStage.Input.ActivityCreditWallet.Data.Amount.Amount,
								historyStage.Input.ActivityCreditWallet.Data.Amount.Asset,
								subjectName(historyStage.Input.ActivityCreditWallet.Data.Sources[0]),
							)
						case historyStage.Input.ActivityDebitWallet != nil:
							fctl.Printf("Debit wallet %s (balance: %s) of %d %s to %s",
								*historyStage.Input.ActivityDebitWallet.Id,
								historyStage.Input.ActivityDebitWallet.Data.Balances[0],
								historyStage.Input.ActivityDebitWallet.Data.Amount.Amount,
								historyStage.Input.ActivityDebitWallet.Data.Amount.Asset,
								subjectName(*historyStage.Input.ActivityDebitWallet.Data.Destination),
							)
						case historyStage.Input.ActivityGetAccount != nil:
							fctl.Printf("Read account %s of ledger %s",
								historyStage.Input.ActivityGetAccount.Id,
								historyStage.Input.ActivityGetAccount.Ledger,
							)
						case historyStage.Input.ActivityGetPayment != nil:
							fctl.Printf("Read payment %s", historyStage.Input.ActivityGetPayment.Id)
						case historyStage.Input.ActivityGetWallet != nil:
							fctl.Printf("Read wallet %s", historyStage.Input.ActivityGetWallet.Id)
						case historyStage.Input.ActivityRevertTransaction != nil:
							fctl.Printf("Revert transaction %s", historyStage.Input.ActivityRevertTransaction.Id)
						case historyStage.Input.ActivityVoidHold != nil:
							fctl.Printf("Cancel debit hold %s", historyStage.Input.ActivityVoidHold.Id)
						}
					}

					if terminated {
						return nil
					}
				default:
					// Display error?
				}
			}

			return nil
		}),
	)
}

func stageSourceName(src any) string {
	switch src := src.(type) {
	case *formance.WalletSource:
		return fmt.Sprintf("wallet %s (balance: %s)", src.Id, src.Balance)
	case *formance.LedgerAccountSource:
		return fmt.Sprintf("account %s (ledger: %s)", src.Id, src.Ledger)
	case *formance.PaymentSource:
		return fmt.Sprintf("payment %s", src.Id)
	default:
		return "unknown_source_type"
	}
}

func stageDestinationName(dst any) string {
	switch src := dst.(type) {
	case *formance.WalletSource:
		return fmt.Sprintf("wallet %s (balance: %s)", src.Id, src.Balance)
	case *formance.LedgerAccountSource:
		return fmt.Sprintf("account %s (ledger: %s)", src.Id, src.Ledger)
	case *formance.PaymentDestination:
		return src.Psp
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
