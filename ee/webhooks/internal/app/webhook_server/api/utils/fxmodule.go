package utils

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DefaultServerParams struct {
	Addr   string
	Info   ServiceInfo
	Auth   auth.Auth
	Logger logging.Logger
}
