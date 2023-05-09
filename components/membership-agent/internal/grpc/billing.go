package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{}   `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		ByIndex struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"by_index"`
	} `json:"aggregations"`
}

func CountDocument(indexName string) (int64, error) {
	// Generate Date
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Create a client
	client, err := opensearch.NewDefaultClient()
	if err != nil {
		return 0, err
	}

	// Create Query for OpenSearch
	content := fmt.Sprintf(`{
		"size": "0",
		"query": {
		  "bool": {
			"filter": [
			  {
				"match_phrase": {
				  "kind": "TRANSACTION"
				}
			  },
			  {
				"range": {
				  "when": {
					"gte": "%s",
					"lte": "%s",
					"format": "strict_date_optional_time"
				  }
				}
			  }
			]
		  }
		},
		"aggs" : {
			  "by_index" : {
				  "terms" : {
					  "field" : "_index"
				 }
			  }
		   }
	}`, firstOfMonth.Format(time.RFC3339), lastOfMonth.Format(time.RFC3339))
	contentBuffer := strings.NewReader(content)
	search := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body:  contentBuffer,
	}

	// Execyte Query
	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		return 0, err
	}

	// Parse Response
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, searchResponse.Body)
	var response Response
	err = json.Unmarshal([]byte(buf.String()), &response)

	return int64(response.Hits.Total.Value), err
}
