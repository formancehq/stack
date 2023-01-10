// docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:v6.2.1 generate  -i ./pkg/api/controllers/swagger.yaml -g go -o ./client --git-user-id=formancehq --git-repo-id=ledger -p packageVersion=latest -p isGoSubmodule=true -p packageName=client
//
//go:generate docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:v6.2.1 validate -i ./pkg/api/controllers/swagger.yaml
package main

import (
	"github.com/numary/ledger/cmd"
)

func main() {
	cmd.Execute()
}
