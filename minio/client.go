package minio

import (
	"github.com/minio/minio-go/v7"
)

type Client struct {
	Bucket      string
	minioCore   *minio.Core
	minioClient *minio.Client
}
