//go:generate task generate-client
package main

import "github.com/numary/auth/cmd"

func main() {
	cmd.Execute()
}
