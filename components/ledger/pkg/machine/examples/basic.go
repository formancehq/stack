package main

import (
	"encoding/json"
	"fmt"

	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/machine/script/compiler"
	"github.com/numary/ledger/pkg/machine/vm"
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
			Balances: map[string]core.MonetaryInt{
				"COIN": *core.NewMonetaryInt(10),
			},
			Meta: map[string]core.Value{},
		},
		"bob": {
			Balances: map[string]core.MonetaryInt{
				"COIN": *core.NewMonetaryInt(100),
			},
			Meta: map[string]core.Value{},
		},
	})

	m := vm.NewMachine(ledger, map[string]core.Value{
		"dest": core.AccountAddress("charlie"),
	})

	err = m.Run(*program)
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
