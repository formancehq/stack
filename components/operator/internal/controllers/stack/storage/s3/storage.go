package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	PutFile(fileName string, file []byte) error
}

type S3Storage struct {
	bucket string
	client *s3manager.Uploader
}

func (s S3Storage) PutFile(fileName string, file []byte) error {
	input := &s3manager.UploadInput{
		Bucket:          aws.String(s.bucket),
		Key:             aws.String(fileName),
		ContentEncoding: aws.String("gzip"),
		// ContentType: aws.String(""),
		Body: bytes.NewReader(file),
	}
	_, err := s.client.Upload(input)
	return err
}

func NewS3Storage(client *s3manager.Uploader, bucket string) Storage {
	return &S3Storage{
		client: client,
		bucket: bucket,
	}
}
