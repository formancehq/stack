package vm

import (
	"fmt"

	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/machine/vm/program"
)

type LedgerReader interface {
	GetBalance(account core.AccountAddress, asset core.Asset) core.Number
	GetMeta(account core.AccountAddress, key string) core.Value
}

type Machine struct {
	ledger      LedgerReader
	vars        map[string]core.Value
	balances    map[core.AccountAddress]map[core.Asset]core.Number
	TxMeta      map[string]core.Value
	AccountMeta map[core.AccountAddress]map[string]core.Value
	Postings    []core.Posting
}

func NewMachine(ledger LedgerReader, vars map[string]core.Value) Machine {
	return Machine{
		ledger:      ledger,
		vars:        vars,
		balances:    make(map[core.AccountAddress]map[core.Asset]core.Number),
		Postings:    make([]core.Posting, 0),
		TxMeta:      make(map[string]core.Value),
		AccountMeta: make(map[core.AccountAddress]map[string]core.Value),
	}
}

func (m *Machine) Run(script program.Program) {
	for _, var_decl := range script.VarsDecl {
		switch o := var_decl.Origin.(type) {
		case program.VarOriginMeta:
			account := EvalAs[core.AccountAddress](m, o.Account)
			value := m.ledger.GetMeta(account, o.Key)
			m.vars[var_decl.Name] = value
		case program.VarOriginBalance:
			account := EvalAs[core.AccountAddress](m, o.Account)
			asset := EvalAs[core.Asset](m, o.Asset)
			balance := m.ledger.GetBalance(account, asset)
			m.vars[var_decl.Name] = balance
		case nil:
		default:
			panic("internal error")
		}
	}

	for _, stmt := range script.Statements {
		switch s := stmt.(type) {
		case *program.StatementFail:
			panic("fail")
		case *program.StatementPrint:
			fmt.Printf("%v\n", s.Expr)
		case *program.StatementAllocate:
			funding := EvalAs[core.Funding](m, s.Funding)
			m.Allocate(funding, s.Destination)
		case *program.StatementLet:
			value := m.Eval(s.Expr)
			m.vars[s.Name] = value
		case *program.StatementSetTxMeta:
			value := m.Eval(s.Value)
			m.TxMeta[s.Key] = value
		case *program.StatementSetAccountMeta:
			account := EvalAs[core.AccountAddress](m, s.Account)
			value := m.Eval(s.Value)
			if _, ok := m.AccountMeta[account]; !ok {
				m.AccountMeta[account] = make(map[string]core.Value)
			}
			m.AccountMeta[account][s.Key] = value
		default:
			panic("internal error")
		}
	}
}

func (m *Machine) Send(funding core.Funding, account core.AccountAddress) {
	for _, part := range funding.Parts {
		m.Postings = append(m.Postings, core.Posting{
			Source:      string(part.Account),
			Destination: string(account),
			Asset:       string(funding.Asset),
			Amount:      part.Amount,
		})
	}
}

