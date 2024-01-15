package organizations

import (
	"github.com/pterm/pterm"
)

func PrintOrganization(store *DescribeStore) error {
	pterm.DefaultSection.Println("Organization")

	data := [][]string{
		{"ID", store.Organization.Id},
		{"Name", store.Organization.Name},
		{"Domain", func() string {
			if store.Organization.Domain == nil {
				return "None"
			}
			return *store.Organization.Domain
		}()},
		{"Default Stack Role", func() string {
			return string(*store.Organization.DefaultStackAccess)
		}()},
		{"Default Organization Role", func() string {
			return string(*store.Organization.DefaultOrganizationAccess)
		}()},
	}

	if store.OrganizationExpandedAllOf != nil {
		data = append(data,
			[]string{"Owner ID", store.OrganizationExpandedAllOf.Owner.Id},
			[]string{"Owner Email", store.OrganizationExpandedAllOf.Owner.Email},
			[]string{"Total Users", func() string { return string(*store.OrganizationExpandedAllOf.TotalUsers) }()},
			[]string{"Total Stacks", func() string { return string(*store.OrganizationExpandedAllOf.TotalStacks) }()},
		)
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
