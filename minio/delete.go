package minio

import (
	"context"
	"go-minio-sync/config"

	"github.com/minio/minio-go/v7"
)

func (c *Client) DeleteObject(ctx context.Context, cfg *config.Config, filePath string) error {
	return c.minioClient.RemoveObject(
		ctx,
		cfg.Minio.Bucket,
		cfg.Minio.UserPrefix+filePath,
		minio.RemoveObjectOptions{})
}
