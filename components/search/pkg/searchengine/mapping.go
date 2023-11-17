package searchengine

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

//go:embed indexed_mapping.json
var indexedMappingJSON string

func CreateIndex(ctx context.Context, client *opensearch.Client, index string) error {
	indexCreateBody, err := GetIndexDefinition()
	if err != nil {
		return err
	}

	_, err = client.Indices.Create(
		index,
		client.Indices.Create.WithContext(ctx),
		func(request *opensearchapi.IndicesCreateRequest) {
			request.Body = bytes.NewReader(indexCreateBody)
		})
	return err
}

func UpdateMapping(ctx context.Context, client *opensearch.Client, index string) error {
	updateMapping, err := json.Marshal(getMapping())
	if err != nil {
		return err
	}

	res, err := client.Indices.PutMapping(
		bytes.NewReader(updateMapping),
		client.Indices.PutMapping.WithContext(ctx),
		client.Indices.PutMapping.WithIndex(index),
	)

	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("request ended with status : %s", res.Status())
	}

	return nil
}

func GetIndexDefinition() ([]byte, error) {
	return json.Marshal(struct {
		Mapping Mappings `json:"mappings"`
	}{
		Mapping: getMapping(),
	})
}

func getMapping() Mappings {
	indexedMapping := map[string]Property{}
	if err := json.Unmarshal([]byte(indexedMappingJSON), &indexedMapping); err != nil {
		panic(err)
	}

	f := false
	return Mappings{
		Properties: map[string]Property{
			"kind": {
				Type: "keyword",
			},
			"ledger": {
				Type: "keyword",
			},
			"stack": {
				Type: "keyword",
			},
			"when": {
				Type: "date",
			},
			"data": {
				Type:    "object",
				Enabled: &f,
			},
			"indexed": {
				Type: "object",
				Mappings: Mappings{
					Properties: indexedMapping,
				},
			},
		},
	}
}

type Property struct {
	Mappings
	Type    string `json:"type,omitempty"`
	Store   bool   `json:"store,omitempty"`
	CopyTo  string `json:"copy_to,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type DynamicTemplate map[string]interface{}

type Mappings struct {
	DynamicTemplates []DynamicTemplate   `json:"dynamic_templates,omitempty"`
	Properties       map[string]Property `json:"properties,omitempty"`
}

type Template struct {
	IndexPatterns []string `json:"index_patterns"`
	Mappings      Mappings `json:"mappings"`
}

//
//func DefaultMapping(patterns ...string) Template {
//	f := false
//	return Template{
//		IndexPatterns: patterns,
//		Mappings: Mappings{
//			DynamicTemplates: []DynamicTemplate{
//				{
//					"strings": map[string]interface{}{
//						"match_mapping_type": "string",
//						"mapping": map[string]interface{}{
//							"type": "keyword",
//						},
//					},
//				},
//			},
//			Properties: map[string]Property{
//				"kind": {
//					Type: "keyword",
//				},
//				"ledger": {
//					Type: "keyword",
//				},
//				"when": {
//					Type: "date",
//				},
//				"data": {
//					Type:    "object",
//					Enabled: &f,
//				},
//				"indexed": {
//					Type: "object",
//				},
//			},
//		},
//	}
//}
//
//func LoadMapping(ctx context.Context, client *opensearch.Client, m Template) error {
//	data, err := json.Marshal(m)
//	if err != nil {
//		return err
//	}
//
//	res, err := opensearchapi.IndicesPutTemplateRequest{
//		Body: bytes.NewReader(data),
//		Name: "search_mapping",
//	}.Do(ctx, client)
//
//	if err != nil {
//		return err
//	}
//	if res.IsError() {
//		return errors.New(res.String())
//	}
//	return nil
//}
//
//func LoadDefaultMapping(ctx context.Context, client *opensearch.Client, indices ...string) error {
//	return LoadMapping(ctx, client, DefaultMapping(indices...))
//}
