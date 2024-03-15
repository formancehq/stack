package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_payments"

type PaymentsNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns PaymentsNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns PaymentsNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns PaymentsNodeStore) OrganizationId() string {
	return cns.organizationId
}

func PaymentsNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *PaymentsNodeStore {
	return &PaymentsNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *PaymentsNodeStore {
	return fctl.GetStore(ctx, key).(*PaymentsNodeStore)
}

func ContextWithStore(ctx context.Context, store *PaymentsNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
