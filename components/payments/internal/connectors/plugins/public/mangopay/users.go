package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type usersState struct {
	LastPage         int       `json:"lastPage"`
	LastCreationDate time.Time `json:"lastCreationDate"`
}

func (p Plugin) fetchNextUsers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	var oldState usersState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextOthersResponse{}, err
		}
	} else {
		oldState = usersState{
			LastPage: 1,
		}
	}

	newState := usersState{
		LastPage:         oldState.LastPage,
		LastCreationDate: oldState.LastCreationDate,
	}

	var users []models.PSPOther
	hasMore := false
	page := oldState.LastPage
	for {
		pagedUsers, err := p.client.GetUsers(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextOthersResponse{}, err
		}

		if len(pagedUsers) == 0 {
			break
		}

		for _, user := range pagedUsers {
			userCreationDate := time.Unix(user.CreationDate, 0)
			switch userCreationDate.Compare(oldState.LastCreationDate) {
			case -1, 0:
				// creationDate <= state.LastCreationDate, nothing to do,
				// we already processed this user.
				continue
			default:
			}

			raw, err := json.Marshal(user)
			if err != nil {
				return models.FetchNextOthersResponse{}, err
			}

			users = append(users, models.PSPOther{
				ID:    user.ID,
				Other: raw,
			})

			newState.LastCreationDate = userCreationDate

			if len(users) >= req.PageSize {
				break
			}
		}

		if len(users) < req.PageSize {
			break
		}

		if len(users) >= req.PageSize {
			hasMore = true
			break
		}

		page++
	}

	newState.LastPage = page

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	return models.FetchNextOthersResponse{
		Others:   users,
		NewState: payload,
		HasMore:  hasMore,
	}, nil
}
