package store

import (
	"context"
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

const key = "_stack"

type StackNodeStore struct {
	*fctl.MembershipStore
	organizationId string
}

func (cns StackNodeStore) Client() *membershipclient.DefaultApiService {
	return cns.MembershipClient.DefaultApi
}

func (cns StackNodeStore) OrganizationId() string {
	return cns.organizationId
}

func StackNode(store *fctl.MembershipStore, organization string) *StackNodeStore {
	return &StackNodeStore{
		MembershipStore: store,
		organizationId:  organization,
	}
}

func GetStore(ctx context.Context) *StackNodeStore {
	return fctl.GetStore(ctx, key).(*StackNodeStore)
}

func ContextWithStore(ctx context.Context, store *StackNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}

func NewMembershipStackStore(cmd *cobra.Command) error {
	err := fctl.NewMembershipStore(cmd)
	if err != nil {
		return err
	}

	store := fctl.GetMembershipStore(cmd.Context())
	organization, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return err
	}

	cmd.SetContext(ContextWithStore(cmd.Context(), StackNode(store, organization)))

	return nil
}

func (cns *StackNodeStore) CheckRegionCapability(key string, checker func([]any) bool) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		stack, err := fctl.ResolveStack(cmd, cns.Config, cns.organizationId)
		if err != nil {
			return
		}

		region, _, err := cns.Client().GetRegion(cmd.Context(), cns.organizationId, stack.RegionID).Execute()
		if err != nil {
			return
		}

		capabilities, err := fctl.StructToMap(region.Data.Capabilities)
		if err != nil {
			return
		}

		if value, ok := capabilities[key]; ok {
			if values := value.([]interface{}); len(values) > 0 {
				if !checker(values) {
					return fmt.Errorf("unsupported membership server version: %s", value)
				}

			}
		}
		return
	}
}
