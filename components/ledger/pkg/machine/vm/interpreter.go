package vm

import (
	"errors"
	"fmt"

	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/machine/vm/program"
)

const InternalError = "internal interpreter error, please report to the issue tracker"

type LedgerReader interface {
	GetBalance(account core.AccountAddress, asset core.Asset) core.Number
	GetMeta(account core.AccountAddress, key string) core.Value
}

type Machine struct {
	ledger      LedgerReader
	vars        map[string]core.Value
	balances    map[core.AccountAddress]map[core.Asset]core.Number
	Postings    []core.Posting
	TxMeta      map[string]core.Value
	AccountMeta map[core.AccountAddress]map[string]core.Value
	Printed     []core.Value
}

func NewMachine(ledger LedgerReader, vars map[string]core.Value) Machine {
	return Machine{
		ledger:      ledger,
		vars:        vars,
		balances:    make(map[core.AccountAddress]map[core.Asset]core.Number),
		Postings:    make([]core.Posting, 0),
		TxMeta:      make(map[string]core.Value),
		AccountMeta: make(map[core.AccountAddress]map[string]core.Value),
		Printed:     make([]core.Value, 0),
	}
}

func (m *Machine) Run(script program.Program) error {
	for _, var_decl := range script.VarsDecl {
		switch o := var_decl.Origin.(type) {
		case program.VarOriginMeta:
			account, err := EvalAs[core.AccountAddress](m, o.Account)
			if err != nil {
				return err
			}
			value := m.ledger.GetMeta(*account, o.Key)
			m.vars[var_decl.Name] = value
		case program.VarOriginBalance:
			account, err := EvalAs[core.AccountAddress](m, o.Account)
			if err != nil {
				return err
			}
			asset, err := EvalAs[core.Asset](m, o.Asset)
			if err != nil {
				return err
			}
			balance := m.ledger.GetBalance(*account, *asset)
			m.vars[var_decl.Name] = balance
		case nil:
		default:
			return errors.New(InternalError)
		}
	}

	for _, stmt := range script.Statements {
		switch s := stmt.(type) {
		case program.StatementFail:
			return errors.New("failed")
		case program.StatementPrint:
			v, err := m.Eval(s.Expr)
			if err != nil {
				return err
			}
			m.Printed = append(m.Printed, v)
			fmt.Printf("%v\n", s.Expr)
		case program.StatementAllocate:
			funding, err := EvalAs[core.Funding](m, s.Funding)
			if err != nil {
				return err
			}
			m.Allocate(*funding, s.Destination)
		case program.StatementLet:
			value, err := m.Eval(s.Expr)
			if err != nil {
				return err
			}
			m.vars[s.Name] = value
		case program.StatementSetTxMeta:
			value, err := m.Eval(s.Value)
			if err != nil {
				return err
			}
			m.TxMeta[s.Key] = value
		case program.StatementSetAccountMeta:
			account, err := EvalAs[core.AccountAddress](m, s.Account)
			if err != nil {
				return err
			}
			value, err := m.Eval(s.Value)
			if err != nil {
				return err
			}
			if _, ok := m.AccountMeta[*account]; !ok {
				m.AccountMeta[*account] = make(map[string]core.Value)
			}
			m.AccountMeta[*account][s.Key] = value
		default:
			return errors.New(InternalError)
		}
	}
	return nil
}

func (m *Machine) Send(funding core.Funding, account core.AccountAddress) {
	if funding.Total().Eq(core.NewNumber(0)) {
		return //no empty postings
	}
	for _, part := range funding.Parts {
		m.Postings = append(m.Postings, core.Posting{
			Source:      string(part.Account),
			Destination: string(account),
			Asset:       string(funding.Asset),
			Amount:      part.Amount,
		})
	}
}

func (m *Machine) Allocate(funding core.Funding, destination program.Destination) error {
	switch d := destination.(type) {
	case program.DestinationAccount:
		account, err := EvalAs[core.AccountAddress](m, d.Expr)
		if err != nil {
			return err
		}
		m.Send(funding, *account)
	case program.DestinationInOrder:
		for _, part := range d.Parts {
			max, err := EvalAs[core.Monetary](m, part.Max)
			if err != nil {
				return err
			}
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
				portion, err := EvalAs[core.Portion](m, part.Portion.Expr)
				if err != nil {
					return err
				}
				portions = append(portions, *portion)
			}
			sub_dests = append(sub_dests, part.Kod)
		}
		allotment, err := core.NewAllotment(portions)
		if err != nil {
			return fmt.Errorf("failed to create allotment: %v", err)
		}
		for i, part := range allotment.Allocate(funding.Total()) {
			taken, remainder, err := funding.Take(part)
			if err != nil {
				return fmt.Errorf("failed to allocate to destination: %v", err)
			}
			funding = remainder
			m.AllocateOrKeep(taken, sub_dests[i])
		}
	}
	return nil
}

