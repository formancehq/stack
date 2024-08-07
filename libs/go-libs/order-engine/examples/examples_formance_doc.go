package examples

import (
	. "ledgerorder"
)


func DocExample1(){
	
    /*
		
		BASIC TRANSACTION SPLITTING 
		https://docs.formance.com/ledger/numscript/multi-destination/



		send [COIN 75] (
  			source = @centralbank
 			destination = {
    			50% to @player:leslieknope
    			remaining to @player:annperkins
  			}
		)


	*/

	account_1 := AccountIdentifier("orga_schema", "ledger-01", "@centralbank")
	account_2 := AccountIdentifier("orga_schema", "ledger-01","@player:leslieknope")
	account_3 := AccountIdentifier("orga_schema", "ledger-01", "@player:annperkins")

	orderDraft := Order()
	
	transaction1 := Transaction().
			SetSource(account_1).
			SetDestination(account_2).
			SetAsset("COIN").
		SetAmount(MicroAmount.FromValue(38))

	transaction2 := Transaction(). 
			SetSource(account_1). 
			SetDestination(account_3). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromValue(37))
	
	orderDraft.AddTransaction(transaction1)
	orderDraft.AddTransaction(transaction2)



}

func DocExample2(){
	
    /*
		
		NESTED TRANSACTION SPLITTING
		https://docs.formance.com/ledger/numscript/multi-destination/



		send [COIN 75] (
			source = @centralbank
			destination = {
				15% to @salestax
				43% to @player:leslieknope
				remaining to @player:annperkins
			}
		)


	*/

	account_1 := AccountIdentifier("orga_schema", "ledger-01", "@centralbank")
	account_2 := AccountIdentifier("orga_schema", "ledger-01", "@salestax")
	account_3 := AccountIdentifier("orga_schema", "ledger-01","@player:leslieknope")
	account_4 := AccountIdentifier("orga_schema", "ledger-01", "@player:annperkins")

	OrderDraft := Order()
	
	transaction1 := Transaction().
			SetSource(account_1).
			SetDestination(account_2).
			SetAsset("COIN").
		SetAmount(MicroAmount.FromValue(12))

	transaction2 := Transaction(). 
			SetSource(account_1). 
			SetDestination(account_3). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromValue(32))
	
	transaction3 := Transaction(). 
		SetSource(account_1). 
		SetDestination(account_4). 
		SetAsset("COIN"). 
		SetAmount(MicroAmount.FromValue(31))

	OrderDraft.AddTransaction(transaction1)
	OrderDraft.AddTransaction(transaction2)
	OrderDraft.AddTransaction(transaction3)

	// OrderDraft.Commit()

}

func DocExample3(){
	
    /*
		
		Specifying backup accounts
		https://docs.formance.com/ledger/numscript/multi-source/


		send [COIN 100] (
		source = {
			@player:andydwyer
			@player:andydwyer:chest
		}
		destination = @player:aprilludgate
		)


	*/

	source1 := AccountIdentifier("orga_schema", "ledger-01", "@player:andydwyer")
	source2 := AccountIdentifier("orga_schema", "ledger-01", "@player:andydwyer:chest")
	destination1 := AccountIdentifier("orga_schema", "ledger-01","@player:aprilludgate")

	source1Balance := VariableAmount.New("source1Balance", 
							ValueExpression.New(Operators.Balance, source1))

	amountSource1 := VariableAmount.New("amountSource1", 
                ValueExpression.New(Operators.Min, source1Balance, MicroAmount.FromValue(100)))

	amountSource2 := VariableAmount.New("amountSource2", 
				ValueExpression.New(Operators.Sub, MicroAmount.FromValue(100), amountSource1))

	
	
	transaction1 := Transaction().
			SetSource(source1).
			SetDestination(destination1).
			SetAsset("COIN").
		SetAmount(MicroAmount.FromVariable(amountSource1))

	transaction2 := Transaction(). 
			SetSource(source2). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromVariable(amountSource2))
	
	
	OrderDraft := Order()

	OrderDraft.AddVariableAmount(amountSource1)
	OrderDraft.AddVariableAmount(amountSource2)

	OrderDraft.AddTransaction(transaction1)
	OrderDraft.AddTransaction(transaction2)
	

	// OrderDraft.Commit()

}

func DocExample4(){
	
    /*
		
		Specifying backup accounts
		https://docs.formance.com/ledger/numscript/multi-source/


		send [COIN 75] (
			source = {
				50% from {
				max [COIN 10] from @player:donnameagle
				@player:donnameagle:chest
				}
				remaining from @player:tomhaverford
			}
			destination = @centralbank
		)


	*/

	source1 := AccountIdentifier("orga_schema", "ledger-01", "@player:donnameagle")
	source2 := AccountIdentifier("orga_schema", "ledger-01", "@player:donnameagle:chest")
	source3 := AccountIdentifier("orga_schema", "ledger-01", "@player:tomhaverford")
	
	destination1 := AccountIdentifier("orga_schema", "ledger-01","@centralbank")

	source1Balance := VariableAmount.New("source1Balance", 
							ValueExpression.New(Operators.Balance, source1))

	amountSource1 := VariableAmount.New("amountSource1", 
                ValueExpression.New(Operators.Min, source1Balance, MicroAmount.FromValue(38)))

	amountSource2 := VariableAmount.New("amountSource2", 
				ValueExpression.New(Operators.Sub, MicroAmount.FromValue(38), amountSource1))

	
	
	transaction1 := Transaction().
			SetSource(source1).
			SetDestination(destination1).
			SetAsset("COIN").
		SetAmount(MicroAmount.FromVariable(amountSource1))

	transaction2 := Transaction(). 
			SetSource(source2). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromVariable(amountSource2))

	transaction3 := Transaction(). 
			SetSource(source3). 
			SetDestination(destination1). 
			SetAsset("COIN"). 
			SetAmount(MicroAmount.FromValue(37))
	
	
	OrderDraft := Order()

	OrderDraft.AddVariableAmount(amountSource1)
	OrderDraft.AddVariableAmount(amountSource2)
	
	OrderDraft.AddTransaction(transaction1)
	OrderDraft.AddTransaction(transaction2)
	OrderDraft.AddTransaction(transaction3)
	

	// OrderDraft.Commit()

}