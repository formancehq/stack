package main

import (
	"fmt"

	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/machine/script/compiler"
	"github.com/numary/ledger/pkg/machine/vm"
)

/*
15% from {
	@alice
	@bob
}
remaining from @bob
*/

func main() {
	program, err := compiler.Compile(`
		// This is a comment
		vars {
		asset $asset
		monetary $mon1 = balance(@alice, $asset)
		monetary $mon2 = balance(@bob, $asset)
		}
		send $mon1 + $mon2 - [$asset 1] (
			source = @wallet
			destination = @account_c
		)`)
	if err != nil {
		panic(err)
	}
	fmt.Print(program)

	m := vm.NewMachine(*program)
	m.Debug = true

	if err = m.SetVars(map[string]core.Value{
		"asset": core.Asset("COIN"),
	}); err != nil {
		panic(err)
	}

	initialBalances := map[string]map[string]*core.MonetaryInt{
		"alice": {"COIN": core.NewMonetaryInt(42)},
		"bob":   {"COIN": core.NewMonetaryInt(100)},
		"wallet": {"COIN": core.NewMonetaryInt(1000)},
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
			if req.Asset != "" {
				req.Response <- initialBalances[req.Account][req.Asset]
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
