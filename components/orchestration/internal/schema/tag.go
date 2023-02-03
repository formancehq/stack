package schema

import (
	"strings"
)

type tag struct {
	defaultValue string
}

func parseTag(tagValue string) tag {
	parts := strings.Split(tagValue, ",")
	ret := tag{}
	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "default:"):
			ret.defaultValue = strings.TrimPrefix(part, "default:")
		}
	}
	return ret
}
