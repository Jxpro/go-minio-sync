package minio

import (
	"context"
	"go-minio-sync/config"

	"github.com/minio/minio-go/v7"
)

func (c *Client) MakeBucketIfNotExists(ctx context.Context, cfg *config.Config) error {
	exists, err := c.minioClient.BucketExists(ctx, cfg.Minio.Bucket)
	if err != nil {
		return err
	}
	if !exists {
		return c.minioClient.MakeBucket(ctx, cfg.Minio.Bucket, minio.MakeBucketOptions{})
	}
	return nil
}
