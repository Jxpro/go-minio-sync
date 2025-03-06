package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"go-minio-sync/config"
)

func (c *Client) ListObjects(ctx context.Context, cfg *config.Config) ([]string, error) {
	var objectNames []string
	for object := range c.minioClient.ListObjects(ctx, cfg.Minio.Bucket, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			return nil, object.Err
		}
		objectNames = append(objectNames, object.Key)
	}
	return objectNames, nil
}
