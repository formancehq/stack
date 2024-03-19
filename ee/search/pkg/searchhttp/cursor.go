package searchhttp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/search/pkg/searchengine"
)

type cursorTokenInfo struct {
	Target      string              `json:"target"`
	Sort        []searchengine.Sort `json:"sort"`
	SearchAfter []interface{}       `json:"searchAfter"`
	Ledgers     []string            `json:"ledgers"`
	PageSize    uint64              `json:"pageSize"`
	TermPolicy  string              `json:"termPolicy"`
	Reverse     bool                `json:"reverse"`
	Terms       []string            `json:"terms"`
}

func DecodeCursorToken(v string, c *cursorTokenInfo) error {
	return json.NewDecoder(base64.NewDecoder(base64.URLEncoding, bytes.NewBufferString(v))).Decode(&c)
}

func EncodeCursorToken(c cursorTokenInfo) string {
	buf := bytes.NewBufferString("")
	enc := base64.NewEncoder(base64.URLEncoding, buf)
	err := json.NewEncoder(enc).Encode(c)
	if err != nil {
		panic(err)
	}
	enc.Close()
	return buf.String()
}

type Total struct {
	Value uint64 `json:"value,omitempty"`
	Rel   string `json:"relation,omitempty"`
}

type Cursor[T any] struct {
	bunpaginate.Cursor[T]
	Total Total `json:"value"`
}

type BaseResponse[T any] struct {
	Cursor *Cursor[T] `json:"cursor,omitempty"`
}
