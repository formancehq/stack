package wise

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/formancehq/payments/internal/models"
)

type profilesState struct {
	// Profiles are ordered by their ID
	LastProfileID uint64 `json:"lastProfileID"`
}

func (p Plugin) fetchNextProfiles(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	var oldState profilesState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextOthersResponse{}, err
		}
	}

	newState := profilesState{
		LastProfileID: oldState.LastProfileID,
	}

	var others []models.PSPOther
	hasMore := false
	profiles, err := p.client.GetProfiles(ctx)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	for _, profile := range profiles {
		if profile.ID <= oldState.LastProfileID {
			continue
		}

		raw, err := json.Marshal(profile)
		if err != nil {
			return models.FetchNextOthersResponse{}, err
		}

		others = append(others, models.PSPOther{
			ID:    strconv.FormatUint(profile.ID, 10),
			Other: raw,
		})

		newState.LastProfileID = profile.ID

		if len(others) >= req.PageSize {
			hasMore = true
			break
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	return models.FetchNextOthersResponse{
		Others:   others,
		NewState: payload,
		HasMore:  hasMore,
	}, nil
}
