package fctl

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/formancehq/go-libs/metadata"
	"github.com/pterm/pterm"
)

func ParseMetadata(array []string) (metadata.Metadata, error) {
	metadata := metadata.Metadata{}
	for _, v := range array {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed metadata: %s", v)
		}
		metadata[parts[0]] = parts[1]
	}
	return metadata, nil
}

func PrintMetadata(out io.Writer, metadata metadata.Metadata) error {
	Section.WithWriter(out).Println("Metadata")
	if len(metadata) == 0 {
		Println("No metadata.")
		return nil
	}
	tableData := pterm.TableData{}
	for k, v := range metadata {
		tableData = append(tableData, []string{pterm.LightCyan(k), v})
	}

	return pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render()
}

func MetadataAsShortString[V any](metadata map[string]V) string {
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
