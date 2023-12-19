package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/get-momo/atlar-v1-go-client/client/credit_transfers"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

func (c *Client) PostV1CreditTransfers(ctx context.Context, req *atlar_models.CreatePaymentRequest) (*credit_transfers.PostV1CreditTransfersCreated, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "create_credit_transfer")
	now := time.Now()
	defer f(ctx, now)

	postCreditTransfersParams := credit_transfers.PostV1CreditTransfersParams{
		Context:        ctx,
		CreditTransfer: req,
	}

	return c.client.CreditTransfers.PostV1CreditTransfers(&postCreditTransfersParams)

}

func (c *Client) GetV1CreditTransfersGetByExternalIDExternalID(ctx context.Context, externalID string) (*credit_transfers.GetV1CreditTransfersGetByExternalIDExternalIDOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "get_credit_transfer")
	now := time.Now()
	defer f(ctx, now)

	getCreditTransferParams := credit_transfers.GetV1CreditTransfersGetByExternalIDExternalIDParams{
		Context:    ctx,
		ExternalID: externalID,
	}

	return c.client.CreditTransfers.GetV1CreditTransfersGetByExternalIDExternalID(&getCreditTransferParams)
}
