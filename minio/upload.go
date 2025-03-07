package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"go-minio-sync/config"
	"go-minio-sync/sync"
)

func (c *Client) UploadFileWithResume(ctx context.Context, cfg *config.Config, filePath string) error {
	var st *sync.State
	stateFile := filePath + ".upload.state"
	objectName := cfg.Minio.UserPrefix + filePath

	// 尝试读取状态文件，如果存在则恢复 UploadID
	if data, err := os.ReadFile(stateFile); err == nil {
		st, err = sync.LoadState(data)
		if err != nil {
			_ = os.Remove(stateFile)
		}
	}

	// 如果没有保存的 UploadID，则初始化新的 Multipart Upload 会话
	if st == nil {
		UploadID, err := c.minioCore.NewMultipartUpload(ctx, cfg.Minio.Bucket, objectName, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		st = &sync.State{
			FilePath:    filePath,
			UploadID:    UploadID,
			FileSize:    fileInfo.Size(),
			TrunkSize:   cfg.Chunk.Size,
			TrunkLength: int(fileInfo.Size() / int64(cfg.Chunk.Size)),
		}
		err = st.Save()
		if err != nil {
			return err
		}
	}

	// 若 state 仍为空，则返回错误
	if st == nil {
		return errors.New("state is nil")
	}

	// 通过 ListObjectParts 获取已上传部分的信息
	partsInfo, err := c.minioCore.ListObjectParts(
		ctx, cfg.Minio.Bucket, objectName, st.UploadID, 0, st.TrunkLength)
	if err != nil {
		return err
	}

	uploadedPartsMap := make(map[int]minio.ObjectPart)
	for _, part := range partsInfo.ObjectParts {
		uploadedPartsMap[part.PartNumber] = part
	}

	// 打开待上传文件
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	var parts []minio.CompletePart
	partNumber := 1

	// 迭代读取并上传每个分块
	for {
		buffer := make([]byte, cfg.Chunk.Size)
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 { // 文件读取完毕
			break
		}

		// 如果该分块已存在（断点续传场景），则直接记录，不重复上传
		if existing, exists := uploadedPartsMap[partNumber]; exists {
			parts = append(parts, minio.CompletePart{
				PartNumber: partNumber,
				ETag:       existing.ETag,
			})
		} else {
			// 上传缺失的分块
			partReader := bytes.NewReader(buffer[:n])
			partUpload, err := c.minioCore.PutObjectPart(
				ctx, cfg.Minio.Bucket, objectName, st.UploadID, partNumber, partReader, int64(n), minio.PutObjectPartOptions{})
			if err != nil {
				return err
			}
			parts = append(parts, minio.CompletePart{
				PartNumber: partNumber,
				ETag:       partUpload.ETag,
			})
		}
		partNumber++
	}

	// 完成上传：通知服务器合并所有分块
	_, err = c.minioCore.CompleteMultipartUpload(
		ctx, cfg.Minio.Bucket, objectName, st.UploadID, parts, minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	// 上传成功后删除状态文件
	_ = os.Remove(stateFile)
	return nil
}
