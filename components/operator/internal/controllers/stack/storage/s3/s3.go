package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getConfig(awsClientId, awsSecretId, endpoint, region string, forceStylePath, insecure bool) *aws.Config {
	return aws.NewConfig().WithCredentials(
		credentials.NewStaticCredentials(
			awsClientId,
			awsSecretId,
			"",
		)).WithS3ForcePathStyle(
		forceStylePath,
	).WithEndpoint(endpoint).WithRegion(region).WithDisableSSL(insecure)

}

func NewSession(accessKey, secretKey, endpoint, region string, forceStylePath, insecure bool) (*session.Session, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *getConfig(accessKey, secretKey, endpoint, region, forceStylePath, insecure),
	}))

	_, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return sess, nil
}
