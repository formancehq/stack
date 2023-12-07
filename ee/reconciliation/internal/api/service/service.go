package service

import (
	sdk "github.com/formancehq/formance-sdk-go"
)

type Service struct {
	client *sdk.Formance
}

func NewService(client *sdk.Formance) *Service {
	return &Service{
		client: client,
	}
}
