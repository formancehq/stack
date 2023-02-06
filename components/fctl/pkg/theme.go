package fctl

import (
	"github.com/pterm/pterm"
)

var (
	StyleGreen = pterm.NewStyle(pterm.FgLightGreen)
	StyleRed   = pterm.NewStyle(pterm.FgLightRed)
	StyleCyan  = pterm.NewStyle(pterm.FgLightCyan)

	BasicText      = pterm.DefaultBasicText
	BasicTextGreen = pterm.DefaultBasicText.WithStyle(StyleGreen)
	BasicTextRed   = pterm.DefaultBasicText.WithStyle(StyleRed)
	BasicTextCyan  = pterm.DefaultBasicText.WithStyle(StyleCyan)
	Section        = pterm.SectionPrinter{
		Style:           &pterm.ThemeDefault.SectionStyle,
		Level:           1,
		TopPadding:      0,
		BottomPadding:   0,
		IndentCharacter: "#",
	}
)
