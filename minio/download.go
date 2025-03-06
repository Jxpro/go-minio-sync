package minio

import (
	"context"
	"go-minio-sync/config"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
)

func (c *Client) DownloadObject(ctx context.Context, cfg *config.Config, filePath string) error {
	obj, err := c.minioClient.GetObject(
		ctx, cfg.Minio.Bucket, cfg.Minio.UserPrefix+filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = obj.Close()
	}()

	if _, err = obj.Stat(); err != nil {
		return err
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outFile.Close()
	}()

	_, err = io.Copy(outFile, obj)
	if err != nil {
		return err
	}
	return nil
}
