package internal

import (
	"fmt"
	"io"
	"net/url"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/iancoleman/strcase"
	"github.com/pterm/pterm"
)

func PrintStackInformation(out io.Writer, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) error {
	baseUrlStr := profile.ServicesBaseUrl(stack)

	err := printInformation(out, stack)

	if err != nil {
		return err
	}

	err = printVersion(out, baseUrlStr, versions, stack)

	if err != nil {
		return err
	}

	err = printMetadata(out, stack)
	if err != nil {
		return err
	}

	return nil
}

func printInformation(out io.Writer, stack *membershipclient.Stack) error {

	fctl.Section.WithWriter(out).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), stack.Id, ""})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), stack.Name, ""})
	tableData = append(tableData, []string{pterm.LightCyan("Region"), stack.RegionID, ""})
	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}

func printVersion(out io.Writer, url *url.URL, versions *shared.GetVersionsResponse, stack *membershipclient.Stack) error {
	fmt.Fprintln(out)
	fctl.Section.WithWriter(out).Println("Versions")

	tableData := pterm.TableData{}

	for _, service := range versions.Versions {
		tableData = append(tableData, []string{pterm.LightCyan(strcase.ToCamel(service.Name)), service.Version,
			fmt.Sprintf("%s/api/%s", url.String(), service.Name)})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}

func printMetadata(out io.Writer, stack *membershipclient.Stack) error {

	fctl.Section.WithWriter(out).Println("Metadata")

	tableData := pterm.TableData{}

	for k, v := range stack.Metadata {
		tableData = append(tableData, []string{pterm.LightCyan(k), v})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}
