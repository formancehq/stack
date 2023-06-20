package controllers

import (
	"encoding/json"
	"net/http"
)

type StargateControllerConfig struct {
	version string
}

func NewStargateControllerConfig(
	version string,
) StargateControllerConfig {
	return StargateControllerConfig{
		version: version,
	}
}

type StargateController struct {
	config StargateControllerConfig
}

func NewStargateController(
	config StargateControllerConfig,
) *StargateController {
	return &StargateController{
		config: config,
	}
}

type ServiceInfo struct {
	Version string `json:"version"`
}

func (s *StargateController) GetInfo(w http.ResponseWriter, r *http.Request) {
	info := ServiceInfo{
		Version: s.config.version,
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		panic(err)
	}
}
