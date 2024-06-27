package order

import(
	. "ledgerorder/utils"
)

type VariableName string 

type VariableAmount struct {
    Name       VariableName
    Expression *ValueExpression
}

func (v VariableAmount) New(name string, expression *ValueExpression) *VariableAmount {
    return &VariableAmount{Name: VariableName(name), Expression: expression}
}

func NewVariableAmount(name string, expression *ValueExpression) *VariableAmount {
    return &VariableAmount{Name: VariableName(name), Expression: expression}
}