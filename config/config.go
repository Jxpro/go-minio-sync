package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Minio MinioConfig
	Chunk ChunkConfig
	Watch WatchConfig
	TLS   TLSConfig
}

type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UserPrefix string
	Bucket     string
	UseSSL     bool
}
type ChunkConfig struct {
	Size int
}

type WatchConfig struct {
	Dir   string
	Delay int
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

// LoadConfig 读取并解析配置文件
func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
