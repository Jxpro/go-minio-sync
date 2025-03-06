package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试加载配置文件，验证返回的配置数据是否有效
func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../config.yaml")
	assert.NoError(t, err, "加载配置时不应出错")
	assert.NotNil(t, cfg, "配置对象不应为空")
}
