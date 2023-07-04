package internal

import (
	"encoding/json"
	"fmt"
)

type Metadata map[string]interface{}

func MetadataInterfaceAsShortString(metadata map[string]interface{}) string {
	metadataAsString := ""
	for k, v := range metadata {
		asJson, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		metadataAsString += fmt.Sprintf("%s=%s ", k, string(asJson))
	}
	if len(metadataAsString) > 100 {
		metadataAsString = metadataAsString[:100] + "..."
	}
	return metadataAsString
}
