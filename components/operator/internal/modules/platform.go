package modules

type Platform struct {
	// Cloud region where the stack is deployed
	Region string
	// Cloud environment where the stack is deployed: staging, production,
	// sandbox, etc.
	Environment string
}
