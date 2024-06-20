package order 


type Order struct {
    Spacename string
    Ledger string 
    ConditionOrder   *ConditionOrder
    VariableAmounts  []*VariableAmount
    Transactions     []TransactionConditionPair
}

type TransactionConditionPair struct {
    Condition   *ConditionTransaction
    Transaction *Transaction
}

func NewOrder() *Order {
    return &Order{}
}

func (o *Order) SetSpaceName(name string){
    o.Spacename = name 
}

func (o *Order) SetLedgerName(name string){
    o.Ledger = name 
}

func (o *Order) AddConditionOrder(cond *ConditionOrder) {
    o.ConditionOrder = cond
}

func (o *Order) AddVariableAmount(varAmount *VariableAmount) {
    o.VariableAmounts = append(o.VariableAmounts, varAmount)
}

func (o *Order) AddTransaction(trans *Transaction) {
    o.Transactions = append(o.Transactions, TransactionConditionPair{nil, trans})
}

func (o *Order) AddTransactionWithCondition(trans *Transaction, cond *ConditionTransaction){
    o.Transactions = append(o.Transactions, TransactionConditionPair{cond, trans})
}

