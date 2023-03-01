package main

import (
	"fmt"

	"github.com/formancehq/formance-sdk-go"
)

func main() {
	configuration := formance.Configuration{}
	_ = formance.NewAPIClient(&configuration)
	fmt.Println("TODO: Actually just checking we can compile the SDK and create a client")
}
