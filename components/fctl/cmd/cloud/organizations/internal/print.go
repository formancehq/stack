package internal

import (
	"strings"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/pterm/pterm"
)

func PrintOrganization(organization *membershipclient.Organization) error {
	pterm.DefaultSection.Println("Organization")

	data := pterm.TableData{
		{"ID", organization.Id},
		{"Name", organization.Name},
		{"Default Stack Role", func() string {
			if len(organization.DefaultStackAccess) == 0 {
				return "None"
			}
			return strings.Join(organization.DefaultStackAccess, ", ")
		}()},
		{"Default Organization Role", func() string {
			if len(organization.DefaultOrganizationAccess) == 0 {
				return "None"
			}
			return strings.Join(organization.DefaultOrganizationAccess, ", ")
		}()},
		{"Domain", func() string {
			if organization.Domain == nil {
				return "None"
			}
			return *organization.Domain
		}()},
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
