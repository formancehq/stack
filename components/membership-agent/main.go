//go:generate protoc --go_out=internal/grpc/generated --go_opt=paths=source_relative --go-grpc_out=internal/grpc/generated --go-grpc_opt=paths=source_relative server.proto
package main

import (
	"github.com/formancehq/stack/components/membership-agent/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
