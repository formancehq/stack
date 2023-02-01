package fctl

import (
	"github.com/pterm/pterm"
)

var (
	BasicText      = pterm.DefaultBasicText
	BasicTextGreen = pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightGreen))
	BasicTextRed   = pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightRed))
	BasicTextCyan  = pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightCyan))
	Section        = pterm.SectionPrinter{
		Style:           &pterm.ThemeDefault.SectionStyle,
		Level:           1,
		TopPadding:      0,
		BottomPadding:   0,
		IndentCharacter: "#",
	}
)
