package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_ledger"

type LedgerNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns LedgerNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns LedgerNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns LedgerNodeStore) OrganizationId() string {
	return cns.organizationId
}

func LedgerNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *LedgerNodeStore {
	return &LedgerNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *LedgerNodeStore {
	return fctl.GetStore(ctx, key).(*LedgerNodeStore)
}

func ContextWithStore(ctx context.Context, store *LedgerNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
