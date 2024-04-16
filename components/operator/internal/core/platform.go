package core

type Platform struct {
	// Cloud region where the stack is deployed
	Region string
	// Cloud environment where the stack is deployed: staging, production,
	// sandbox, etc.
	Environment string
	// The licence information
	Licence Licence
}

type Licence struct {
	// The licence token
	Token string
	// The licence issuer
	Issuer string
	// The licence clusterID
	ClusterID string
}
