package store

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
)

const key = "_cloud"

type CloudNodeStore struct {
	Config           *fctl.Config
	MembershipClient *fctl.MembershipClient
}

func (cns CloudNodeStore) Client() *membershipclient.DefaultApiService {
	return cns.MembershipClient.DefaultApi
}

func CloudNode(config *fctl.Config, apiClient *fctl.MembershipClient) *CloudNodeStore {
	return &CloudNodeStore{
		Config:           config,
		MembershipClient: apiClient,
	}
}

func GetStore(ctx context.Context) *CloudNodeStore {
	return fctl.GetStore(ctx, key).(*CloudNodeStore)
}

func ContextWithStore(ctx context.Context, store *CloudNodeStore) context.Context {
	return fctl.ContextWithStore(ctx, key, store)
}
