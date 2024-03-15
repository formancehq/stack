package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
)

const key = "_stack"

type StackNodeStore struct {
	Config           *fctl.Config
	MembershipClient *fctl.MembershipClient
	organizationId   string
}

func (cns StackNodeStore) Client() *membershipclient.DefaultApiService {
	return cns.MembershipClient.DefaultApi
}

func (cns StackNodeStore) OrganizationId() string {
	return cns.organizationId
}

func StackNode(config *fctl.Config, apiClient *fctl.MembershipClient, organization string) *StackNodeStore {
	return &StackNodeStore{
		Config:           config,
		MembershipClient: apiClient,
		organizationId:   organization,
	}
}

func GetStore(ctx context.Context) *StackNodeStore {
	return fctl.GetStore(ctx, key).(*StackNodeStore)
}

func ContextWithStore(ctx context.Context, store *StackNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
