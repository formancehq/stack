package main

import (
	"fmt"

	"github.com/numary/ledger/pkg/machine"
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
	fmt.Print(program)

	m := vm.NewMachine(*program)
	m.Debug = true

	if err = m.SetVars(map[string]machine.Value{
		"dest": machine.Account("charlie"),
	}); err != nil {
		panic(err)
	}

	initialBalances := map[string]map[string]*machine.MonetaryInt{
		"alice": {"COIN": machine.NewMonetaryInt(10)},
		"bob":   {"COIN": machine.NewMonetaryInt(100)},
	}

	{
		ch, err := m.ResolveResources()
		if err != nil {
			panic(err)
		}
		for req := range ch {
			if req.Error != nil {
				panic(req.Error)
			}
		}
	}

	{
		ch, err := m.ResolveBalances()
		if err != nil {
			panic(err)
		}
		for req := range ch {
			val := initialBalances[req.Account][req.Asset]
			if req.Error != nil {
				panic(req.Error)
			}
			req.Response <- val
		}
	}

	exitCode, err := m.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println("Exit code:", exitCode)
	fmt.Println(m.Postings)
	fmt.Println(m.TxMeta)
}
