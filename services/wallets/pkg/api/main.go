package api

import (
	"github.com/formancehq/wallets/pkg/wallet"
)

type MainHandler struct {
	funding    *wallet.FundingService
	repository *wallet.Repository
}

func NewMainHandler(
	funding *wallet.FundingService,
	repository *wallet.Repository,
) *MainHandler {
	return &MainHandler{
		funding:    funding,
		repository: repository,
	}
}
