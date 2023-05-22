package internal

import (
	"fmt"
	"io"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/iancoleman/strcase"
	"github.com/pterm/pterm"
)

func PrintStackInformation(out io.Writer, profile *fctl.Profile, stack *membershipclient.Stack, versions *formance.GetVersionsResponse) error {
	baseUrlStr := profile.ServicesBaseUrl(stack).String()

	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), stack.Id, ""})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), stack.Name, ""})
	tableData = append(tableData, []string{pterm.LightCyan("Region"), stack.RegionID, ""})
	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Println()
	fctl.Section.WithWriter(out).Println("Versions")
	tableData = pterm.TableData{}
	for _, service := range versions.Versions {
		tableData = append(tableData, []string{pterm.LightCyan(strcase.ToCamel(service.Name)), service.Version,
			fmt.Sprintf("%s/api/%s", baseUrlStr, service.Name)})
	}
	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Println()
	fctl.Section.WithWriter(out).Println("Metadata")
	tableData = pterm.TableData{}
	for k, v := range stack.Metadata {
		tableData = append(tableData, []string{pterm.LightCyan(k), v})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}
