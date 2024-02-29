package settings

import "github.com/formancehq/operator/internal/core"

func GetAWSRole(ctx core.Context, stackName string) (string, error) {
	return GetStringOrEmpty(ctx, stackName, "aws.service-account")
}
