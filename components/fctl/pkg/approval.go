package fctl

import (
	"flag"
	"fmt"
	"strings"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
)

var ErrMissingApproval = errors.New("Missing approval.")

var interactiveContinue = pterm.InteractiveContinuePrinter{
	DefaultValueIndex: 0,
	DefaultText:       "Do you want to continue",
	TextStyle:         &pterm.ThemeDefault.PrimaryStyle,
	Options:           []string{"y", "n"},
	OptionsStyle:      &pterm.ThemeDefault.SuccessMessageStyle,
	SuffixStyle:       &pterm.ThemeDefault.SecondaryStyle,
}

const (
	ProtectedStackMetadata = "github.com/formancehq/fctl/protected"
	confirmFlag            = "confirm"
)

func IsProtectedStack(stack *membershipclient.Stack) bool {
	return stack.Metadata != nil && (stack.Metadata)[ProtectedStackMetadata] == "Yes"
}

func NeedConfirm(flags *flag.FlagSet, stack *membershipclient.Stack) bool {
	if !IsProtectedStack(stack) {
		return false
	}
	if GetBool(flags, confirmFlag) {
		return false
	}
	return true
}

func CheckStackApprobation(flags *flag.FlagSet, stack *membershipclient.Stack, disclaimer string, args ...any) bool {
	if !IsProtectedStack(stack) {
		return true
	}
	if GetBool(flags, confirmFlag) {
		return true
	}

	disclaimer = fmt.Sprintf(disclaimer, args...)

	result, err := interactiveContinue.WithDefaultText(disclaimer + ".\r\n" + pterm.DefaultInteractiveContinue.DefaultText).Show()
	if err != nil {
		panic(err)
	}
	return strings.ToLower(result) == "y"
}

func CheckOrganizationApprobation(flags *flag.FlagSet, disclaimer string, args ...any) bool {
	if GetBool(flags, confirmFlag) {
		return true
	}

	result, err := interactiveContinue.WithDefaultText(disclaimer + ".\r\n" + pterm.DefaultInteractiveContinue.DefaultText).Show()
	if err != nil {
		panic(err)
	}
	return strings.ToLower(result) == "y"
}
