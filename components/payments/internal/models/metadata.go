package models

const (
	formanceMetadataSpecNamespace = "com.formance.spec/"
)

func ExtractNamespacedMetadata(metadata map[string]string, key string) string {
	value, ok := metadata[key]
	if !ok {
		return ""
	}
	return value
}
