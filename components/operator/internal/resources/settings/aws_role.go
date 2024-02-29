package settings

import "github.com/formancehq/operator/internal/core"

func GetAWSServiceAccount(ctx core.Context, stackName string) (string, error) {
	return GetStringOrEmpty(ctx, stackName, "aws.service-account")
}
