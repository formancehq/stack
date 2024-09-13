package models

import "encoding/json"

type PSPOther struct {
	ID    string          `json:"id"`
	Other json.RawMessage `json:"other"`
}
