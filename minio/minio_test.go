package minio

import (
	"context"
	. "go-minio-sync/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg = Config{
	Minio: MinioConfig{
		Endpoint:   "localhost:9000",
		AccessKey:  "minioadmin",
		SecretKey:  "minioadmin",
		UserPrefix: "test-user/",
		Bucket:     "test-bucket",
		UseSSL:     false,
	},
	Chunk: ChunkConfig{
		Size: 5 * 1024 * 1024,
	},
	TLS: TLSConfig{
		CertFile: "",
		KeyFile:  "",
	},
}

const (
	testFilePath = "C:\\Users\\xin\\Downloads\\开题报告.zip"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(&cfg)
	require.NoError(t, err, "初始化 MinIO 客户端出错")

	err = client.MakeBucketIfNotExists(context.Background(), &cfg)
	assert.NoError(t, err, "创建或检查桶时出错")
}

func TestUploadObject(t *testing.T) {
	client, err := NewClient(&cfg)
	require.NoError(t, err, "初始化 MinIO 客户端出错")

	err = client.UploadFileWithResume(context.Background(), &cfg, testFilePath)
	require.NoError(t, err, "上传对象时出错")
}

func TestDownloadObject(t *testing.T) {
	client, err := NewClient(&cfg)
	require.NoError(t, err, "初始化 MinIO 客户端出错")

	err = client.DownloadObject(context.Background(), &cfg, testFilePath)
	require.NoError(t, err, "下载对象时出错")
}

func TestDeleteObject(t *testing.T) {
	client, err := NewClient(&cfg)
	require.NoError(t, err, "初始化 MinIO 客户端出错")

	err = client.DeleteObject(context.Background(), &cfg, testFilePath)
	require.NoError(t, err, "删除对象时出错")

	err = client.DownloadObject(context.Background(), &cfg, testFilePath)
	assert.Error(t, err, "删除后下载对象应失败")
}
