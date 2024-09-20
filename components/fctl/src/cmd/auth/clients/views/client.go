package views

import (
	"fmt"
	"io"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
)

func PrintClient(out io.Writer, client *shared.Client) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), client.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), client.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), fctl.StringPointerToString(client.Description)})
	tableData = append(tableData, []string{pterm.LightCyan("Public"), fctl.BoolPointerToString(client.Public)})
	if len(client.Scopes) > 0 {
		tableData = append(tableData, []string{pterm.LightCyan("Scopes"), strings.Join(client.Scopes, " ")})
	}

	fctl.Section.WithWriter(out).Println("Information :")
	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	fmt.Fprintln(out, "")

	return nil
}
