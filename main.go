//go:generate docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:latest generate  -i ./swagger.yaml -g go -o ./client --git-user-id=formancehq --git-repo-id=search -p packageVersion=latest -p isGoSubmodule=true -p packageName=client
package main

import "github.com/numary/search/cmd"

func main() {
	cmd.Execute()
}
