package fctl

import (
	"io"

	"github.com/pterm/pterm"
)

func Println(args ...any) {
	pterm.Println(args...)
}

func Printfln(fmt string, args ...any) {
	pterm.Printfln(fmt, args...)
}

func Highlightln(out io.Writer, format string, args ...any) {
	BasicTextCyan.WithWriter(out).Printfln(format, args...)
}

func SuccessWriter(out io.Writer) *pterm.PrefixPrinter {
	return pterm.Success.WithWriter(out)
}

func Success(out io.Writer, format string, args ...any) {
	SuccessWriter(out).Printfln(format, args...)
}

func ErrorWriter(out io.Writer) *pterm.PrefixPrinter {
	return pterm.Error.WithWriter(out)
}

func Error(out io.Writer, format string, args ...any) {
	ErrorWriter(out).Printfln(format, args...)
}
