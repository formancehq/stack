//go:generate docker run --rm -w /local -v ${PWD}:/local openapitools/openapi-generator-cli:v6.2.1 validate -i ./swagger.yaml
package main

import "github.com/formancehq/search/cmd"

func main() {
	cmd.Execute()
}
