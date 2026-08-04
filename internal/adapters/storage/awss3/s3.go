package awss3

import (
	"github.com/antoniomiletta/fileman/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(cfg config.S3Config) *S3Storage {
	return &S3Storage{
		client: s3.New(s3.Options{}),
		bucket: cfg.S3Bucket,
	}
}
