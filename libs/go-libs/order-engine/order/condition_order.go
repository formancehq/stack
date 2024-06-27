package order

import . "ledgerorder/utils"

type ConditionOrder struct {
    Expression *LogicalExpression
}

func NewConditionOrder(expression *LogicalExpression) *ConditionOrder {

    return &ConditionOrder{Expression: expression}
}

func (ConditionOrder) With(expression *LogicalExpression) *ConditionOrder {
    return &ConditionOrder{Expression: expression}
}