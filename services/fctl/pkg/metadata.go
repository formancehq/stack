package fctl

import (
	"fmt"
	"strings"
)

func ParseMetadata(array []string) (map[string]any, error) {
	metadata := map[string]interface{}{}
	for _, v := range array {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed metadata: %s", v)
		}
		metadata[parts[0]] = parts[1]
	}
	return metadata, nil
}
