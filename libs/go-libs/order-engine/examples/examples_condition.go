package examples

import (
	. "ledgerorder"
)


func CondExample1(){
	
    /*
		
		FICTIVE EXAMPLE BASED ON FUTUR PROBABLE FEATURES
		Conditionnal Transaction

		send [COIN 100] (
			condition = {
				$lt : { AggBalance[@player:leslieknope], 500 }	
			}

  			source = @centralbank
 			destination = {
    			25% to @player:leslieknope
    			remaining to @player:leslieknope:chest
  			}
		)


	*/

	source1 := AccountIdentifier("orga_schema", "ledger-01", "@centralbank")
	destination1 := AccountIdentifier("orga_schema", "ledger-01","@player:leslieknope")
	destination2 := AccountIdentifier("orga_schema", "ledger-01", "@player:leslieknope:chest")

	aggIdentifier := AccountIdentifier("orga_schema", "ledger-01","@player:leslieknope")


	aggBalance1 := VariableAmount.New("aggBalance1",
				ValueExpression.New(Operators.AggBalance, aggIdentifier))

	condTrans1 := ConditionTransaction.With(LogicalExpression.New(Operators.Lt, aggBalance1, MicroAmount.FromValue(500)))

	orderDraft := Order()

	orderDraft.AddVariableAmount(aggBalance1)
	
	transaction1 := Transaction().
			SetSource(source1).
			SetDestination(destination1).
			SetAsset("COIN").
		SetAmount(MicroAmount.FromValue(25))

	transaction2 := Transaction(). 
			SetSource(source1). 
			SetDestination(destination2). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromValue(75))
	
	orderDraft.AddTransactionWithCondition(transaction1, condTrans1)
	orderDraft.AddTransactionWithCondition(transaction2, condTrans1)

	// OrderDraft.Commit()

}

func CondExample2(){
	
    /*
		
		FICTIVE EXAMPLE BASED ON FUTUR PROBABLE FEATURES
		Conditionnal Order

		{

			globalCondition = {
				$gt : { AggBalance[@bank], 2000 }	
			}

			
			send [COIN Percent[Balance[@bank:creditagricole], 20]] (
				source = @bank:creditagricole
				destination = @bercy
			)

			send [COIN Percent[Balance[@bank:societegeneral], 20]] (
				source = @bank:societegeneral
				destination = @bercy
			)

			send [COIN Percent[Balance[@bank:creditmutuel], 20]] (
				source = @bank:creditmutuel
				destination = @bercy
			)

			send [COIN Percent[Balance[@bank:bankpopulaire], 20]] (
				source = @bank:bankpopulaire
				destination = @bercy
			)

		
		}


	*/

	destination1 := AccountIdentifier("orga_schema", "ledger-01", "@bercy")
	
	source1 := AccountIdentifier("orga_schema", "ledger-01", "@bank:creditagricole")
	source2 := AccountIdentifier("orga_schema", "ledger-01", "@bank:societegeneral")
	source3 := AccountIdentifier("orga_schema", "ledger-01", "@bank:creditmutuel")
	source4 := AccountIdentifier("orga_schema", "ledger-01", "@bank:bankpopulaire")

	aggIdentifier := AccountIdentifier("orga_schema", "ledger-01","@bank")


	aggBalance1 := VariableAmount.New("aggBalance1",
				ValueExpression.New(Operators.AggBalance, aggIdentifier))

	_balance1 := VariableAmount.New("balance1",
				ValueExpression.New(Operators.Balance, source1))

	_balance2 := VariableAmount.New("balance2",
				ValueExpression.New(Operators.Balance, source2))

	_balance3 := VariableAmount.New("balance3",
				ValueExpression.New(Operators.Balance, source3))

	_balance4 := VariableAmount.New("balance4",
				ValueExpression.New(Operators.Balance, source4))

	varAmount1 := VariableAmount.New("varAmount1",
					ValueExpression.New(Operators.Per, _balance1, 20))

	varAmount2 := VariableAmount.New("varAmount2",
					ValueExpression.New(Operators.Per, _balance2, 20))

	varAmount3 := VariableAmount.New("varAmount3",
					ValueExpression.New(Operators.Per, _balance3, 20))

	varAmount4 := VariableAmount.New("varAmount4",
					ValueExpression.New(Operators.Per, _balance4, 20))


	condOrder1 := ConditionOrder.With(LogicalExpression.New(Operators.Gt, aggBalance1, MicroAmount.FromValue(2000)))

	

	transaction1 := Transaction().
			SetSource(source1).
			SetDestination(destination1).
			SetAsset("COIN").
			SetAmount(MicroAmount.FromVariable(varAmount1))

	transaction2 := Transaction(). 
			SetSource(source2). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromVariable(varAmount2))

	transaction3 := Transaction(). 
			SetSource(source3). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromVariable(varAmount3))	
	
	transaction4 := Transaction(). 
			SetSource(source4). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromVariable(varAmount4))


	orderDraft := Order()
	orderDraft.AddConditionOrder(condOrder1)
	
	orderDraft.AddTransaction(transaction1)
	orderDraft.AddTransaction(transaction2)
	orderDraft.AddTransaction(transaction3)
	orderDraft.AddTransaction(transaction4)

	// OrderDraft.Commit()

}