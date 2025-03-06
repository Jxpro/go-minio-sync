package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadState(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "temp1g_*")
	require.NoError(t, err, "文件创建失败")
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	const oneGB = int64(1024 * 1024 * 1024)
	err = tmpFile.Truncate(oneGB)
	assert.NoError(t, err, "文件截断失败")

	state := &State{
		FilePath:    tmpFile.Name(),
		UploadID:    "12345",
		FileSize:    123456789,
		TrunkSize:   5 * 1024 * 1024, // 5MB
		TrunkLength: 123456789 / (5 * 1024 * 1024),
	}

	err = state.Save()
	require.NoError(t, err, "文件保存失败")

	data, err := os.ReadFile(tmpFile.Name() + ".upload.state")
	require.NoError(t, err, "文件读取失败")
	loaded, err := LoadState(data)
	require.NoError(t, err, "文件解析失败")

	assert.Equal(t, state.UploadID, loaded.UploadID, "UploadID 不一致")
	assert.Equal(t, state.TrunkLength, loaded.TrunkLength, "分块大小不一致")
	assert.Equal(t, state.FileSize, loaded.FileSize, "文件大小不一致")
}
