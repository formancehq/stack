package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_wallets"

type WalletsNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns WalletsNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns WalletsNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns WalletsNodeStore) OrganizationId() string {
	return cns.organizationId
}

func WalletsNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *WalletsNodeStore {
	return &WalletsNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *WalletsNodeStore {
	return fctl.GetStore(ctx, key).(*WalletsNodeStore)
}

func ContextWithStore(ctx context.Context, store *WalletsNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
