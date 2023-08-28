package fctl

import (
	"flag"
	"io"
	"os"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/pkg/errors"
)

func ReadFile(flags *flag.FlagSet, stack *membershipclient.Stack, where string) (string, error) {
	var ret string
	if where == "-" {
		if NeedConfirm(flags, stack) {
			return "", errors.New("You need to use --confirm flag to use stdin")
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil && err != io.EOF {
			return "", errors.Wrapf(err, "reading stdin")
		}

		ret = string(data)
	} else {
		data, err := os.ReadFile(where)
		if err != nil {
			return "", errors.Wrapf(err, "reading file %s", where)
		}
		ret = string(data)
	}
	return ret, nil
}
