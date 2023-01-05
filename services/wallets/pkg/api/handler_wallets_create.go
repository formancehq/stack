package api

import (
	"net/http"

	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/render"
)

type CreateWalletRequest struct {
	Metadata metadata.Metadata `json:"metadata"`
	Name     string            `json:"name"`
}

func (c *CreateWalletRequest) Bind(r *http.Request) error {
	return nil
}

func (m *MainHandler) CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	data := &CreateWalletRequest{}
	if r.ContentLength > 0 {
		if err := render.Bind(r, data); err != nil {
			badRequest(w, ErrorCodeValidation, err)
			return
		}
	}

	wallet, err := m.repository.CreateWallet(r.Context(), &wallet.Data{
		Metadata: data.Metadata,
		Name:     data.Name,
	})
	if err != nil {
		internalError(w, r, err)
		return
	}

	ok(w, wallet)
}
