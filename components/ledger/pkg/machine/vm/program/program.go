package program

import (
	"fmt"

	"github.com/formancehq/ledger/pkg/machine/internal"
)

const (
	OP_ADD = byte(iota + 1)
	OP_SUB
)

type Expr interface {
	isExpr()
}

type ExprLiteral struct {
	Value internal.Value
}

func (e ExprLiteral) isExpr() {}

type ExprNumberAdd struct {
	Lhs Expr
	Rhs Expr
}

func (e ExprNumberAdd) isExpr() {}

type ExprNumberSub struct {
	Lhs Expr
	Rhs Expr
}

func (e ExprNumberSub) isExpr() {}

type ExprMonetaryAdd struct {
	Lhs Expr
	Rhs Expr
}

func (e ExprMonetaryAdd) isExpr() {}

type ExprMonetarySub struct {
	Lhs Expr
	Rhs Expr
}

func (e ExprMonetarySub) isExpr() {}

type ExprVariable string

func (e ExprVariable) isExpr() {}

type ExprTake struct {
	Amount Expr
	Source ValueAwareSource
}

func (e ExprTake) isExpr() {}

type ExprTakeAll struct {
	Asset  Expr
	Source Source
}

func (e ExprTakeAll) isExpr() {}

type ExprMonetaryNew struct {
	Asset  Expr
	Amount Expr
}

func (e ExprMonetaryNew) isExpr() {}

type Overdraft struct {
	Unbounded bool
	UpTo      *Expr // invariant: if unbounbed then up_to == nil
}

type Source interface {
	isSource()
}

type SourceAccount struct {
	Account   Expr
	Overdraft *Overdraft
}

func (s SourceAccount) isSource() {}

type SourceMaxed struct {
	Source Source
	Max    Expr
}

func (s SourceMaxed) isSource() {}

type SourceInOrder []Source

func (s SourceInOrder) isSource() {}

// invariant: if remaining then expr == nil
type AllotmentPortion struct {
	Expr      Expr
	Remaining bool
}

type ValueAwareSource interface {
	isValueAwareSource()
}

type ValueAwareSourceSource struct {
	Source Source
}

func (v ValueAwareSourceSource) isValueAwareSource() {}

type ValueAwareSourcePart struct {
	Portion AllotmentPortion
	Source  Source
}
type ValueAwareSourceAllotment []ValueAwareSourcePart

func (v ValueAwareSourceAllotment) isValueAwareSource() {}

type KeptOrDestination struct {
	Kept        bool
	Destination Destination
}

type Destination interface {
	isDestination()
}

type DestinationAccount struct{ Expr Expr }

func (d DestinationAccount) isDestination() {}

type DestinationInOrderPart struct {
	Max Expr
	Kod KeptOrDestination
}

type DestinationInOrder struct {
	Parts     []DestinationInOrderPart
	Remaining KeptOrDestination
}

func (d DestinationInOrder) isDestination() {}

type DestinationAllotmentPart struct {
	Portion AllotmentPortion
	Kod     KeptOrDestination
}
type DestinationAllotment []DestinationAllotmentPart

func (d DestinationAllotment) isDestination() {}

type Statement interface {
	isStatement()
}

type StatementFail struct{}

func (s StatementFail) isStatement() {}

type StatementPrint struct{ Expr Expr }

func (s StatementPrint) isStatement() {}

type StatementSave struct {
	Amount  Expr
	Account Expr
}

func (s StatementSave) isStatement() {}

type StatementSaveAll struct {
	Asset   Expr
	Account Expr
}

func (s StatementSaveAll) isStatement() {}

type StatementAllocate struct {
	Funding     Expr
	Destination Destination
}

func (s StatementAllocate) isStatement() {}

type StatementSetTxMeta struct {
	Key   string
	Value Expr
}

func (s StatementSetTxMeta) isStatement() {}

type StatementSetAccountMeta struct {
	Account Expr
	Key     string
	Value   Expr
}

func (s StatementSetAccountMeta) isStatement() {}

type VarOrigin interface {
	isVarOrigin()
}

type VarOriginMeta struct {
	Account Expr
	Key     string
}

func (v VarOriginMeta) isVarOrigin() {}

type VarOriginBalance struct {
	Account Expr
	Asset   Expr
}

func (v VarOriginBalance) isVarOrigin() {}

type VarDecl struct {
	Typ    internal.Type
	Name   string
	Origin VarOrigin
}

type Program struct {
	VarsDecl   []VarDecl
	Statements []Statement
}

// func (p *Program) ParseVariables(vars map[string]internal.Value) (map[string]internal.Value, error) {
// 	variables := make(map[string]internal.Value)
// 	for _, res := range p.Resources {
// 		if variable, ok := res.(Variable); ok {
// 			if val, ok := vars[variable.Name]; ok && val.GetType() == variable.Typ {
// 				variables[variable.Name] = val
// 				switch val.GetType() {
// 				case internal.TypeAccount:
// 					if err := internal.ParseAccountAddress(val.(internal.AccountAddress)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, string(val.(internal.AccountAddress)))
// 					}
// 				case internal.TypeAsset:
// 					if err := internal.ParseAsset(val.(internal.Asset)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, string(val.(internal.Asset)))
// 					}
// 				case internal.TypeMonetary:
// 					if err := internal.ParseMonetary(val.(internal.Monetary)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, val.(internal.Monetary).String())
// 					}
// 				case internal.TypePortion:
// 					if err := internal.ValidatePortionSpecific(val.(internal.Portion)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, val.(internal.Portion).String())
// 					}
// 				case internal.TypeString:
// 				case internal.TypeNumber:
// 				default:
// 					return nil, fmt.Errorf("unsupported type for variable $%s: %s",
// 						variable.Name, val.GetType())
// 				}
// 				delete(vars, variable.Name)
// 			} else if val, ok := vars[variable.Name]; ok && val.GetType() != variable.Typ {
// 				return nil, fmt.Errorf("wrong type for variable $%s: %s instead of %s",
// 					variable.Name, variable.Typ, val.GetType())
// 			} else {
// 				return nil, fmt.Errorf("missing variable $%s", variable.Name)
// 			}
// 		}
// 	}
// 	for name := range vars {
// 		return nil, fmt.Errorf("extraneous variable $%s", name)
// 	}
// 	return variables, nil
// }

func (p *Program) ParseVariablesJSON(vars map[string]string) (map[string]internal.Value, error) {
	variables := make(map[string]internal.Value)
	for _, var_decl := range p.VarsDecl {
		if var_decl.Origin != nil {
			continue
		}
		data, ok := vars[var_decl.Name]
		if !ok {
			return nil, fmt.Errorf("missing variable $%s", var_decl.Name)
		}
		val, err := internal.NewValueFromString(var_decl.Typ, data)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid JSON value for variable $%s of type %v: %w",
				var_decl.Name, var_decl.Typ, err)
		}
		variables[var_decl.Name] = val
		delete(vars, var_decl.Name)
	}
	for name := range vars {
		return nil, fmt.Errorf("extraneous variable $%s", name)
	}
	return variables, nil
}
