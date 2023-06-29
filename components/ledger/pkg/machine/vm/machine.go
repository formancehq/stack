package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/formancehq/ledger/pkg/machine/internal"
	"github.com/formancehq/ledger/pkg/machine/vm/program"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Posting struct {
	Source      string                `json:"source"`
	Destination string                `json:"destination"`
	Amount      *internal.MonetaryInt `json:"amount"`
	Asset       string                `json:"asset"`
}

const InternalError = "internal interpreter error, please report to the issue tracker"

// type LedgerReader interface {
// 	GetBalance(account internal.AccountAddress, asset internal.Asset) internal.Number
// 	GetMeta(account internal.AccountAddress, key string) internal.Value
// }

type Machine struct {
	store        Store
	ctx          context.Context
	providedVars map[string]string
	vars         map[string]internal.Value
	balances     map[internal.AccountAddress]map[internal.Asset]internal.Number
	Postings     []Posting
	TxMeta       map[string]internal.Value
	AccountMeta  map[internal.AccountAddress]map[string]internal.Value
	Printed      []internal.Value
}

func NewMachine(store Store, vars map[string]string) Machine {
	return Machine{
		store:        store,
		providedVars: vars,
		balances:     make(map[internal.AccountAddress]map[internal.Asset]internal.Number),
		Postings:     make([]Posting, 0),
		TxMeta:       make(map[string]internal.Value),
		AccountMeta:  make(map[internal.AccountAddress]map[string]internal.Value),
		Printed:      make([]internal.Value, 0),
	}
}

func (m *Machine) checkVar(value internal.Value) error {
	switch v := value.(type) {
	case internal.Monetary:
		if v.Amount.Ltz() {
			return fmt.Errorf("monetary amounts must be non-negative but is %v", v.Amount)
		}
	}
	return nil
}

func (m *Machine) Execute(script program.Program) error {
	for _, var_decl := range script.VarsDecl {
		switch o := var_decl.Origin.(type) {
		case program.VarOriginMeta:
			account, err := EvalAs[internal.AccountAddress](m, o.Account)
			if err != nil {
				return err
			}
			metadata, err := m.store.GetMetadataFromLogs(m.ctx, string(*account), o.Key)
			if err != nil {
				return fmt.Errorf("failed to get metadata of account %s for key %s: %s", account, o.Key, err)
			}
			value, err := internal.NewValueFromString(var_decl.Typ, metadata)
			if err != nil {
				return fmt.Errorf("failed to parse variable: %s", err)
			}
			err = m.checkVar(value)
			if err != nil {
				return fmt.Errorf("failed to get metadata of account %s for key %s: %s", account, o.Key, err)
			}
			m.vars[var_decl.Name] = value
		case program.VarOriginBalance:
			account, err := EvalAs[internal.AccountAddress](m, o.Account)
			if err != nil {
				return err
			}
			asset, err := EvalAs[internal.Asset](m, o.Asset)
			if err != nil {
				return err
			}
			amt, err := m.store.GetBalanceFromLogs(m.ctx, string(*account), string(*asset))
			if err != nil {
				return err
			}
			balance := internal.Monetary{
				Asset:  *asset,
				Amount: internal.NewMonetaryIntFromBigInt(amt),
			}
			err = m.checkVar(balance)
			if err != nil {
				return fmt.Errorf("failed to get balance of account %s for asset %s: %s", account, asset, err)
			}
			m.vars[var_decl.Name] = balance
		case nil:
			val, err := internal.NewValueFromString(var_decl.Typ, m.providedVars[var_decl.Name])
			if err != nil {
				return fmt.Errorf("failed to parse variable: %s", err)
			}
			err = m.checkVar(val)
			if err != nil {
				return fmt.Errorf("variable passed is incorrect: %s", err)
			}
			m.vars[var_decl.Name] = val
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
			funding, err := EvalAs[internal.Funding](m, s.Funding)
			if err != nil {
				return err
			}
			err = m.Allocate(*funding, s.Destination)
			if err != nil {
				return err
			}
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
			account, err := EvalAs[internal.AccountAddress](m, s.Account)
			if err != nil {
				return err
			}
			value, err := m.Eval(s.Value)
			if err != nil {
				return err
			}
			if _, ok := m.AccountMeta[*account]; !ok {
				m.AccountMeta[*account] = make(map[string]internal.Value)
			}
			m.AccountMeta[*account][s.Key] = value
		default:
			return errors.New(InternalError)
		}
	}
	return nil
}

func (m *Machine) Send(funding internal.Funding, account internal.AccountAddress) error {
	if funding.Total().Eq(internal.NewNumber(0)) {
		return nil //no empty postings
	}
	for _, part := range funding.Parts {
		m.Postings = append(m.Postings, Posting{
			Source:      string(part.Account),
			Destination: string(account),
			Asset:       string(funding.Asset),
			Amount:      part.Amount,
		})
		bal, err := m.BalanceOf(account, funding.Asset)
		if err != nil {
			return err
		}
		*bal = *bal.Add(part.Amount)
	}
	return nil
}

