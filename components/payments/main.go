//go:generate docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:v6.2.1 generate  -i ./swagger.yml -g go -o ./client --git-user-id=formancehq --git-repo-id=payments -p packageVersion=latest -p isGoSubmodule=true -p packageName=client
//go:generate docker run --rm -w /local -v ${PWD}:/local cytopia/goimports -w -e ./client
package main

import "github.com/formancehq/payments/cmd"

func main() {
	cmd.Execute()
}
