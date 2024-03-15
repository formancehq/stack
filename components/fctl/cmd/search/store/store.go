package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_search"

type SearchNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns SearchNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns SearchNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns SearchNodeStore) OrganizationId() string {
	return cns.organizationId
}

func SearchNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *SearchNodeStore {
	return &SearchNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *SearchNodeStore {
	return fctl.GetStore(ctx, key).(*SearchNodeStore)
}

func ContextWithStore(ctx context.Context, store *SearchNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