func (m *Machine) Allocate(funding core.Funding, destination program.Destination) {
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
				portion := EvalAs[core.Portion](m, part.Portion.Expr)
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

func (m *Machine) AllocateOrKeep(funding core.Funding, kod program.KeptOrDestination) {
	if kod.Kept {
		m.Repay(funding)
	} else {
		m.Allocate(funding, kod.Destination)
	}
}

func (m *Machine) TakeFromValueAwareSource(source program.ValueAwareSource, mon core.Monetary) core.Funding {
	var err error
	switch s := source.(type) {
	case program.ValueAwareSourceSource:
		available, fallback := m.TakeFromSource(s.Source, mon.Asset)
		taken, remainder := available.TakeMax(mon.Amount)
		if !taken.Total().Eq(mon.Amount) {
			missing := core.Monetary{
				Asset:  mon.Asset,
				Amount: mon.Amount.Sub(taken.Total()),
			}
			if fallback != nil {
				missing_taken := m.WithdrawAlways(*fallback, missing)
				taken, err = taken.Concat(missing_taken)
				if err != nil {
					panic("??")
				}
			} else {
				panic(fmt.Sprintf("insufficient funds: needed %v and got %v", mon.Amount, taken.Total()))
			}
		}
		m.Repay(remainder)
		return taken
	case program.ValueAwareSourceAllotment:
		portions := make([]core.Portion, 0)
		for _, part := range s {
			if part.Portion.Remaining {
				portions = append(portions, core.NewPortionRemaining())
			} else {
				portion := EvalAs[core.Portion](m, part.Portion.Expr)
				portions = append(portions, portion)
			}
		}
		allotment, err := core.NewAllotment(portions)
		if err != nil {
			panic("invalid allotment")
		}
		funding := core.Funding{
			Asset: mon.Asset,
			Parts: make([]core.FundingPart, 0),
		}
		for i, amt := range allotment.Allocate(mon.Amount) {
			taken := m.TakeFromValueAwareSource(program.ValueAwareSourceSource{Source: s[i].Source}, core.Monetary{Asset: mon.Asset, Amount: amt})
			funding, err = funding.Concat(taken)
			if err != nil {
				panic("funding error")
			}
		}
		return funding
	default:
		panic("ice")
	}
}

func (m *Machine) TakeFromSource(source program.Source, asset core.Asset) (core.Funding, *core.AccountAddress) {
	var err error
	switch s := source.(type) {
	case program.SourceAccount:
		account := EvalAs[core.AccountAddress](m, s.Account)
		overdraft := core.Monetary{
			Asset:  asset,
			Amount: core.NewNumber(0),
		}
		var fallback *core.AccountAddress
		if s.Overdraft != nil {
			if s.Overdraft.Unbounded {
				fallback = &account
			} else {
				overdraft = EvalAs[core.Monetary](m, *s.Overdraft.UpTo)
			}
		}
		if overdraft.Asset != asset {
			panic("mismatching asset")
		}
		funding := m.WithdrawAll(account, asset, overdraft.Amount)
		return funding, fallback
	case program.SourceMaxed:
		taken, fallback := m.TakeFromSource(s.Source, asset)
		max := EvalAs[core.Monetary](m, s.Max)
		if max.Asset != asset {
			panic("mismatching assets")
		}
		maxed, remainder := taken.TakeMax(max.Amount)
		m.Repay(remainder)
		if maxed.Total().Lte(max.Amount) {
			if fallback != nil {
				missing := core.Monetary{
					Asset:  asset,
					Amount: max.Amount.Sub(maxed.Total()),
				}
				maxed, err = maxed.Concat(m.WithdrawAlways(*fallback, missing))
				if err != nil {
					panic("")
				}
			}
		}
		return maxed, nil
	case program.SourceInOrder:
		total := core.Funding{
			Asset: asset,
			Parts: make([]core.FundingPart, 0),
		}
		var fallback *core.AccountAddress
		nb_sources := len(s)
		for i, source := range s {
			subsource_taken, subsource_fallback := m.TakeFromSource(source, asset)
			if subsource_fallback != nil && i != nb_sources-1 {
				panic("fallback not in last position")
			}
			fallback = subsource_fallback
			total, err = total.Concat(subsource_taken)
			if err != nil {
				panic("mismatching assets")
			}
		}
		return total, fallback
	// case program.SourceArrayInOrder:
	// 	list := EvalAs[core.AccountAddress](m, s.Array)
	// 	total := core.Funding{
	// 		Asset: asset,
	// 		Parts: make([]core.FundingPart, 0),
	// 	}
	// 	for _, account := range list {
	// 		withdrawn := m.WithdrawAll(account, asset, core.NewNumber(0))
	// 		total, err = total.Concat(withdrawn)
	// 		if err != nil {
	// 			panic("mismatching assets")
	// 		}
	// 	}
	// 	return total, nil
	default:
		panic("ice")
	}
}

func (m *Machine) Repay(funding core.Funding) {
	for _, part := range funding.Parts {
		balance := m.BalanceOf(part.Account, funding.Asset)
		*balance = *balance.Add(part.Amount)
	}
}

func (m *Machine) WithdrawAll(account core.AccountAddress, asset core.Asset, overdraft core.Number) core.Funding {
	balance := m.BalanceOf(account, asset)
	withdrawn := balance.Sub(overdraft)
	*balance = *balance.Sub(overdraft)
	return core.Funding{
		Asset: asset,
		Parts: []core.FundingPart{
			{
				Account: account,
				Amount:  withdrawn,
			},
		},
	}
}

func (m *Machine) WithdrawAlways(account core.AccountAddress, mon core.Monetary) core.Funding {
	balance := m.BalanceOf(account, mon.Asset)
	*balance = *balance.Sub(mon.Amount)
	return core.Funding{
		Asset: mon.Asset,
		Parts: []core.FundingPart{
			{
				Account: account,
				Amount:  mon.Amount,
			},
		},
	}
}

func (m *Machine) BalanceOf(account core.AccountAddress, asset core.Asset) core.Number {
	if _, ok := m.balances[account]; !ok {
		m.balances[account] = make(map[core.Asset]core.Number)
	}
	if _, ok := m.balances[account][asset]; !ok {
		m.balances[account][asset] = m.ledger.GetBalance(account, asset)
	}
	return m.balances[account][asset]
}

func (m *Machine) Eval(expr program.Expr) core.Value {
	switch expr := expr.(type) {
	case program.ExprLiteral:
		return expr.Value
	case program.ExprInfix:
		switch expr.Op {
		case program.OP_ADD:
			lhs := EvalAs[core.Number](m, expr.Lhs)
			rhs := EvalAs[core.Number](m, expr.Rhs)
			return lhs.Add(rhs)
		case program.OP_SUB:
			lhs := EvalAs[core.Number](m, expr.Lhs)
			rhs := EvalAs[core.Number](m, expr.Rhs)
			return lhs.Sub(rhs)
		}
	case program.ExprMonetaryNew:
		asset := EvalAs[core.Asset](m, expr.Asset)
		amount := EvalAs[core.Number](m, expr.Amount)
		return core.Monetary{
			Asset:  asset,
			Amount: amount,
		}
	case program.ExprVariable:
		return m.vars[string(expr)]
	case program.ExprTake:
		amt := EvalAs[core.Monetary](m, expr.Amount)
		return m.TakeFromValueAwareSource(expr.Source, amt)
	case program.ExprTakeAll:
		asset := EvalAs[core.Asset](m, expr.Asset)
		funding, fallback := m.TakeFromSource(expr.Source, asset)
		if fallback != nil {
			panic("oops infinite money")
		}
		return funding
	}
	panic("nope")
}

func EvalAs[T core.Value](i *Machine, expr program.Expr) T {
	x := i.Eval(expr)
	if v, ok := x.(T); ok {
		return v
	}
	panic(fmt.Errorf("unexpected type '%T'", x))
}
