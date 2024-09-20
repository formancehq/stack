package client

import (
	"github.com/adyen/adyen-go-api-library/v7/src/adyen"
	"github.com/adyen/adyen-go-api-library/v7/src/common"
	"github.com/formancehq/go-libs/logging"
)

type Client struct {
	client *adyen.APIClient

	HMACKey string

	logger logging.Logger
}

func NewClient(apiKey, hmacKey, liveEndpointPrefix string, logger logging.Logger) (*Client, error) {
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
		client:  client,
		HMACKey: hmacKey,
		logger:  logger,
	}, nil
}
