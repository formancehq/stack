package searchhttp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/numary/search/pkg/searchengine"
)

type Total struct {
	Value uint64 `json:"value"`
	Rel   string `json:"relation"`
}

type Page struct {
	PageSize int         `json:"pageSize"`
	HasMore  bool        `json:"hasMore"`
	Total    Total       `json:"total,omitempty"`
	Next     string      `json:"next,omitempty"`
	Data     interface{} `json:"data"`
	Previous string      `json:"previous"`
}

type cursorTokenInfo struct {
	Target      string              `json:"target"`
	Sort        []searchengine.Sort `json:"sort"`
	SearchAfter []interface{}       `json:"searchAfter"`
	Ledgers     []string            `json:"ledgers"`
	Size        uint64              `json:"size"`
	TermPolicy  string              `json:"termPolicy"`
	Reverse     bool                `json:"reverse"`
	Terms       []string            `json:"terms"`
}

func DecodeCursorToken(v string, c *cursorTokenInfo) error {
	return json.NewDecoder(base64.NewDecoder(base64.URLEncoding, bytes.NewBufferString(v))).Decode(&c)
}

func EncodePaginationToken(c cursorTokenInfo) string {
	buf := bytes.NewBufferString("")
	enc := base64.NewEncoder(base64.URLEncoding, buf)
	err := json.NewEncoder(enc).Encode(c)
	if err != nil {
		panic(err)
	}
	enc.Close()
	return buf.String()
}
