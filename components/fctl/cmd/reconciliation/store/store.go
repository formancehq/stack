package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_reconciliation"

type ReconciliationNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns ReconciliationNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns ReconciliationNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns ReconciliationNodeStore) OrganizationId() string {
	return cns.organizationId
}

func ReconciliationNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *ReconciliationNodeStore {
	return &ReconciliationNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *ReconciliationNodeStore {
	return fctl.GetStore(ctx, key).(*ReconciliationNodeStore)
}

func ContextWithStore(ctx context.Context, store *ReconciliationNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
