package main

import (
	"github.com/formancehq/stack/components/agent/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
