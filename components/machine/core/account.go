package core

import (
	"fmt"
	"regexp"
)

const accountPattern = "^[a-zA-Z_]+[a-zA-Z0-9_:]*$"

var accountRegexp = regexp.MustCompile(accountPattern)

type Account string

func (Account) GetType() Type { return TypeAccount }
func (a Account) String() string {
	return fmt.Sprintf("@%v", string(a))
}

func ParseAccount(acc Account) error {
	// TODO: handle properly in ledger v1.10
	if acc == "" {
		return nil
	}
	if !accountRegexp.MatchString(string(acc)) {
		return fmt.Errorf("accounts should respect pattern %s", accountPattern)
	}
	return nil
}
