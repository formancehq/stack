package program

import "github.com/formancehq/ledger/pkg/machine/internal"

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

type ExprNumberOperation struct {
	Op  byte
	Lhs Expr
	Rhs Expr
}

func (e ExprNumberOperation) isExpr() {}

type ExprMonetaryOperation struct {
	Op  byte
	Lhs Expr
	Rhs Expr
}

func (e ExprMonetaryOperation) isExpr() {}

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
