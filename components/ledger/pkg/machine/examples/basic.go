package main

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/ledger/pkg/machine/internal"
	"github.com/formancehq/ledger/pkg/machine/script/compiler"
	"github.com/formancehq/ledger/pkg/machine/vm"
)

func main() {
	program, err := compiler.Compile(`
		// This is a comment
		vars {
			account $dest
		}
		send [COIN 99] (
			source = {
				15% from {
					@alice
					@bob
				}
				remaining from @bob
			}
			destination = $dest
		)`)
	if err != nil {
		panic(err)
	}

	s, _ := json.MarshalIndent(program, "", "\t")
	fmt.Println(string(s))

	ledger := vm.MockupLedger(map[string]vm.MockupAccount{
		"alice": {
			Balances: map[string]internal.MonetaryInt{
				"COIN": *internal.NewMonetaryInt(10),
			},
			Meta: map[string]internal.Value{},
		},
		"bob": {
			Balances: map[string]internal.MonetaryInt{
				"COIN": *internal.NewMonetaryInt(100),
			},
			Meta: map[string]internal.Value{},
		},
	})

	m := vm.NewMachine(ledger, map[string]internal.Value{
		"dest": internal.AccountAddress("charlie"),
	})

	err = m.Execute(*program)
	if err != nil {
		panic(err)
	}

	fmt.Println("Postings:")
	for _, posting := range m.Postings {
		fmt.Printf("[%v %v] %v -> %v\n", posting.Asset, posting.Amount, posting.Source, posting.Destination)
	}
	fmt.Println("Tx Meta:")
	for key, value := range m.TxMeta {
		fmt.Printf("%v: %v", key, value)
	}
}
