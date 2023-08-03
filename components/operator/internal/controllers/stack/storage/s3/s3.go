package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

// According to https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// AWS_S3_Force_Path_Style
func NewClient(accessKey, secretKey, endpoint, region string, forceStylePath, insecure bool) (*s3manager.Uploader, error) {

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *getConfig(accessKey, secretKey, endpoint, region, forceStylePath, insecure),
	}))

	_, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	// Create an uploader with the session and default options
	return s3manager.NewUploader(sess), nil
}
