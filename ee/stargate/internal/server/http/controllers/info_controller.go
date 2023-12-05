package controllers

import (
	"encoding/json"
	"net/http"
)

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
