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
	balances []core.Posting
	tx_meta map[string]core.Value
	account_meta map[core.AccountAddress]map[string]core.Value
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
			// ice
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
		case program.StatementTxMeta:
			value := m.Eval(s.Expr)
			m.tx_meta[s.Key] = value
		case program.SetAccountMeta:
			account := EvalAs(m, s.Account)
			value := m.Eval(s.Value)
			if 
			account_meta := []
		}
	}
}

func (m Interpreter) Allocate(funding core.Funding, destination program.Destination) {
	// todo	
}


func (i Interpreter) Eval(expr program.Expr) core.Value {
	// todo
}

func EvalAs[T core.Value](i Interpreter, expr program.Expr) T {
	x := i.Eval(expr)
	if v, ok := x.(T); ok {
		return v
	}
	panic(fmt.Errorf("unexpected type '%T'", x))
}
