package ledgerorder 

import (
	"ledgerorder/order"
	"ledgerorder/utils"
)


func Order() *order.Order {
	return &order.Order{}
}

func Transaction() *order.Transaction {
	return &order.Transaction{}
}

func AccountIdentifier(org, ledger, account string) *utils.AccountIdentifier {
	return &utils.AccountIdentifier{OrganisationName: org, Ledger: ledger, AccountName: account}
}

var MicroAmount order.MicroAmount = order.MicroAmount{}

var VariableAmount order.VariableAmount = order.VariableAmount{}

var ValueExpression utils.ValueExpression = utils.ValueExpression{}

var LogicalExpression utils.LogicalExpression = utils.LogicalExpression{}

var Operators = utils.Operators


var ConditionTransaction order.ConditionTransaction = order.ConditionTransaction{}

var ConditionOrder order.ConditionOrder = order.ConditionOrder{}