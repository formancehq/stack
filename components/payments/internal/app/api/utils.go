package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/query"
)

func getQueryBuilder(r *http.Request) (query.Builder, error) {
	return query.ParseJSON(r.URL.Query().Get("query"))
}
