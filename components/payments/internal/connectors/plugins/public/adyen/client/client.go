package client

import (
	"sync"

	"github.com/adyen/adyen-go-api-library/v7/src/adyen"
	"github.com/adyen/adyen-go-api-library/v7/src/common"
	"github.com/adyen/adyen-go-api-library/v7/src/management"
)

type Client struct {
	client *adyen.APIClient

	webhookUsername string
	webhookPassword string

	companyID string

	webhooks management.Webhook
	rwMutex  sync.RWMutex
}

func New(apiKey, username, password, companyID string, liveEndpointPrefix string) (*Client, error) {
	adyenConfig := &common.Config{
		ApiKey:      apiKey,
		Environment: common.TestEnv,
		Debug:       true,
	}

	if liveEndpointPrefix != "" {
		adyenConfig.Environment = common.LiveEnv
		adyenConfig.LiveEndpointURLPrefix = liveEndpointPrefix
		adyenConfig.Debug = false
	}

	client := adyen.NewClient(adyenConfig)

	return &Client{
		client:          client,
		webhookUsername: username,
		webhookPassword: password,
	}, nil
}
