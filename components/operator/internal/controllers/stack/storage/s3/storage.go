package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	PutFile(fileName string, file []byte) error
	Exist(fileName string) (bool, error)
}

type S3Storage struct {
	bucket  string
	session *session.Session
}

func (s S3Storage) PutFile(fileName string, file []byte) error {

	client := s3manager.NewUploader(s.session)
	input := &s3manager.UploadInput{
		Bucket:          aws.String(s.bucket),
		Key:             aws.String(fileName),
		ContentEncoding: aws.String("gzip"),
		// ContentType: aws.String(""),
		Body: bytes.NewReader(file),
	}
	_, err := client.Upload(input)
	return err
}

func (s S3Storage) Exist(fileName string) (bool, error) {
	client := s3.New(s.session)

	_, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		// Should check if the error is 404
		// if err.Error() == "NotFound: Not Found" {
		// 	return false, nil
		// }

		return false, err
	}

	return true, nil
}

func NewS3Storage(session *session.Session, bucket string) Storage {
	return &S3Storage{
		bucket:  bucket,
		session: session,
	}
}
