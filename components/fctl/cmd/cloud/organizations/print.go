package organizations

import (
	"strconv"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/pterm/pterm"
)

func PrintOrganization(org *membershipclient.Organization) error {
	pterm.DefaultSection.Println("Organization")

	data := [][]string{
		{"ID", org.Id},
		{"Name", org.Name},
		{"Domain", func() string {
			if org.Domain == nil {
				return "None"
			}
			return *org.Domain
		}()},
		{"Default Stack Role", func() string {
			return string(*org.DefaultStackAccess)
		}()},
		{"Default Organization Role", func() string {
			return string(*org.DefaultOrganizationAccess)
		}()},
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

func PrintOrganizationFromStore(store *DescribeStore) error {
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
			return string(*store.DefaultStackAccess)
		}()},
		{"Default Organization Role", func() string {
			return string(*store.DefaultOrganizationAccess)
		}()},
	}

	if store.expanded {
		data = append(data,
			[]string{"Owner Email", store.Owner.Email},
			[]string{"Total Users", func() string {
				return strconv.Itoa(int(store.TotalUsers))
			}()},
			[]string{"Total Stacks", func() string {
				return strconv.Itoa(int(store.TotalStacks))
			}()},
		)
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
