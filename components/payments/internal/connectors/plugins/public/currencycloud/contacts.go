package currencycloud

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) fetchContactID(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextOthersResponse{}, errors.New("missing from payload when fetching contacts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	contact, err := p.client.GetContactID(ctx, from.Reference)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	raw, err := json.Marshal(contact)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	return models.FetchNextOthersResponse{
		Others: []models.PSPOther{
			{
				ID:    contact.ID,
				Other: raw,
			},
		},
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
