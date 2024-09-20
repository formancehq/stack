package internal

import (
	"strconv"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/pterm/pterm"
)

func PrintOrganization(store *membershipclient.OrganizationExpanded) error {
	pterm.DefaultSection.Println("Organization")

	data := [][]string{
		{"ID", store.Id},
		{"Name", store.Name},
		{"Domain", func() string {
			if store.Domain == nil {
				return "None"
			}
			return *store.Domain
		}()},
		{"Default Stack Role", func() string {
			if store.DefaultStackAccess == nil {
				return "None"
			}
			return string(*store.DefaultStackAccess)
		}()},
		{"Default Organization Role", func() string {
			if store.DefaultOrganizationAccess == nil {
				return "None"
			}
			return string(*store.DefaultOrganizationAccess)
		}()},
	}

	if store.Owner != nil {
		data = append(data, []string{"Owner", store.Owner.Email})
	}

	if store.TotalUsers != nil {
		data = append(data, []string{"Total Users", strconv.Itoa(int(*store.TotalUsers))})
	}

	if store.TotalStacks != nil {
		data = append(data, []string{"Total Stacks", strconv.Itoa(int(*store.TotalStacks))})
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
