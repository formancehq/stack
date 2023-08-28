package fctl

import (
	"errors"
	"flag"
	"os"
	"sync"
)

type fValue[T any] struct {
	value T
}

func (f *fValue[T]) String() T {
	return f.value
}

func (f *fValue[T]) Set(s T) error {
	f.value = s
	return nil
}

func (f *fValue[T]) Get() *T {
	return &f.value
}

var scopeFlags *flag.FlagSet
var lock = &sync.Mutex{}

var (
	stackFlagV        = &fValue[string]{value: ""}
	organizationFlagV = &fValue[string]{value: ""}
	ledgerFlagV       = &fValue[string]{value: ""}
	Stack             = getScopeFlags(stackFlag)
	Ledger            = getScopeFlags("ledger")
	Organization      = getScopeFlags(organizationFlag)
)

func getScopesFlagsInstance() *flag.FlagSet {
	if scopeFlags == nil {
		lock.Lock()
		defer lock.Unlock()
		if scopeFlags == nil {
			scopeFlags = flag.NewFlagSet("scopes", flag.ContinueOnError)
			scopeFlags.StringVar(stackFlagV.Get(), "stack", "", "Specific stack id (not required if only one stack is present)")
			scopeFlags.StringVar(organizationFlagV.Get(), "organization", "", "Selected organization (not required if only one organization is present)")
			scopeFlags.StringVar(ledgerFlagV.Get(), "ledger", "default", "Specific ledger name")
		}
	}

	return scopeFlags
}

func getScopeFlags(name string) *flag.Flag {
	return getScopesFlagsInstance().Lookup(name)
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(errors.New("unable to get home directory"))
	}
	return homedir
}
