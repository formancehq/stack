package internal

import (
	"github.com/formancehq/fctl/membershipclient"
	"github.com/pterm/pterm"
)

func PrintOrganization(organization *membershipclient.Organization) error {
	pterm.DefaultSection.Println("Organization")

	data := pterm.TableData{
		{"ID", organization.Id},
		{"Name", organization.Name},
		{"Default Stack Role", func() string {
			return string(*organization.DefaultStackAccess)
		}()},
		{"Default Organization Role", func() string {
			return string(*organization.DefaultOrganizationAccess)
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
