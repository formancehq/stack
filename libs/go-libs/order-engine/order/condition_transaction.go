package order

import . "ledgerorder/utils"

type ConditionTransaction struct {
    Expression *LogicalExpression
}


func NewConditionTransaction(expression *LogicalExpression) *ConditionTransaction {
    return &ConditionTransaction{Expression: expression}
}

func (ConditionTransaction) With(expression *LogicalExpression) *ConditionTransaction {
    return &ConditionTransaction{Expression: expression}
}