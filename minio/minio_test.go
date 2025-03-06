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

func TestNewClient(t *testing.T) {
	client, err := NewClient(&cfg)
	require.NoError(t, err, "初始化 MinIO 客户端出错")

	err = client.MakeBucketIfNotExists(context.Background(), &cfg)
	assert.NoError(t, err, "创建或检查桶时出错")
}
