package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"

	"github.com/formancehq/go-libs/api"
	"github.com/pkg/errors"
)

type listAccountsRepository interface {
	ListAccounts(ctx context.Context, sort storage.Sorter, pagination storage.Paginator) ([]*models.Account, error)
}

type accountResponse struct {
	ID        string             `json:"id"`
	Reference string             `json:"reference"`
	CreatedAt time.Time          `json:"createdAt"`
	Provider  string             `json:"provider"`
	Type      models.AccountType `json:"type"`
}

func listAccountsHandler(repo listAccountsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var sorter storage.Sorter

		if sortParams := r.URL.Query()["sort"]; sortParams != nil {
			for _, s := range sortParams {
				parts := strings.SplitN(s, ":", 2)

				var order storage.SortOrder

				if len(parts) > 1 {
					//nolint:goconst // allow duplicate string
					switch parts[1] {
					case "asc", "ASC":
						order = storage.SortOrderAsc
					case "dsc", "desc", "DSC", "DESC":
						order = storage.SortOrderDesc
					default:
						handleValidationError(w, r, errors.New("sort order not well specified, got "+parts[1]))

						return
					}
				}

				column := parts[0]

				sorter.Add(column, order)
			}
		}

		skip, err := integerWithDefault(r, "skip", 0)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		limit, err := integerWithDefault(r, "limit", maxPerPage)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		if limit > maxPerPage {
			limit = maxPerPage
		}

		ret, err := repo.ListAccounts(r.Context(), sorter, storage.Paginate(skip, limit))
		if err != nil {
			handleServerError(w, r, err)

			return
		}

		data := make([]*accountResponse, len(ret))

		for i := range ret {
			data[i] = &accountResponse{
				ID:        ret[i].ID.String(),
				Reference: ret[i].Reference,
				CreatedAt: ret[i].CreatedAt,
				Provider:  ret[i].Provider,
				Type:      ret[i].Type,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[[]*accountResponse]{
			Data: &data,
		})
		if err != nil {
			handleServerError(w, r, err)

			return
		}
	}
}
