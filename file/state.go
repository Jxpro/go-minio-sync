package file

import (
	"encoding/json"
	"os"
)

type State struct {
	// FilePath 表示待上传文件的本地路径
	FilePath string `json:"file_path"`
	// FileSize 表示文件的总大小（单位：字节）
	FileSize int64 `json:"file_size"`
	// UploadID 用于标识当前的 multipart upload 任务
	UploadID string `json:"upload_id"`
	// TrunkSize 表示每个分块的大小（单位：字节）
	TrunkSize int `json:"trunk_size"`
	// TrunkLength 表示每个分块的大小（单位：字节）
	TrunkLength int `json:"trunk_length"`
}

func (s *State) Save() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath+".upload.state", data, 0644)
}

func LoadState(stateFile []byte) (*State, error) {
	var s State
	if err := json.Unmarshal(stateFile, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