func (m *Machine) Allocate(funding internal.Funding, destination program.Destination) error {
	switch d := destination.(type) {
	case program.DestinationAccount:
		account, err := EvalAs[internal.AccountAddress](m, d.Expr)
		if err != nil {
			return err
		}
		m.Send(funding, *account)
	case program.DestinationInOrder:
		for _, part := range d.Parts {
			max, err := EvalAs[internal.Monetary](m, part.Max)
			if err != nil {
				return err
			}
			taken, remainder := funding.TakeMax(max.Amount)
			funding = remainder
			m.AllocateOrKeep(taken, part.Kod)
		}
		m.AllocateOrKeep(funding, d.Remaining)
	case program.DestinationAllotment:
		portions := make([]internal.Portion, 0)
		sub_dests := make([]program.KeptOrDestination, 0)
		for _, part := range d {
			if part.Portion.Remaining {
				portions = append(portions, internal.NewPortionRemaining())
			} else {
				portion, err := EvalAs[internal.Portion](m, part.Portion.Expr)
				if err != nil {
					return err
				}
				portions = append(portions, *portion)
			}
			sub_dests = append(sub_dests, part.Kod)
		}
		allotment, err := internal.NewAllotment(portions)
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

func (m *Machine) AllocateOrKeep(funding internal.Funding, kod program.KeptOrDestination) error {
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

func (m *Machine) TakeFromValueAwareSource(source program.ValueAwareSource, mon internal.Monetary) (*internal.Funding, error) {
	switch s := source.(type) {
	case program.ValueAwareSourceSource:
		available, fallback, err := m.TakeFromSource(s.Source, mon.Asset)
		if err != nil {
			return nil, fmt.Errorf("failed to take from source: %v", err)
		}
		taken, remainder := available.TakeMax(mon.Amount)
		if !taken.Total().Eq(mon.Amount) {
			missing := internal.Monetary{
				Asset:  mon.Asset,
				Amount: mon.Amount.Sub(taken.Total()),
			}
			if fallback != nil {
				missing_taken, err := m.WithdrawAlways(*fallback, missing)
				if err != nil {
					return nil, err
				}
				taken, err = taken.Concat(*missing_taken)
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
		portions := make([]internal.Portion, 0)
		for _, part := range s {
			if part.Portion.Remaining {
				portions = append(portions, internal.NewPortionRemaining())
			} else {
				portion, err := EvalAs[internal.Portion](m, part.Portion.Expr)
				if err != nil {
					return nil, err
				}
				portions = append(portions, *portion)
			}
		}
		allotment, err := internal.NewAllotment(portions)
		if err != nil {
			return nil, fmt.Errorf("could not create allotment: %v", err)
		}
		funding := internal.Funding{
			Asset: mon.Asset,
			Parts: make([]internal.FundingPart, 0),
		}
		for i, amt := range allotment.Allocate(mon.Amount) {
			taken, err := m.TakeFromValueAwareSource(program.ValueAwareSourceSource{Source: s[i].Source}, internal.Monetary{Asset: mon.Asset, Amount: amt})
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

func (m *Machine) TakeFromSource(source program.Source, asset internal.Asset) (*internal.Funding, *internal.AccountAddress, error) {
	switch s := source.(type) {
	case program.SourceAccount:
		account, err := EvalAs[internal.AccountAddress](m, s.Account)
		if err != nil {
			return nil, nil, err
		}
		overdraft := internal.Monetary{
			Asset:  asset,
			Amount: internal.NewNumber(0),
		}
		var fallback *internal.AccountAddress
		if s.Overdraft != nil {
			if s.Overdraft.Unbounded {
				fallback = account
			} else {
				ov, err := EvalAs[internal.Monetary](m, *s.Overdraft.UpTo)
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
		funding, err := m.WithdrawAll(*account, asset, overdraft.Amount)
		if err != nil {
			return nil, nil, err
		}
		return funding, fallback, nil
	case program.SourceMaxed:
		taken, fallback, err := m.TakeFromSource(s.Source, asset)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to take from source: %v", err)
		}
		max, err := EvalAs[internal.Monetary](m, s.Max)
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
				missing := internal.Monetary{
					Asset:  asset,
					Amount: max.Amount.Sub(maxed.Total()),
				}
				withdrawn, err := m.WithdrawAlways(*fallback, missing)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to withdraw %s", err)
				}
				maxed, err = maxed.Concat(*withdrawn)
				if err != nil {
					return nil, nil, fmt.Errorf("funding error: %v", err)
				}
			}
		}
		return &maxed, nil, nil
	case program.SourceInOrder:
		total := internal.Funding{
			Asset: asset,
			Parts: make([]internal.FundingPart, 0),
		}
		var fallback *internal.AccountAddress
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
		// 	list := EvalAs[internal.AccountAddress](m, s.Array)
		// 	total := internal.Funding{
		// 		Asset: asset,
		// 		Parts: make([]internal.FundingPart, 0),
		// 	}
		// 	for _, account := range list {
		// 		withdrawn := m.WithdrawAll(account, asset, internal.NewNumber(0))
		// 		total, err = total.Concat(withdrawn)
		// 		if err != nil {
		// 			panic("mismatching assets")
		// 		}
		// 	}
		// 	return total, nil
	}
	return nil, nil, errors.New(InternalError)
}

func (m *Machine) Repay(funding internal.Funding) error {
	for _, part := range funding.Parts {
		balance, err := m.BalanceOf(part.Account, funding.Asset)
		if err != nil {
			return err
		}
		*balance = *balance.Add(part.Amount)
	}
	return nil
}

func (m *Machine) WithdrawAll(account internal.AccountAddress, asset internal.Asset, overdraft internal.Number) (*internal.Funding, error) {
	balance, err := m.BalanceOf(account, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw %s", err)
	}
	withdrawn := balance.Add(overdraft)
	*balance = *overdraft.Neg()
	return &internal.Funding{
		Asset: asset,
		Parts: []internal.FundingPart{
			{
				Account: account,
				Amount:  withdrawn,
			},
		},
	}, nil
}

func (m *Machine) WithdrawAlways(account internal.AccountAddress, mon internal.Monetary) (*internal.Funding, error) {
	balance, err := m.BalanceOf(account, mon.Asset)
	if err != nil {
		return nil, err
	}
	*balance = *balance.Sub(mon.Amount)
	return &internal.Funding{
		Asset: mon.Asset,
		Parts: []internal.FundingPart{
			{
				Account: account,
				Amount:  mon.Amount,
			},
		},
	}, nil
}

func (m *Machine) BalanceOf(account internal.AccountAddress, asset internal.Asset) (internal.Number, error) {
	if _, ok := m.balances[account]; !ok {
		m.balances[account] = make(map[internal.Asset]internal.Number)
	}
	if _, ok := m.balances[account][asset]; !ok {
		amt, err := m.store.GetBalanceFromLogs(m.ctx, string(account), string(asset))
		if err != nil {
			return nil, fmt.Errorf("failed to get balance from store: %s", err)
		}
		m.balances[account][asset] = internal.NewMonetaryIntFromBigInt(amt)
	}
	return m.balances[account][asset], nil
}

func (m *Machine) Eval(expr program.Expr) (internal.Value, error) {
	switch expr := expr.(type) {
	case program.ExprLiteral:
		return expr.Value, nil
	case program.ExprInfix:
		switch expr.Op {
		case program.OP_ADD:
			lhs, err := EvalAs[internal.Number](m, expr.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err := EvalAs[internal.Number](m, expr.Rhs)
			if err != nil {
				return nil, err
			}
			return (*lhs).Add(*rhs), nil
		case program.OP_SUB:
			lhs, err := EvalAs[internal.Number](m, expr.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err := EvalAs[internal.Number](m, expr.Rhs)
			if err != nil {
				return nil, err
			}
			return (*lhs).Sub(*rhs), nil
		}
	case program.ExprMonetaryNew:
		asset, err := EvalAs[internal.Asset](m, expr.Asset)
		if err != nil {
			return nil, err
		}
		amount, err := EvalAs[internal.Number](m, expr.Amount)
		if err != nil {
			return nil, err
		}
		return internal.Monetary{
			Asset:  *asset,
			Amount: *amount,
		}, nil
	case program.ExprVariable:
		return m.vars[string(expr)], nil
	case program.ExprTake:
		amt, err := EvalAs[internal.Monetary](m, expr.Amount)
		if err != nil {
			return nil, err
		}
		taken, err := m.TakeFromValueAwareSource(expr.Source, *amt)
		if err != nil {
			return nil, fmt.Errorf("failed to take from source: %v", err)
		}
		return *taken, nil
	case program.ExprTakeAll:
		asset, err := EvalAs[internal.Asset](m, expr.Asset)
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

func EvalAs[T internal.Value](i *Machine, expr program.Expr) (*T, error) {
	x, err := i.Eval(expr)
	if err != nil {
		return nil, err
	}
	if v, ok := x.(T); ok {
		return &v, nil
	}
	return nil, fmt.Errorf("internal interpreter error: expected type '%T' and got '%T'", *new(T), x)
}

func (m *Machine) GetTxMetaJSON() metadata.Metadata {
	meta := make(metadata.Metadata)
	for k, v := range m.TxMeta {
		var err error
		meta[k], err = internal.NewStringFromValue(v)
		if err != nil {
			panic(err)
		}
	}
	return meta
}

func (m *Machine) GetAccountsMetaJSON() map[string]metadata.Metadata {
	res := make(map[string]metadata.Metadata)
	for account, meta := range m.AccountMeta {
		for k, v := range meta {
			if _, ok := res[string(account)]; !ok {
				res[string(account)] = metadata.Metadata{}
			}

			var err error
			res[string(account)][k], err = internal.NewStringFromValue(v)
			if err != nil {
				panic(err)
			}
		}
	}

	return res
}
