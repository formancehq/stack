package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_webhooks"

type WebhooksNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns WebhooksNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns WebhooksNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns WebhooksNodeStore) OrganizationId() string {
	return cns.organizationId
}

func WebhooksNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *WebhooksNodeStore {
	return &WebhooksNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *WebhooksNodeStore {
	return fctl.GetStore(ctx, key).(*WebhooksNodeStore)
}

func ContextWithStore(ctx context.Context, store *WebhooksNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
