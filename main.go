//go:generate docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:latest generate  -i ./pkg/server/openapi.yaml -g go -o ./client --git-user-id=formancehq --git-repo-id=webhooks -p packageVersion=latest -p isGoSubmodule=true -p packageName=client
package main

import "github.com/formancehq/webhooks/cmd"

func main() {
	cmd.Execute()
}
