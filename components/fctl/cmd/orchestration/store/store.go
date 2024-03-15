package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

const key = "_orchestration"

type OrchestrationNodeStore struct {
	Config         *fctl.Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns OrchestrationNodeStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns OrchestrationNodeStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns OrchestrationNodeStore) OrganizationId() string {
	return cns.organizationId
}

func OrchestrationNode(config *fctl.Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *OrchestrationNodeStore {
	return &OrchestrationNodeStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStore(ctx context.Context) *OrchestrationNodeStore {
	return fctl.GetStore(ctx, key).(*OrchestrationNodeStore)
}

func ContextWithStore(ctx context.Context, store *OrchestrationNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
