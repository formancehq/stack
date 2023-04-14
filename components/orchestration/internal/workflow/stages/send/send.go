package send

import (
	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/schema"
	"github.com/formancehq/orchestration/internal/workflow/stages"
)

type WalletSource struct {
	ID      string `json:"id" validate:"required"`
	Balance string `json:"balance" spec:"default:main" validate:"required"`
}

type WalletDestination = WalletSource

type LedgerAccountSource struct {
	ID     string `json:"id" validate:"required"`
	Ledger string `json:"ledger" spec:"default:default" validate:"required"`
}

type LedgerAccountDestination = LedgerAccountSource

type PaymentSource struct {
	ID string `json:"id" validate:"required"`
}

type PaymentDestination struct {
	PSP      string `json:"psp"`
	Metadata string `json:"metadata" spec:"default:stripeConnectID"`
}

type Source struct {
	Wallet  *WalletSource        `json:"wallet,omitempty"`
	Account *LedgerAccountSource `json:"account,omitempty"`
	Payment *PaymentSource       `json:"payment,omitempty"`
}

func NewSource() *Source {
	return &Source{}
}

func (s Source) WithWallet(src *WalletSource) Source {
	s.Wallet = src
	return s
}

func (s Source) WithPayment(src *PaymentSource) Source {
	s.Payment = src
	return s
}

func (s Source) WithAccount(src *LedgerAccountSource) Source {
	s.Account = src
	return s
}

type Destination struct {
	Wallet  *WalletDestination        `json:"wallet,omitempty"`
	Account *LedgerAccountDestination `json:"account,omitempty"`
	Payment *PaymentDestination       `json:"payment,omitempty"`
}

func NewDestination() *Destination {
	return &Destination{}
}

func (s Destination) WithWallet(src *WalletDestination) Destination {
	s.Wallet = src
	return s
}

func (s Destination) WithPayment(src *PaymentDestination) Destination {
	s.Payment = src
	return s
}

func (s Destination) WithAccount(src *LedgerAccountDestination) Destination {
	s.Account = src
	return s
}

type Send struct {
	Source      Source        `json:"source"`
	Destination Destination   `json:"destination"`
	Amount      *sdk.Monetary `json:"amount,omitempty"`
}

func (s Send) GetWorkflow() any {
	return RunSend
}

func init() {
	schema.RegisterOneOf(&Source{}, &Destination{})
	stages.Register("send", Send{})
}
