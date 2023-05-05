package program

import (
	"github.com/numary/ledger/pkg/core"
)

const (
	OP_ADD = byte(iota + 1)
	OP_SUB
)

type Expr interface {
	isExpr()
}

type ExprLiteral struct{ core.Value }

func (e ExprLiteral) isExpr() {}

type ExprInfix struct {
	Op  byte
	Lhs Expr
	Rhs Expr
}

func (e ExprInfix) isExpr() {}

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

type SourceArrayInOrder struct {
	Array Expr
}

func (s SourceArrayInOrder) isSource() {}

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

type DestinationAccount struct{ Expr }

func (d DestinationAccount) isDestination() {}

type DestinationInOrder struct {
	Parts []struct {
		Max Expr
		Kod KeptOrDestination
	}
	Remaining KeptOrDestination
}

func (d DestinationInOrder) isDestination() {}

type DestinationAllotment []struct {
	Portion AllotmentPortion
	Kod     KeptOrDestination
}

func (d DestinationAllotment) isDestination() {}

type Statement interface {
	isStatement()
}

type StatementFail struct{}

func (s StatementFail) isStatement() {}

type StatementPrint struct{ Expr }

func (s StatementPrint) isStatement() {}

type StatementAllocate struct {
	Funding     Expr
	Destination Destination
}

func (s StatementAllocate) isStatement() {}

type StatementLet struct {
	Name string
	Expr Expr
}

func (s StatementLet) isStatement() {}

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
	Ty     core.Type
	Name   string
	Origin VarOrigin
}

type Program struct {
	VarsDecl   []VarDecl
	Statements []Statement
}

// func (p *Program) ParseVariables(vars map[string]core.Value) (map[string]core.Value, error) {
// 	variables := make(map[string]core.Value)
// 	for _, res := range p.Resources {
// 		if variable, ok := res.(Variable); ok {
// 			if val, ok := vars[variable.Name]; ok && val.GetType() == variable.Typ {
// 				variables[variable.Name] = val
// 				switch val.GetType() {
// 				case core.TypeAccount:
// 					if err := core.ParseAccountAddress(val.(core.AccountAddress)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, string(val.(core.AccountAddress)))
// 					}
// 				case core.TypeAsset:
// 					if err := core.ParseAsset(val.(core.Asset)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, string(val.(core.Asset)))
// 					}
// 				case core.TypeMonetary:
// 					if err := core.ParseMonetary(val.(core.Monetary)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, val.(core.Monetary).String())
// 					}
// 				case core.TypePortion:
// 					if err := core.ValidatePortionSpecific(val.(core.Portion)); err != nil {
// 						return nil, errors.Wrapf(err, "invalid variable $%s value '%s'",
// 							variable.Name, val.(core.Portion).String())
// 					}
// 				case core.TypeString:
// 				case core.TypeNumber:
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

// func (p *Program) ParseVariablesJSON(vars map[string]json.RawMessage) (map[string]core.Value, error) {
// 	variables := make(map[string]core.Value)
// 	for _, res := range p.Resources {
// 		if param, ok := res.(Variable); ok {
// 			data, ok := vars[param.Name]
// 			if !ok {
// 				return nil, fmt.Errorf("missing variable $%s", param.Name)
// 			}
// 			val, err := core.NewValueFromJSON(param.Typ, data)
// 			if err != nil {
// 				return nil, fmt.Errorf(
// 					"invalid JSON value for variable $%s of type %v: %w",
// 					param.Name, param.Typ, err)
// 			}
// 			variables[param.Name] = *val
// 			delete(vars, param.Name)
// 		}
// 	}
// 	for name := range vars {
// 		return nil, fmt.Errorf("extraneous variable $%s", name)
// 	}
// 	return variables, nil
// }
