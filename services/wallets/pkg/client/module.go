package client

import (
	"github.com/formancehq/formance-sdk-go"
	"go.uber.org/fx"
)

func NewModule(clientID string, clientSecret string, tokenURL string) fx.Option {
	return fx.Provide(func() (*formance.APIClient, error) {
		return NewStackClient(clientID, clientSecret, tokenURL)
	})
}
