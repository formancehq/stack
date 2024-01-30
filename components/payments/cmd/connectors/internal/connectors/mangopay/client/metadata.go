package client

const (
	mangopayMetadataSpecNamespace = "com.mangopay.spec/"
)

func ExtractNamespacedMetadata(metadata map[string]string, key string) string {
	value, ok := metadata[mangopayMetadataSpecNamespace+key]
	if !ok {
		return ""
	}
	return value
}
