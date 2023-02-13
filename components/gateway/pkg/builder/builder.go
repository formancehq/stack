package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/caddyserver/xcaddy"
)

type Builder struct {
	configPath string
	outputPath string
}

func NewBuilder(configPath, outputPath string) *Builder {
	return &Builder{
		configPath: configPath,
		outputPath: outputPath,
	}
}

func (b *Builder) Build(ctx context.Context) error {
	xcaddyBuilder, err := loadConfig(b.configPath)
	if err != nil {
		return err
	}

	return xcaddyBuilder.Build(ctx, b.outputPath)
}

//------------------------------------------------------------------------------

func loadConfig(path string) (xcaddy.Builder, error) {
	var builder xcaddy.Builder

	configFile, err := os.Open(path)
	if err != nil {
		return builder, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&builder); err != nil {
		return builder, fmt.Errorf("failed to parse config file: %w", err)
	}

	return builder, nil
}
