package awss3

import (
	"github.com/antoniomiletta/fileman/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Backend struct {
	client *s3.Client
	bucket string
}

func NewS3Backend(cfg config.StorageConfig) (*S3Backend, error)