func (m *Machine) AllocateOrKeep(funding core.Funding, kod program.KeptOrDestination) error {
	if kod.Kept {
		m.Repay(funding)
	} else {
		err := m.Allocate(funding, kod.Destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) TakeFromValueAwareSource(source program.ValueAwareSource, mon core.Monetary) (*core.Funding, error) {
	switch s := source.(type) {
	case program.ValueAwareSourceSource:
		available, fallback, err := m.TakeFromSource(s.Source, mon.Asset)
		if err != nil {
			return nil, fmt.Errorf("failed to take from source: %v", err)
		}
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
					return nil, errors.New("mismatching assets")
				}
			} else {
				return nil, fmt.Errorf("insufficient funds: needed %v and got %v", mon.Amount, taken.Total())
			}
		}
		m.Repay(remainder)
		return &taken, nil
	case program.ValueAwareSourceAllotment:
		portions := make([]core.Portion, 0)
		for _, part := range s {
			if part.Portion.Remaining {
				portions = append(portions, core.NewPortionRemaining())
			} else {
				portion, err := EvalAs[core.Portion](m, part.Portion.Expr)
				if err != nil {
					return nil, err
				}
				portions = append(portions, *portion)
			}
		}
		allotment, err := core.NewAllotment(portions)
		if err != nil {
			return nil, fmt.Errorf("could not create allotment: %v", err)
		}
		funding := core.Funding{
			Asset: mon.Asset,
			Parts: make([]core.FundingPart, 0),
		}
		for i, amt := range allotment.Allocate(mon.Amount) {
			taken, err := m.TakeFromValueAwareSource(program.ValueAwareSourceSource{Source: s[i].Source}, core.Monetary{Asset: mon.Asset, Amount: amt})
			if err != nil {
				return nil, fmt.Errorf("failed to take from source: %v", err)
			}
			funding, err = funding.Concat(*taken)
			if err != nil {
				return nil, fmt.Errorf("funding error: %v", err)
			}
		}
		return &funding, nil
	}
	return nil, errors.New(InternalError)
}

func (m *Machine) TakeFromSource(source program.Source, asset core.Asset) (*core.Funding, *core.AccountAddress, error) {
	switch s := source.(type) {
	case program.SourceAccount:
		account, err := EvalAs[core.AccountAddress](m, s.Account)
		if err != nil {
			return nil, nil, err
		}
		overdraft := core.Monetary{
			Asset:  asset,
			Amount: core.NewNumber(0),
		}
		var fallback *core.AccountAddress
		if s.Overdraft != nil {
			if s.Overdraft.Unbounded {
				fallback = account
			} else {
				ov, err := EvalAs[core.Monetary](m, *s.Overdraft.UpTo)
				if err != nil {
					return nil, nil, err
				}
				overdraft = *ov
			}
		}
		if string(*account) == "world" {
			fallback = account
		}
		if overdraft.Asset != asset {
			return nil, nil, errors.New("mismatching asset")
		}
		funding := m.WithdrawAll(*account, asset, overdraft.Amount)
		return &funding, fallback, nil
	case program.SourceMaxed:
		taken, fallback, err := m.TakeFromSource(s.Source, asset)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to take from source: %v", err)
		}
		max, err := EvalAs[core.Monetary](m, s.Max)
		if err != nil {
			return nil, nil, err
		}
		if max.Asset != asset {
			return nil, nil, errors.New("mismatching asset")
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
					return nil, nil, fmt.Errorf("funding error: %v", err)
				}
			}
		}
		return &maxed, nil, nil
	case program.SourceInOrder:
		total := core.Funding{
			Asset: asset,
			Parts: make([]core.FundingPart, 0),
		}
		var fallback *core.AccountAddress
		nb_sources := len(s)
		for i, source := range s {
			subsource_taken, subsource_fallback, err := m.TakeFromSource(source, asset)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to take from source: %v", err)
			}
			if subsource_fallback != nil && i != nb_sources-1 {
				return nil, nil, errors.New("fallback is not in the last position") // shouldn't we let this slide?
			}
			fallback = subsource_fallback
			total, err = total.Concat(*subsource_taken)
			if err != nil {
				return nil, nil, errors.New("mismatching assets")
			}
		}
		return &total, fallback, nil
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
	}
	return nil, nil, errors.New(InternalError)
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

func (m *Machine) Eval(expr program.Expr) (core.Value, error) {
	switch expr := expr.(type) {
	case program.ExprLiteral:
		return expr.Value, nil
	case program.ExprInfix:
		switch expr.Op {
		case program.OP_ADD:
			lhs, err := EvalAs[core.Number](m, expr.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err := EvalAs[core.Number](m, expr.Rhs)
			if err != nil {
				return nil, err
			}
			return (*lhs).Add(*rhs), nil
		case program.OP_SUB:
			lhs, err := EvalAs[core.Number](m, expr.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err := EvalAs[core.Number](m, expr.Rhs)
			if err != nil {
				return nil, err
			}
			return (*lhs).Sub(*rhs), nil
		}
	case program.ExprMonetaryNew:
		asset, err := EvalAs[core.Asset](m, expr.Asset)
		if err != nil {
			return nil, err
		}
		amount, err := EvalAs[core.Number](m, expr.Amount)
		if err != nil {
			return nil, err
		}
		return core.Monetary{
			Asset:  *asset,
			Amount: *amount,
		}, nil
	case program.ExprVariable:
		return m.vars[string(expr)], nil
	case program.ExprTake:
		amt, err := EvalAs[core.Monetary](m, expr.Amount)
		if err != nil {
			return nil, err
		}
		taken, err := m.TakeFromValueAwareSource(expr.Source, *amt)
		if err != nil {
			return nil, fmt.Errorf("failed to take from source: %v", err)
		}
		return *taken, nil
	case program.ExprTakeAll:
		asset, err := EvalAs[core.Asset](m, expr.Asset)
		if err != nil {
			return nil, err
		}
		funding, fallback, err := m.TakeFromSource(expr.Source, *asset)
		if err != nil {
			return nil, fmt.Errorf("failed to take from source: %v", err)
		}
		if fallback != nil {
			panic("oops infinite money")
		}
		return *funding, nil
	}
	return nil, errors.New(InternalError)
}

func EvalAs[T core.Value](i *Machine, expr program.Expr) (*T, error) {
	x, err := i.Eval(expr)
	if err != nil {
		return nil, err
	}
	if v, ok := x.(T); ok {
		return &v, nil
	}
	return nil, fmt.Errorf("internal interpreter error: extected type '%T' and got '%T'", *new(T), x)
}
