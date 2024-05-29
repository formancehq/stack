package models

import (
	"encoding/json"
	"fmt"
)

type Filter struct {
	MatchPhrase map[string]interface{} `json:"match_phrase"`
}

const (
	esQuery = `{
		"query": {
			"bool": {
				"filter": %s
			}
		}
	}
	`
)

func BuildESQuery(filters []Filter) (map[string]interface{}, error) {
	b, err := json.Marshal(filters)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(esQuery, string(b))

	var res map[string]interface{}
	if err := json.Unmarshal([]byte(query), &res); err != nil {
		return nil, err
	}

	return res, nil
}
