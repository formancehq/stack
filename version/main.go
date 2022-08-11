package main

import (
	"encoding/json"
	"net/http"
)

type Info struct {
	Name     string
	Version  string
	Type     string
	Features []Feature
}

type Feature struct {
	Name    string
	Enabled bool
}

func getVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := Info{
		Name:    "Formance Stack",
		Version: "1.0.0",
		Type:    "Community",
		Features: []Feature{
			{Name: "Reports", Enabled: true},
			{Name: "Webhook", Enabled: true},
			{Name: "Reconciliation", Enabled: true},
		},
	}
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", getVersions)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
