package v1beta1

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func init() {
	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
}

type Client struct {
	rest.Interface
}

func NewClient(restClient rest.Interface) *Client {
	return &Client{
		Interface: restClient,
	}
}

func (c *Client) Stacks() StackInterface {
	return &stackClient{
		restClient: c.Interface,
	}
}

func (c *Client) Auths() AuthInterface {
	return &AuthClient{
		restClient: c.Interface,
	}
}

func (c *Client) Gateways() GatewayInterface {
	return &gatewayClient{
		restClient: c.Interface,
	}
}

func (c *Client) Ledgers() LedgerInterface {
	return &LedgerClient{
		restClient: c.Interface,
	}
}

func (c *Client) Orchestrations() OrchestrationInterface {
	return &OrchestrationClient{
		restClient: c.Interface,
	}
}

func (c *Client) Payments() PaymentsInterface {
	return &paymentsClient{
		restClient: c.Interface,
	}
}

func (c *Client) Reconciliations() ReconciliationInterface {
	return &reconciliationClient{
		restClient: c.Interface,
	}
}

func (c *Client) Searches() SearchInterface {
	return &SearchClient{
		restClient: c.Interface,
	}
}

func (c *Client) Wallets() WalletsInterface {
	return &walletsClient{
		restClient: c.Interface,
	}
}

func (c *Client) Webhooks() WebhooksInterface {
	return &webhooksClient{
		restClient: c.Interface,
	}
}
