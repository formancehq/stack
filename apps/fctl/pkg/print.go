package fctl

import (
	"github.com/pterm/pterm"
)

func Println(args ...any) {
	pterm.Println(args...)
}

func Printfln(fmt string, args ...any) {
	pterm.Printfln(fmt, args...)
}
