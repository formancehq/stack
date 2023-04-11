package vm

import (
	"fmt"
	"github.com/numary/ledger/pkg/machine/vm/program"
	"github.com/numary/ledger/pkg/core"
)

type LedgerReader interface {
	GetBalance(account core.AccountAddress, asset core.Asset) core.Number
	GetMeta(account core.AccountAddress, key string) core.Value
}


type Interpreter struct {
	ledger LedgerReader
	vars map[string]core.Value
	balances map[core.AccountAddress]map[core.Asset]core.Number
	tx_meta map[string]core.Value
	AccountMeta map[core.AccountAddress]map[string]core.Value
	Postings []core.Posting
}

func (m Interpreter) Run(script program.Script) {
	for _, var_decl := range script.VarsDecl {
		switch o := (*var_decl.Origin).(type) {
		case program.VarOriginMeta:
			account := EvalAs[core.AccountAddress](m, o.Account)
			value := m.ledger.GetMeta(account, o.Key)
			m.vars[var_decl.Name] = value
		case program.VarOriginBalance:
			account := EvalAs[core.AccountAddress](m, o.Account)
			asset := EvalAs[core.Asset](m, o.Asset)
			balance := m.ledger.GetBalance(account, asset)
			m.vars[var_decl.Name] = balance
		default:
			panic("ice")
		}
	}

	for _, stmt := range script.Statements {
		switch s := stmt.(type) {
		case program.StatementFail:
			panic("fail")
		case program.StatementPrint:
			fmt.Printf("%v\n", s.Expr)
		case program.StatementAllocate:
			funding := EvalAs[core.Funding](m, s.Funding)
			m.Allocate(funding, s.Destination)
		case program.StatementLet:
			value := m.Eval(s.Expr)
			m.vars[s.Name] = value
		case program.StatementSetTxMeta:
			value := m.Eval(s.Value)
			m.tx_meta[s.Key] = value
		case program.StatementSetAccountMeta:
			account := EvalAs[core.AccountAddress](m, s.Account)
			value := m.Eval(s.Value)
			if _, ok := m.AccountMeta[account]; !ok {
				m.AccountMeta[account] = make(map[string]core.Value)
			}
			m.AccountMeta[account][s.Key] = value
		default:
			panic("ice")
		}
	}
}

func (m Interpreter) Send(funding core.Funding, account core.AccountAddress) {
	for _, part := range funding.Parts {
		m.Postings = append(m.Postings, core.Posting {
			Source: string(part.Account),
			Destination: string(account),
			Asset: string(funding.Asset),
			Amount: part.Amount,
		})
	}
}

func (m Interpreter) Allocate(funding core.Funding, destination program.Destination) {
	switch d := destination.(type) {
	case program.DestinationAccount:
		account := EvalAs[core.AccountAddress](m, d.Expr)
		m.Send(funding, account)
	case program.DestinationInOrder:
		for _, part := range d.Parts {
			max := EvalAs[core.Monetary](m, part.Max)
			taken, remainder := funding.TakeMax(max.Amount)
			funding = remainder
			m.AllocateOrKeep(taken, part.Kod)
		}
		m.AllocateOrKeep(funding, d.Remaining)
	case program.DestinationAllotment:
		portions := make([]core.Portion, 0)
		sub_dests := make([]program.KeptOrDestination, 0)
		for _, part := range d {
			if part.Portion.Remaining {
				portions = append(portions, core.NewPortionRemaining())
			} else {
				portion := EvalAs[core.Portion](m, *part.Portion.Expr)
				portions = append(portions, portion)
			}
			sub_dests = append(sub_dests, part.Kod)
		}
		allotment, err := core.NewAllotment(portions)
		if err != nil {
			panic("invalid allotment")
		}
		for i, part := range allotment.Allocate(funding.Total()) {
			taken, remainder, err := funding.Take(part)
			if err != nil {
				panic("insufficient funds")
			}
			funding = remainder
			m.AllocateOrKeep(taken, sub_dests[i])
		}
	}
}

func (m Interpreter) AllocateOrKeep(funding core.Funding, kod program.KeptOrDestination) {
	if kod.Kept {
		m.Repay(funding)
	} else {
		m.Allocate(funding, kod.Destination)
	}
}

func (m Interpreter) TakeFromValueAwareSource(source program.ValueAwareSource, mon core.Monetary) core.Funding {
	var err error
	switch s := source.(type) {
	case program.ValueAwareSourceSource:
		available, fallback := m.TakeFromSource(s.Source, mon.Asset)
		taken, remainder := available.TakeMax(mon.Amount)
		if taken.Total() != mon.Amount {
			missing := core.Monetary{
				Asset: mon.Asset,
				Amount: mon.Amount.Sub(taken.Total()),
			}
			if fallback != nil {
				missing_taken := m.WithdrawAlways(*fallback, missing)
				taken, err = taken.Concat(missing_taken)
				if err != nil {
					panic("??")
				}
			} else {
				panic("insufficient funds")
			}
		}
		m.Repay(remainder)
		return taken
	case program.ValueAwareSourceAllotment:
		portions := make([]core.Portion, 0)
		sub_sources := make([]program.Source, 0)
		for _, part := range s {
			if part.Portion.Remaining {
				portions = append(portions, core.NewPortionRemaining())
			} else {
				portion := EvalAs[core.Portion](m, *part.Portion.Expr)
				portions = append(portions, portion)
			}
			sub_sources = append(sub_sources, part.Source)
		}
		allotment, err := core.NewAllotment(portions)
		funding := core.Funding {
			Asset: mon.Asset,
			Parts: make([]core.FundingPart, 0),
		}
		
	}
}

func (m Interpreter) WithdrawAlways(account core.AccountAddress, mon core.Monetary) core.Funding {
	panic("todo")
}

func (m Interpreter) TakeFromSource(source program.Source, asset core.Asset) (core.Funding, *core.AccountAddress) {
	panic("todo")
}

func (m Interpreter) Repay(funding core.Funding) {
	for _, part := range funding.Parts {
		balance := m.BalanceOf(part.Account, funding.Asset)
		balance = balance.Add(part.Amount)
	}
}

func (m Interpreter) BalanceOf(account core.AccountAddress, asset core.Asset) core.Number {
	if _, ok := m.balances[account]; !ok {
		m.balances[account] = make(map[core.Asset]core.Number)
	}
	if _, ok := m.balances[account][asset]; !ok {
		m.balances[account][asset] = m.ledger.GetBalance(account, asset)
	}
	return m.balances[account][asset]
}


func (i Interpreter) Eval(expr program.Expr) core.Value {
	panic("todo")
}

func EvalAs[T core.Value](i Interpreter, expr program.Expr) T {
	x := i.Eval(expr)
	if v, ok := x.(T); ok {
		return v
	}
	panic(fmt.Errorf("unexpected type '%T'", x))
}
