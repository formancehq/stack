package order

 import(
	. "ledgerorder/utils"

 ) 
type Transaction struct {
    Source      *AccountIdentifier
    Destination *AccountIdentifier
    Asset       string
    Amount      MicroAmount
}

func (t *Transaction) SetSource(acc *AccountIdentifier) *Transaction{
	t.Source = acc
	return t
}

func (t *Transaction) SetDestination(acc *AccountIdentifier) *Transaction{
	t.Destination = acc
	return t
}

func (t *Transaction) SetAsset(asset string) *Transaction{
	t.Asset = asset 
	return t 
}

func (t *Transaction) SetAmount(a MicroAmount) *Transaction{
	t.Amount = a
	return t 
}

type MicroAmount struct {
    Value           int64
    VariableAmount  *VariableAmount
}

func (m MicroAmount) FromValue(v int64) MicroAmount{
	return MicroAmount{Value: v}
}

func (m MicroAmount) FromVariable(v *VariableAmount) MicroAmount{
	return MicroAmount{VariableAmount: v}
}

