package vm

import "github.com/formancehq/ledger/pkg/machine/internal"

type MockupLedger map[string]MockupAccount

func NewMockupLedger() MockupLedger {
	return MockupLedger(make(map[string]MockupAccount))
}

func (l MockupLedger) GetBalance(account internal.AccountAddress, asset internal.Asset) *internal.MonetaryInt {
	balance := l[string(account)].Balances[string(asset)]
	return &balance
}

func (l MockupLedger) GetMeta(account internal.AccountAddress, key string) internal.Value {
	return l[string(account)].Meta[key]
}

type MockupAccount struct {
	Balances map[string]internal.MonetaryInt
	Meta     map[string]internal.Value
}

func NewMockupAccount() MockupAccount {
	return MockupAccount{
		Balances: make(map[string]internal.MonetaryInt),
		Meta:     make(map[string]internal.Value),
	}
}
